package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// Config defines the configuration structure for the application
type Config struct {
	ServiceName           string   `json:"service_name"`
	HTTPPort              string   `json:"http_port"`
	GRPCPort              string   `json:"grpc_port"`
	RedisAddrs            string   `json:"redis_addrs"`
	EtcdAddrs             string   `json:"etcd_addrs"`
	MemcacheAddrs         string   `json:"memcache_addrs"`
	CacheType             string   `json:"cache_type"`             // redis, memory, etcd, memcache
	CacheTTL              int      `json:"cache_ttl"`              // seconds
	CacheCleanupInt       int      `json:"cache_cleanup_interval"` // seconds
	CacheKeyPrefix        string   `json:"cache_key_prefix"`
	ServiceDiscovery      string   `json:"service_discovery"` // etcd, zookeeper, consul, nacos
	ServiceDiscoveryAddrs string   `json:"service_discovery_addrs"`
	RateLimitQPS          int      `json:"rate_limit_qps"`
	RateLimitBurst        int      `json:"rate_limit_burst"`
	LoadBalancer          string   `json:"load_balancer"` // round-robin, consistent-hash
	APIKeys               []string `json:"api_keys"`      // API keys for authentication
	EnableCors            bool     `json:"enable_cors"`   // cross-domain resources
}

// DynamicConfig .
type DynamicConfig struct {
	Config Config
	mu     sync.RWMutex
}

// NewDynamicConfig .
func NewDynamicConfig(file string) (*DynamicConfig, error) {
	cfg, err := Load(file)
	if err != nil {
		return nil, err
	}
	dc := &DynamicConfig{Config: cfg}
	go dc.watchFile(file)
	return dc, nil
}

// Get retrieves the current configuration
func (dc *DynamicConfig) Get() Config {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.Config
}

// Update updates the configuration
func (dc *DynamicConfig) Update(cfg Config) error {
	if err := Validate(cfg); err != nil {
		return err
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.Config = cfg
	return nil
}

// watchFile monitors the Config file for changes
func (dc *DynamicConfig) watchFile(file string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	absPath, err := filepath.Abs(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get absolute path: %v\n", err)
		return
	}
	dir := filepath.Dir(absPath)

	if err := watcher.Add(dir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to watch directory: %v\n", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Name == absPath && (event.Op&fsnotify.Write == fsnotify.Write) {
				cfg, err := Load(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to reload Config: %v\n", err)
					continue
				}
				if err := dc.Update(cfg); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to update Config: %v\n", err)
					continue
				}
				fmt.Printf("Configuration reloaded successfully\n")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
		}
	}
}

// Load reads the configuration from a file
func Load(file string) (Config, error) {
	var config Config
	data, err := os.ReadFile(file)
	if err != nil {
		return config, fmt.Errorf("failed to read Config file: %v", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse Config file: %v", err)
	}
	return config, nil
}

// Validate checks the configuration for validity
func Validate(config Config) error {
	if !isValidPort(config.HTTPPort) {
		return fmt.Errorf("invalid http_port: %s", config.HTTPPort)
	}
	if !isValidPort(config.GRPCPort) {
		return fmt.Errorf("invalid grpc_port: %s", config.GRPCPort)
	}

	validCacheTypes := map[string]bool{
		"redis":    true,
		"memory":   true,
		"etcd":     true,
		"memcache": true,
	}
	if !validCacheTypes[config.CacheType] {
		return fmt.Errorf("invalid cache_type: %s, must be redis, memory, etcd, or memcache", config.CacheType)
	}

	switch config.CacheType {
	case "redis":
		if !isValidAddrs(config.RedisAddrs) {
			return fmt.Errorf("invalid redis_addrs: %s", config.RedisAddrs)
		}
	case "etcd":
		if !isValidAddrs(config.EtcdAddrs) {
			return fmt.Errorf("invalid etcd_addrs: %s", config.EtcdAddrs)
		}
	case "memcache":
		if !isValidAddrs(config.MemcacheAddrs) {
			return fmt.Errorf("invalid memcache_addrs: %s", config.MemcacheAddrs)
		}
	}

	if config.CacheTTL <= 0 {
		return fmt.Errorf("cache_ttl must be positive: %d", config.CacheTTL)
	}
	if config.CacheCleanupInt <= 0 && config.CacheType == "memory" {
		return fmt.Errorf("cache_cleanup_interval must be positive for memory cache: %d", config.CacheCleanupInt)
	}

	validDiscoveryTypes := map[string]bool{
		"etcd":      true,
		"zookeeper": true,
		"consul":    true,
		"nacos":     true,
	}
	if config.ServiceDiscovery != "" && !validDiscoveryTypes[config.ServiceDiscovery] {
		return fmt.Errorf("invalid service_discovery: %s, must be etcd, zookeeper, consul, or nacos", config.ServiceDiscovery)
	}
	if config.ServiceDiscovery != "" && !isValidAddrs(config.ServiceDiscoveryAddrs) {
		return fmt.Errorf("invalid service_discovery_addrs: %s", config.ServiceDiscoveryAddrs)
	}

	if config.RateLimitQPS <= 0 {
		return fmt.Errorf("rate_limit_qps must be positive: %d", config.RateLimitQPS)
	}
	if config.RateLimitBurst <= 0 {
		return fmt.Errorf("rate_limit_burst must be positive: %d", config.RateLimitBurst)
	}

	if len(config.APIKeys) == 0 {
		return fmt.Errorf("api_keys must not be empty")
	}
	for _, key := range config.APIKeys {
		if key == "" {
			return fmt.Errorf("api_keys contain empty key")
		}
	}

	validBalancerTypes := map[string]bool{
		"round-robin":     true,
		"consistent-hash": true,
	}
	if config.LoadBalancer != "" && !validBalancerTypes[config.LoadBalancer] {
		return fmt.Errorf("invalid load_balancer: %s, must be round-robin or consistent-hash", config.LoadBalancer)
	}

	return nil
}

// isValidPort checks if a port number is valid
func isValidPort(port string) bool {
	p, err := strconv.Atoi(port)
	return err == nil && p > 0 && p <= 65535
}

// isValidAddrs checks if addresses are valid
func isValidAddrs(addrs string) bool {
	if addrs == "" {
		return false
	}
	addrRegex := regexp.MustCompile(`^([a-zA-Z0-9.-]+:[0-9]+)(,[a-zA-Z0-9.-]+:[0-9]+)*$`)
	return addrRegex.MatchString(addrs)
}

// MergeWithFlags merges command-line flags into the configuration
func MergeWithFlags(config Config, flags map[string]interface{}) Config {
	if v, ok := flags["http-port"].(string); ok && v != "" {
		config.HTTPPort = v
	}
	if v, ok := flags["service-name"].(string); ok && v != "" {
		config.ServiceName = v
	}
	if v, ok := flags["grpc-port"].(string); ok && v != "" {
		config.GRPCPort = v
	}
	if v, ok := flags["redis-addrs"].(string); ok && v != "" {
		config.RedisAddrs = v
	}
	if v, ok := flags["etcd-addrs"].(string); ok && v != "" {
		config.EtcdAddrs = v
	}
	if v, ok := flags["memcache-addrs"].(string); ok && v != "" {
		config.MemcacheAddrs = v
	}
	if v, ok := flags["cache-type"].(string); ok && v != "" {
		config.CacheType = v
	}
	if v, ok := flags["cache-ttl"].(int); ok && v != 0 {
		config.CacheTTL = v
	}
	if v, ok := flags["cache-cleanup-interval"].(int); ok && v != 0 {
		config.CacheCleanupInt = v
	}
	if v, ok := flags["cache-key-prefix"].(string); ok && v != "" {
		config.CacheKeyPrefix = v
	}
	if v, ok := flags["service-discovery"].(string); ok && v != "" {
		config.ServiceDiscovery = v
	}
	if v, ok := flags["service-discovery-addrs"].(string); ok && v != "" {
		config.ServiceDiscoveryAddrs = v
	}
	if v, ok := flags["rate-limit-qps"].(int); ok && v != 0 {
		config.RateLimitQPS = v
	}
	if v, ok := flags["rate-limit-burst"].(int); ok && v != 0 {
		config.RateLimitBurst = v
	}
	if v, ok := flags["load-balancer"].(string); ok && v != "" {
		config.LoadBalancer = v
	}
	if v, ok := flags["api-keys"].(string); ok && v != "" {
		config.APIKeys = strings.Split(v, ",")
	}
	if v, ok := flags["enable-cors"].(bool); ok {
		config.EnableCors = v
	}
	return config
}

func DefaultConfig() Config {
	return Config{
		ServiceName:           "go-captcha-service",
		HTTPPort:              "8080",
		GRPCPort:              "50051",
		RedisAddrs:            "localhost:6379",
		EtcdAddrs:             "localhost:2379",
		MemcacheAddrs:         "localhost:11211",
		CacheType:             "memory",
		CacheTTL:              60,
		CacheCleanupInt:       10,
		CacheKeyPrefix:        "GO_CAPTCHA_DATA",
		ServiceDiscovery:      "",
		ServiceDiscoveryAddrs: "localhost:2379",
		RateLimitQPS:          1000,
		RateLimitBurst:        1000,
		LoadBalancer:          "round-robin",
		APIKeys:               []string{"my-secret-key-123"},
		EnableCors:            false,
	}
}
