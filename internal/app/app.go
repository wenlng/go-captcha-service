/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package app

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sony/gobreaker"
	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/middleware"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	config2 "github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha/config"
	"github.com/wenlng/go-captcha-service/internal/server"
	"github.com/wenlng/go-captcha-service/proto"
	"github.com/wenlng/go-service-link/dynaconfig"
	"github.com/wenlng/go-service-link/servicediscovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// App manages the application components
type App struct {
	logger         *zap.Logger
	dynamicCfg     *config.DynamicConfig
	dynamicCaptCfg *config2.DynamicCaptchaConfig
	cacheMgr       *cache.CacheManager
	discovery      servicediscovery.ServiceDiscovery
	configManager  *dynaconfig.ConfigManager
	httpServer     *http.Server
	grpcServer     *grpc.Server
	cacheBreaker   *gobreaker.CircuitBreaker
	limiter        *middleware.DynamicLimiter
	captcha        *gocaptcha.GoCaptcha
}

// NewApp initializes the application
func NewApp() (*App, error) {

	// Parse command-line flags
	configFile := flag.String("config", "config.json", "Path to config file")
	gocaptchaConfigFile := flag.String("gocaptcha-config", "gocaptcha.json", "Path to gocaptcha config file")
	serviceName := flag.String("service-name", "", "Name for service")
	serviceNode := flag.Int64("service-node", 1, "Node number for service")
	httpPort := flag.String("http-port", "", "Port for HTTP server")
	grpcPort := flag.String("grpc-port", "", "Port for gRPC server")

	cacheType := flag.String("cache-type", "", "CacheManager type: redis, memory, etcd, memcache")
	cacheAddrs := flag.String("cache-addrs", "", "Comma-separated Cache cluster addresses")
	cacheUsername := flag.String("cache-username", "", "Cache service username")
	cachePassword := flag.String("cache-password", "", "Cache service password")
	cacheDB := flag.String("cache-dbb", "0", "Cache service db name")
	cacheTTL := flag.Int("cache-ttl", 0, "CacheManager TTL in seconds")
	cacheKeyPrefix := flag.String("cache-key-prefix", "GO_CAPTCHA_DATA:", "Key prefix for cache")

	enableDynamicConfig := flag.String("enable-dynamic-config", "false", "Enable dynamic config")
	dynamicConfigType := flag.String("dynamic-config-type", "", "Service discovery: etcd, zookeeper, consul, nacos")
	dynamicConfigAddrs := flag.String("dynamic-config-addrs", "", "Comma-separated list of service dynamic config addresses")
	dynamicConfigTTL := flag.Int("dynamic-config-ttl", 10, "Time-to-live in seconds for dynamic config registrations")
	dynamicConfigKeepAlive := flag.Int("dynamic-config-keep-alive", 3, "Duration in seconds for dynamic config keep-alive interval")
	dynamicConfigMaxRetries := flag.Int("dynamic-config-max-retries", 3, "Maximum number of retries for dynamic config operations")
	dynamicConfigBaseRetryDelay := flag.Int("dynamic-config-base-retry-delay", 3, "Base delay in milliseconds for dynamic config retry attempts")
	dynamicConfigUsername := flag.String("dynamic-config-username", "", "Username for dynamic config authentication")
	dynamicConfigPassword := flag.String("dynamic-config-password", "", "Password for dynamic config authentication")
	dynamicConfigTlsServerName := flag.String("dynamic-config-tls-server-name", "", "TLS server name for dynamic config connection")
	dynamicConfigTlsAddress := flag.String("dynamic-config-tls-address", "", "TLS address for service dynamic config")
	dynamicConfigTlsCertFile := flag.String("dynamic-config-tls-cert-file", "", "Path to TLS certificate file for dynamic config")
	dynamicConfigTlsKeyFile := flag.String("dynamic-config-tls-key-file", "", "Path to TLS key file for dynamic config")
	dynamicConfigTlsCaFile := flag.String("dynamic-config-tls-ca-file", "", "Path to TLS CA file for dynamic config")

	enableServiceDiscovery := flag.String("enable-service-discovery", "false", "Enable service discovery")
	serviceDiscoveryType := flag.String("service-discovery-type", "", "Service discovery: etcd, zookeeper, consul, nacos")
	serviceDiscoveryAddrs := flag.String("service-discovery-addrs", "", "Comma-separated list of service discovery server addresses")
	serviceDiscoveryTTL := flag.Int("service-discovery-ttl", 10, "Time-to-live in seconds for service discovery registrations")
	serviceDiscoveryKeepAlive := flag.Int("service-discovery-keep-alive", 3, "Duration in seconds for service discovery keep-alive interval")
	serviceDiscoveryMaxRetries := flag.Int("service-discovery-max-retries", 3, "Maximum number of retries for service discovery operations")
	serviceDiscoveryBaseRetryDelay := flag.Int("service-discovery-base-retry-delay", 3, "Base delay in milliseconds for service discovery retry attempts")
	serviceDiscoveryUsername := flag.String("service-discovery-username", "", "Username for service discovery authentication")
	serviceDiscoveryPassword := flag.String("service-discovery-password", "", "Password for service discovery authentication")
	serviceDiscoveryTlsServerName := flag.String("service-discovery-tls-server-name", "", "TLS server name for service discovery connection")
	serviceDiscoveryTlsAddress := flag.String("service-discovery-tls-address", "", "TLS address for service discovery server")
	serviceDiscoveryTlsCertFile := flag.String("service-discovery-tls-cert-file", "", "Path to TLS certificate file for service discovery")
	serviceDiscoveryTlsKeyFile := flag.String("service-discovery-tls-key-file", "", "Path to TLS key file for service discovery")
	serviceDiscoveryTlsCaFile := flag.String("service-discovery-tls-ca-file", "", "Path to TLS CA file for service discovery")

	rateLimitQPS := flag.Int("rate-limit-qps", 0, "Rate limit QPS")
	rateLimitBurst := flag.Int("rate-limit-burst", 0, "Rate limit burst")
	apiKeys := flag.String("api-keys", "", "Comma-separated API keys")
	authApis := flag.String("auth-apis", "", "Comma-separated Auth APIs")
	logLevel := flag.String("log-level", "", "Set log level: error, debug, warn, info")
	healthCheckFlag := flag.String("health-check", "false", "Run health check and exit")
	enableCorsFlag := flag.String("enable-cors", "true", "Enable cross-domain resources")

	flag.Parse()

	// Read environment variables
	if v, exists := os.LookupEnv("CONFIG"); exists {
		*configFile = v
	}
	if v, exists := os.LookupEnv("GO_CAPTCHA_CONFIG"); exists {
		*gocaptchaConfigFile = v
	}

	if v, exists := os.LookupEnv("SERVICE_NAME"); exists {
		*serviceName = v
	}
	if v, exists := os.LookupEnv("SERVICE_NODE"); exists {
		n, _ := strconv.ParseInt(v, 10, 64)
		*serviceNode = n
	}
	if v, exists := os.LookupEnv("HTTP_PORT"); exists {
		*httpPort = v
	}
	if v, exists := os.LookupEnv("GRPC_PORT"); exists {
		*grpcPort = v
	}
	if v, exists := os.LookupEnv("API_KEYS"); exists {
		*apiKeys = v
	}
	if v, exists := os.LookupEnv("AUTH_APIS"); exists {
		*authApis = v
	}
	if v, exists := os.LookupEnv("CACHE_TYPE"); exists {
		*cacheType = v
	}
	if v, exists := os.LookupEnv("CACHE_ADDRS"); exists {
		*cacheAddrs = v
	}
	if v, exists := os.LookupEnv("CACHE_USERNAME"); exists {
		*cacheUsername = v
	}
	if v, exists := os.LookupEnv("CACHE_PASSWORD"); exists {
		*cachePassword = v
	}
	if v, exists := os.LookupEnv("CACHE_DB"); exists {
		*cacheDB = v
	}
	if v, exists := os.LookupEnv("LOG_LEVEL"); exists {
		*logLevel = v
	}

	if v, exists := os.LookupEnv("ENABLE_CORS"); exists {
		*enableCorsFlag = v
	}

	if v, exists := os.LookupEnv("ENABLE_DYNAMIC_CONFIG"); exists {
		*enableDynamicConfig = v
	}
	if v, exists := os.LookupEnv("DYNAMIC_CONFIG_TYPE"); exists {
		*dynamicConfigType = v
	}
	if v, exists := os.LookupEnv("DYNAMIC_CONFIG_ADDRS"); exists {
		*dynamicConfigAddrs = v
	}
	if v, exists := os.LookupEnv("DYNAMIC_CONFIG_USERNAME"); exists {
		*dynamicConfigUsername = v
	}
	if v, exists := os.LookupEnv("DYNAMIC_CONFIG_PASSWORD"); exists {
		*dynamicConfigPassword = v
	}

	if v, exists := os.LookupEnv("ENABLE_SERVICE_DISCOVERY"); exists {
		*enableServiceDiscovery = v
	}
	if v, exists := os.LookupEnv("SERVICE_DISCOVERY_TYPE"); exists {
		*serviceDiscoveryType = v
	}
	if v, exists := os.LookupEnv("SERVICE_DISCOVERY_ADDRS"); exists {
		*serviceDiscoveryAddrs = v
	}
	if v, exists := os.LookupEnv("SERVICE_DISCOVERY_USERNAME"); exists {
		*serviceDiscoveryUsername = v
	}
	if v, exists := os.LookupEnv("SERVICE_DISCOVERY_PASSWORD"); exists {
		*serviceDiscoveryPassword = v
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}
	setupLoggerLevel(logger, *logLevel)

	// Load configuration
	dc, err := config.NewDynamicConfig(*configFile, true)
	if err != nil {
		if helper.FileExists(*configFile) {
			logger.Warn("[App] Failed to load of the config.json file, using defaults", zap.Error(err))
		} else {
			logger.Warn("[App] No configuration file 'config.json' was provided for the application. Use the defaults configuration")
		}
		dc = config.DefaultDynamicConfig()
	}
	// Register hot update callback
	dc.RegisterHotCallback("UPDATE_LOG_LEVEL", func(dnCfg *config.DynamicConfig, hotType config.HotCallbackType) {
		setupLoggerLevel(logger, dnCfg.Get().LogLevel)
	})

	// Load configuration
	dgc, err := config2.NewDynamicConfig(*gocaptchaConfigFile, true)
	if err != nil {
		if helper.FileExists(*gocaptchaConfigFile) {
			logger.Warn("[App] Failed to load of the gocaptcha.json file, using defaults", zap.Error(err))
		} else {
			logger.Warn("[App] No configuration file 'gocaptcha.json' was provided for the application. Use the defaults configuration")
		}
		dgc = config2.DefaultDynamicConfig()
	}

	// Merge command-line flags
	cfg := dc.Get()
	cfg = config.MergeWithFlags(cfg, map[string]interface{}{
		"service-name": *serviceName,
		"service-node": *serviceNode,
		"http-port":    *httpPort,
		"grpc-port":    *grpcPort,

		"cache-type":       *cacheType,
		"cache-addrs":      *cacheAddrs,
		"cache-username":   *cacheUsername,
		"cache-password":   *cachePassword,
		"cache-db":         *cacheDB,
		"cache-ttl":        *cacheTTL,
		"cache-key-prefix": *cacheKeyPrefix,

		"enable-dynamic-config":           *enableDynamicConfig,
		"dynamic-config-type":             *dynamicConfigType,
		"dynamic-config-addrs":            *dynamicConfigAddrs,
		"dynamic-config-username":         *dynamicConfigUsername,
		"dynamic-config-password":         *dynamicConfigPassword,
		"dynamic-config-ttl":              *dynamicConfigTTL,
		"dynamic-config-keep-alive":       *dynamicConfigKeepAlive,
		"dynamic-config-max-retries":      *dynamicConfigMaxRetries,
		"dynamic-config-base-retry-delay": *dynamicConfigBaseRetryDelay,
		"dynamic-config-tls-server-name":  *dynamicConfigTlsServerName,
		"dynamic-config-tls-address":      *dynamicConfigTlsAddress,
		"dynamic-config-tls-cert-file":    *dynamicConfigTlsCertFile,
		"dynamic-config-tls-key-file":     *dynamicConfigTlsKeyFile,
		"dynamic-config-tls-ca-file":      *dynamicConfigTlsCaFile,

		"enable-service-discovery":           *enableServiceDiscovery,
		"service-discovery-type":             *serviceDiscoveryType,
		"service-discovery-addrs":            *serviceDiscoveryAddrs,
		"service-discovery-username":         *serviceDiscoveryUsername,
		"service-discovery-password":         *serviceDiscoveryPassword,
		"service-discovery-ttl":              *serviceDiscoveryTTL,
		"service-discovery-keep-alive":       *serviceDiscoveryKeepAlive,
		"service-discovery-max-retries":      *serviceDiscoveryMaxRetries,
		"service-discovery-base-retry-delay": *serviceDiscoveryBaseRetryDelay,
		"service-discovery-tls-server-name":  *serviceDiscoveryTlsServerName,
		"service-discovery-tls-address":      *serviceDiscoveryTlsAddress,
		"service-discovery-tls-cert-file":    *serviceDiscoveryTlsCertFile,
		"service-discovery-tls-key-file":     *serviceDiscoveryTlsKeyFile,
		"service-discovery-tls-ca-file":      *serviceDiscoveryTlsCaFile,

		"rate-limit-qps":   *rateLimitQPS,
		"rate-limit-burst": *rateLimitBurst,
		"enable-cors":      *enableCorsFlag,
		"api-keys":         *apiKeys,
		"auth-apis":        *authApis,
	})
	if err = dc.Update(cfg); err != nil {
		logger.Fatal("[App] Configuration validation failed", zap.Error(err))
	}

	// Initialize rate limiter
	limiter := middleware.NewDynamicLimiter(cfg.RateLimitQPS, cfg.RateLimitBurst)
	dc.RegisterHotCallback("UPDATE_LIMITER", func(dnCfg *config.DynamicConfig, hotType config.HotCallbackType) {
		limiter.Update(dnCfg.Get().RateLimitQPS, dnCfg.Get().RateLimitBurst)
	})

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

	// Setup cache
	cacheMgr, err := setupCacheManager(dc, logger)
	if err != nil {
		logger.Fatal("[App] Create cache manager", zap.Error(err))
	}

	// Setup service discovery
	discovery, err := setupServiceDiscovery(dc, logger)
	if err != nil {
		logger.Fatal("[App] Setup service discovery", zap.Error(err))
	}

	// Setup dynamic config
	configManager, err := setupDynamicConfig(dc, dgc, logger)
	if err != nil {
		logger.Fatal("[App] Setup dynamic config manager", zap.Error(err))
	}

	// Setup captcha
	captcha, err := gocaptcha.Setup(dgc)
	if err != nil {
		logger.Fatal("[App] Failed to setup gocaptcha: ", zap.Error(err))
	}
	dgc.RegisterHotCallback("GENERATE_CAPTCHA", func(captchaConfig *config2.DynamicCaptchaConfig, callbackType config2.HotCallbackType) {
		err = captcha.HotSetup(captchaConfig)
		if err != nil {
			logger.Error("[App] Failed to hot update gocaptcha, without any change: ", zap.Error(err))
		}
	})

	// Perform health check if requested
	if *healthCheckFlag == "true" {
		if err = setupHealthCheck(":"+cfg.HTTPPort, ":"+cfg.GRPCPort); err != nil {
			logger.Error("[App] Filed to health check", zap.Error(err))
			os.Exit(1)
		}
		os.Exit(0)
	}

	return &App{
		logger:         logger,
		dynamicCfg:     dc,
		dynamicCaptCfg: dgc,
		cacheMgr:       cacheMgr,
		discovery:      discovery,
		configManager:  configManager,
		cacheBreaker:   cacheBreaker,
		limiter:        limiter,
		captcha:        captcha,
	}, nil
}

// Start starting the Application
func (a *App) Start(ctx context.Context) error {
	cfg := a.dynamicCfg.Get()

	// Register service with discovery
	err := a.startDiscoveryRegister(ctx, &cfg)
	if err != nil {
		return err
	}

	// Service context
	svcCtx := common.NewSvcContext()
	svcCtx.CacheMgr = a.cacheMgr
	svcCtx.DynamicConfig = a.dynamicCfg
	svcCtx.Logger = a.logger
	svcCtx.Captcha = a.captcha

	// Start HTTP server
	if cfg.HTTPPort != "" && cfg.HTTPPort != "0" {
		if err = a.startHTTPServer(svcCtx, &cfg); err != nil {
			return err
		}
	}

	// Start gRPC server
	if cfg.GRPCPort != "" && cfg.GRPCPort != "0" {
		if err = a.startGRPCServer(svcCtx, &cfg); err != nil {
			return err
		}
	}

	return nil
}

// startDiscoveryRegister start service discovery register
func (a *App) startDiscoveryRegister(ctx context.Context, cfg *config.Config) error {
	var instanceID string
	if a.discovery != nil {
		instanceID = uuid.New().String()
		if err := a.discovery.Register(ctx, cfg.ServiceName, instanceID, "localhost", cfg.HTTPPort, cfg.GRPCPort); err != nil {
			return fmt.Errorf("failed to register service: %v", err)
		}
		go a.watchServiceDiscoveryInstances(ctx, instanceID)
	}
	return nil
}

// startHTTPServer start HTTP server
func (a *App) startHTTPServer(svcCtx *common.SvcContext, cfg *config.Config) error {
	handlers := server.NewHTTPHandlers(svcCtx)
	var middlewares = make([]middleware.HTTPMiddleware, 0)

	// Enable cross-domain resource
	//if cfg.EnableCors {
	//	middlewares = append(middlewares, nil, middleware.CORSMiddleware(a.dynamicCfg, a.logger))
	//}

	middlewares = append(middlewares,
		middleware.CORSMiddleware(a.dynamicCfg, a.logger),
		middleware.APIKeyMiddleware(a.dynamicCfg, a.logger),
		middleware.LoggingMiddleware(a.logger),
		middleware.RateLimitMiddleware(a.limiter, a.logger),
		middleware.CircuitBreakerMiddleware(a.cacheBreaker, a.logger),
	)

	mwChain := middleware.NewChainHTTP(middlewares...)

	http.Handle("/status/health", mwChain.Then(handlers.HealthStatusHandler))
	http.Handle("/rate-limit", mwChain.Then(middleware.RateLimitHandler(a.limiter, a.logger)))

	http.Handle("/api/v1/public/get-data", mwChain.Then(handlers.GetDataHandler))
	http.Handle("/api/v1/public/check-data", mwChain.Then(handlers.CheckDataHandler))
	http.Handle("/api/v1/public/check-status", mwChain.Then(handlers.CheckStatusHandler))

	http.Handle("/api/v1/manage/get-status-info", mwChain.Then(handlers.GetStatusInfoHandler))
	http.Handle("/api/v1/manage/del-status-info", mwChain.Then(handlers.DelStatusInfoHandler))
	http.Handle("/api/v1/manage/upload-resource", mwChain.Then(handlers.UploadResourceHandler))
	http.Handle("/api/v1/manage/delete-resource", mwChain.Then(handlers.DeleteResourceHandler))
	http.Handle("/api/v1/manage/get-resource-list", mwChain.Then(handlers.GetResourceListHandler))
	http.Handle("/api/v1/manage/get-config", mwChain.Then(handlers.GetGoCaptchaConfigHandler))
	http.Handle("/api/v1/manage/update-hot-config", mwChain.Then(handlers.UpdateHotGoCaptchaConfigHandler))

	a.httpServer = &http.Server{
		Addr: ":" + cfg.HTTPPort,
	}

	go func() {
		a.logger.Info("[App] Starting HTTP server", zap.String("port", cfg.HTTPPort))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("[App] HTTP server failed", zap.Error(err))
		}
	}()

	return nil
}

// startGRPCServer start gRPC server
func (a *App) startGRPCServer(svcCtx *common.SvcContext, cfg *config.Config) error {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	interceptor := middleware.UnaryServerInterceptor(a.dynamicCfg, a.logger, a.cacheBreaker)
	a.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	proto.RegisterGoCaptchaServiceServer(a.grpcServer, server.NewGoCaptchaServer(svcCtx))

	go func() {
		a.logger.Info("[App] Starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := a.grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			a.logger.Fatal("[App] gRPC server failed", zap.Error(err))
		}
	}()
	return nil
}

// watchServiceDiscoveryInstances periodically updates service instances
func (a *App) watchServiceDiscoveryInstances(ctx context.Context, instanceID string) {
	if a.discovery == nil {
		return
	}

	cfg := a.dynamicCfg.Get()
	ch, err := a.discovery.Watch(ctx, cfg.ServiceName)
	if err != nil {
		a.logger.Fatal("[App] Failed to service discovery watch", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			if a.discovery != nil {
				if err = a.discovery.Deregister(ctx, cfg.ServiceName, instanceID); err != nil {
					a.logger.Error("[App] Failed to deregister service", zap.Error(err))
				}
			}
			return
		case instances, ok := <-ch:
			if !ok {
				return
			}
			a.logger.Info("[App] Discovered service instances", zap.Int("count", len(instances)))
		}
	}
}

// Shutdown gracefully stops the application
func (a *App) Shutdown() {
	a.logger.Info("[App] Received shutdown signal, shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defer a.logger.Sync()

	// Stop HTTP server
	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			a.logger.Error("[App] HTTP server shutdown error", zap.Error(err))
		} else {
			a.logger.Info("[App] HTTP server shut down successfully")
		}
	}

	// Stop gRPC server
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
		a.logger.Info("[App] gRPC server shut down successfully")
	}

	// Stop cache
	err := a.cacheMgr.Close()
	if err != nil {
		a.logger.Error("[App] CacheManager client close error", zap.Error(err))
	} else {
		a.logger.Info("[App] CacheManager client stopped successfully", zap.Error(err))
	}

	// Stop service discovery
	if a.discovery != nil {
		if err := a.discovery.Close(); err != nil {
			a.logger.Error("[App] Service discovery close error", zap.Error(err))
		} else {
			a.logger.Info("[App] Service discovery closed successfully")
		}
	}

	// Stop config manager
	if a.configManager != nil {
		if err = a.configManager.Close(); err != nil {
			a.logger.Error("[App] Config manager close error", zap.Error(err))
		} else {
			a.logger.Info("[App] Config manager closed successfully")
		}
	}

	a.logger.Info("[App] App service shutdown")
}
