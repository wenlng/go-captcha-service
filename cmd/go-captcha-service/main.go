package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/sony/gobreaker"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/middleware"
	"github.com/wenlng/go-captcha-service/internal/server"
	"github.com/wenlng/go-captcha-service/internal/service_discovery"
	"github.com/wenlng/go-captcha-service/proto"
)

// App manages the application components
type App struct {
	logger       *zap.Logger
	dynamicCfg   *config.DynamicConfig
	cache        cache.Cache
	discovery    service_discovery.ServiceDiscovery
	httpServer   *http.Server
	grpcServer   *grpc.Server
	cacheBreaker *gobreaker.CircuitBreaker
	limiter      *middleware.DynamicLimiter
}

const (
	CacheTypeRedis    string = "redis"
	CacheTypeMemory          = "memory"
	CacheTypeEtcd            = "etcd"
	CacheTypeMemcache        = "memcache"
)

const (
	ServiceDiscoveryTypeEtcd      string = "etcd"
	ServiceDiscoveryTypeZookeeper        = "zookeeper"
	ServiceDiscoveryTypeConsul           = "consul"
	ServiceDiscoveryTypeNacos            = "nacos"
)

// NewApp initializes the application
func NewApp() (*App, error) {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	// Parse command-line flags
	configFile := flag.String("config", "config.json", "Path to config file")
	serviceName := flag.String("service-name", "", "Name for service")
	httpPort := flag.String("http-port", "", "Port for HTTP server")
	grpcPort := flag.String("grpc-port", "", "Port for gRPC server")
	redisAddrs := flag.String("redis-addrs", "", "Comma-separated Redis cluster addresses")
	etcdAddrs := flag.String("etcd-addrs", "", "Comma-separated etcd addresses")
	memcacheAddrs := flag.String("memcache-addrs", "", "Comma-separated Memcached addresses")
	cacheType := flag.String("cache-type", "", "Cache type: redis, memory, etcd, memcache")
	cacheTTL := flag.Int("cache-ttl", 0, "Cache TTL in seconds")
	cacheCleanupInt := flag.Int("cache-cleanup-interval", 0, "Cache cleanup interval in seconds")
	cacheKeyPrefix := flag.Int("cache-key-prefix", 0, "Key prefix for cache")
	serviceDiscovery := flag.String("service-discovery", "", "Service discovery: etcd, zookeeper, consul, nacos")
	serviceDiscoveryAddrs := flag.String("service-discovery-addrs", "", "Service discovery addresses")
	rateLimitQPS := flag.Int("rate-limit-qps", 0, "Rate limit QPS")
	rateLimitBurst := flag.Int("rate-limit-burst", 0, "Rate limit burst")
	loadBalancer := flag.String("load-balancer", "", "Load balancer: round-robin, consistent-hash")
	apiKeys := flag.String("api-keys", "", "Comma-separated API keys")
	healthCheckFlag := flag.Bool("health-check", false, "Run health check and exit")
	enableCorsFlag := flag.Bool("enable-cors", false, "Enable cross-domain resources")
	flag.Parse()

	// Load configuration
	dc, err := config.NewDynamicConfig(*configFile)
	if err != nil {
		logger.Warn("Failed to load config, using defaults", zap.Error(err))
		dc = &config.DynamicConfig{Config: config.DefaultConfig()}
	}

	// Merge command-line flags
	cfg := dc.Get()
	cfg = config.MergeWithFlags(cfg, map[string]interface{}{
		"service-name":            *serviceName,
		"http-port":               *httpPort,
		"grpc-port":               *grpcPort,
		"redis-addrs":             *redisAddrs,
		"etcd-addrs":              *etcdAddrs,
		"memcache-addrs":          *memcacheAddrs,
		"cache-type":              *cacheType,
		"cache-ttl":               *cacheTTL,
		"cache-cleanup-interval":  *cacheCleanupInt,
		"cache-key-prefix":        *cacheKeyPrefix,
		"service-discovery":       *serviceDiscovery,
		"service-discovery-addrs": *serviceDiscoveryAddrs,
		"rate-limit-qps":          *rateLimitQPS,
		"rate-limit-burst":        *rateLimitBurst,
		"load-balancer":           *loadBalancer,
		"api-keys":                *apiKeys,
		"enable-cors":             *enableCorsFlag,
	})
	if err = dc.Update(cfg); err != nil {
		logger.Fatal("Configuration validation failed", zap.Error(err))
	}

	// Initialize rate limiter
	limiter := middleware.NewDynamicLimiter(cfg.RateLimitQPS, cfg.RateLimitBurst)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			newCfg := dc.Get()
			limiter.Update(newCfg.RateLimitQPS, newCfg.RateLimitBurst)
		}
	}()

	// Initialize circuit breaker
	cacheBreaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        *serviceName,
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	})

	// Initialize curCache
	var curCache cache.Cache
	ttl := time.Duration(cfg.CacheTTL) * time.Second
	cleanInt := time.Duration(cfg.CacheCleanupInt) * time.Second
	switch cfg.CacheType {
	case CacheTypeRedis:
		curCache, err = cache.NewRedisClient(cfg.RedisAddrs, cfg.CacheKeyPrefix, ttl)
		if err != nil {
			logger.Fatal("Failed to initialize Redis", zap.Error(err))
		}
	case CacheTypeMemory:
		curCache = cache.NewMemoryCache(cfg.CacheKeyPrefix, ttl, cleanInt)
	case CacheTypeEtcd:
		curCache, err = cache.NewEtcdClient(cfg.EtcdAddrs, cfg.CacheKeyPrefix, ttl)
		if err != nil {
			logger.Fatal("Failed to initialize etcd", zap.Error(err))
		}
	case CacheTypeMemcache:
		curCache, err = cache.NewMemcacheClient(cfg.MemcacheAddrs, cfg.CacheKeyPrefix, ttl)
		if err != nil {
			logger.Fatal("Failed to initialize Memcached", zap.Error(err))
		}
	default:
		logger.Fatal("Invalid curCache type", zap.String("type", cfg.CacheType))
	}

	// Initialize service discovery
	var discovery service_discovery.ServiceDiscovery
	if cfg.ServiceDiscovery != "" {
		switch cfg.ServiceDiscovery {
		case ServiceDiscoveryTypeEtcd:
			discovery, err = service_discovery.NewEtcdDiscovery(cfg.ServiceDiscoveryAddrs, 10)
		case ServiceDiscoveryTypeZookeeper:
			discovery, err = service_discovery.NewZookeeperDiscovery(cfg.ServiceDiscoveryAddrs, 10)
		case ServiceDiscoveryTypeConsul:
			discovery, err = service_discovery.NewConsulDiscovery(cfg.ServiceDiscoveryAddrs, 10)
		case ServiceDiscoveryTypeNacos:
			discovery, err = service_discovery.NewNacosDiscovery(cfg.ServiceDiscoveryAddrs, 10)
		default:
			logger.Fatal("Invalid service discovery type", zap.String("type", cfg.ServiceDiscovery))
		}
		if err != nil {
			logger.Fatal("Failed to initialize service discovery", zap.Error(err))
		}
	}

	// Perform health check if requested
	if *healthCheckFlag {
		if err = healthCheck(":"+cfg.HTTPPort, ":"+cfg.GRPCPort); err != nil {
			logger.Error("Health check failed", zap.Error(err))
			os.Exit(1)
		}
		os.Exit(0)
	}

	return &App{
		logger:       logger,
		dynamicCfg:   dc,
		cache:        curCache,
		discovery:    discovery,
		cacheBreaker: cacheBreaker,
		limiter:      limiter,
	}, nil
}

// Start launches the HTTP and gRPC servers
func (a *App) Start(ctx context.Context) error {
	cfg := a.dynamicCfg.Get()

	// setup captcha
	captcha, err := gocaptcha.Setup()
	if err != nil {
		return errors.New("setup gocaptcha failed")
	}

	// Register service with discovery
	var instanceID string
	if a.discovery != nil {
		instanceID = uuid.New().String()
		httpPortInt, _ := strconv.Atoi(cfg.HTTPPort)
		grpcPortInt, _ := strconv.Atoi(cfg.GRPCPort)
		if err = a.discovery.Register(ctx, cfg.ServiceName, instanceID, "127.0.0.1", httpPortInt, grpcPortInt); err != nil {
			return fmt.Errorf("failed to register service: %v", err)
		}
		go a.updateInstances(ctx, instanceID)
	}

	// service context
	svcCtx := common.NewSvcContext()
	svcCtx.Cache = a.cache
	svcCtx.Config = &cfg
	svcCtx.Logger = a.logger
	svcCtx.Captcha = captcha

	// Register HTTP routes
	handlers := server.NewHTTPHandlers(svcCtx)
	mwChain := middleware.NewChainHTTP(
		nil, // .
		middleware.APIKeyMiddleware(a.dynamicCfg, a.logger),
		middleware.LoggingMiddleware(a.logger),
		middleware.RateLimitMiddleware(a.limiter, a.logger),
		middleware.CircuitBreakerMiddleware(a.cacheBreaker, a.logger),
	)

	// Enable cross-domain resource
	if cfg.EnableCors {
		mwChain.AppendMiddleware(middleware.CORSMiddleware(a.logger))
	}

	// Logic Routes
	http.Handle("/get-data", mwChain.Then(handlers.GetDataHandler))
	http.Handle("/check-data", mwChain.Then(handlers.CheckDataHandler))
	http.Handle("/check-status", mwChain.Then(handlers.CheckStatusHandler))
	http.Handle("/get-status-info", mwChain.Then(handlers.GetStatusInfoHandler))
	http.Handle("/del-status-data", mwChain.Then(handlers.DelStatusInfoHandler))
	http.Handle("/rate-limit", mwChain.Then(middleware.RateLimitHandler(a.limiter, a.logger)))

	// Start HTTP server
	a.httpServer = &http.Server{
		Addr: ":" + cfg.HTTPPort,
	}
	go func() {
		a.logger.Info("Starting HTTP server", zap.String("port", cfg.HTTPPort))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryServerInterceptor(a.dynamicCfg, a.logger, a.cacheBreaker)),
	)
	proto.RegisterGoCaptchaServiceServer(a.grpcServer, server.NewGoCaptchaServer(svcCtx))
	go func() {
		a.logger.Info("Starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := a.grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			a.logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	return nil
}

// updateInstances periodically updates service instances
func (a *App) updateInstances(ctx context.Context, instanceID string) {
	ticker := time.NewTicker(10 * time.Second)
	cfg := a.dynamicCfg.Get()

	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			if a.discovery != nil {
				if err := a.discovery.Deregister(ctx, instanceID); err != nil {
					a.logger.Error("Failed to deregister service", zap.Error(err))
				}
			}
			return
		case <-ticker.C:
			if a.discovery == nil {
				continue
			}
			instances, err := a.discovery.Discover(ctx, cfg.ServiceName)
			if err != nil {
				a.logger.Error("Failed to discover instances", zap.Error(err))
				continue
			}
			a.logger.Info("Discovered instances", zap.Int("count", len(instances)))
		}
	}
}

// Shutdown gracefully stops the application
func (a *App) Shutdown() {
	a.logger.Info("Received shutdown signal, shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop HTTP server
	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			a.logger.Error("HTTP server shutdown error", zap.Error(err))
		} else {
			a.logger.Info("HTTP server shut down successfully")
		}
	}

	// Stop gRPC server
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
		a.logger.Info("gRPC server shut down successfully")
	}

	// Close cache
	if redisClient, ok := a.cache.(*cache.RedisClient); ok {
		if err := redisClient.Close(); err != nil {
			a.logger.Error("Redis client close error", zap.Error(err))
		} else {
			a.logger.Info("Redis client closed successfully")
		}
	}
	if memoryCache, ok := a.cache.(*cache.MemoryCache); ok {
		memoryCache.Stop()
		a.logger.Info("Memory cache stopped successfully")
	}
	if etcdClient, ok := a.cache.(*cache.EtcdClient); ok {
		if err := etcdClient.Close(); err != nil {
			a.logger.Error("etcd client close error", zap.Error(err))
		} else {
			a.logger.Info("etcd client closed successfully")
		}
	}
	if memcacheClient, ok := a.cache.(*cache.MemcacheClient); ok {
		if err := memcacheClient.Close(); err != nil {
			a.logger.Error("Memcached client close error", zap.Error(err))
		} else {
			a.logger.Info("Memcached client closed successfully")
		}
	}

	// Close service discovery
	if a.discovery != nil {
		if err := a.discovery.Close(); err != nil {
			a.logger.Error("Service discovery close error", zap.Error(err))
		} else {
			a.logger.Info("Service discovery closed successfully")
		}
	}
}

// healthCheck performs a health check on HTTP and gRPC servers
func healthCheck(httpAddr, grpcAddr string) error {
	resp, err := http.Get("http://localhost" + httpAddr + "/read?key=test")
	if err != nil || resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("HTTP health check failed: %v", err)
	}
	resp.Body.Close()

	conn, err := net.DialTimeout("tcp", "localhost"+grpcAddr, 1*time.Second)
	if err != nil {
		return fmt.Errorf("gRPC health check failed: %v", err)
	}
	conn.Close()

	return nil
}

func main() {
	app, err := NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize app: %v\n", err)
		os.Exit(1)
	}
	defer app.logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = app.Start(ctx); err != nil {
		app.logger.Fatal("Failed to start app", zap.Error(err))
	}

	// Handle termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	app.Shutdown()
	app.logger.Info("App exited")
}
