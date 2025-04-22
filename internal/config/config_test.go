package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	configContent := `
{
	"http_port": "8080",
	"grpc_port": "50051",
	"redis_addrs": "localhost:6379",
	"etcd_addrs": "localhost:2379",
	"memcache_addrs": "localhost:11211",
	"cache_type": "redis",
	"cache_ttl": 60,
	"cache_cleanup_interval": 10,
	"service_discovery": "etcd",
	"service_discovery_addrs": "localhost:2379",
	"rate_limit_qps": 1000,
	"rate_limit_burst": 1000,
	"load_balancer": "round-robin",
	"api_keys": ["key1", "key2"]
}`
	tmpFile, err := os.CreateTemp("", "Config.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.WriteString(configContent)
	assert.NoError(t, err)
	tmpFile.Close()

	config, err := Load(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "8080", config.HTTPPort)
	assert.Equal(t, "redis", config.CacheType)
	assert.Equal(t, []string{"key1", "key2"}, config.APIKeys)
}

func TestValidate(t *testing.T) {
	config := Config{
		HTTPPort:              "8080",
		GRPCPort:              "50051",
		RedisAddrs:            "localhost:6379",
		CacheType:             "redis",
		CacheTTL:              60,
		ServiceDiscovery:      "etcd",
		ServiceDiscoveryAddrs: "localhost:2379",
		RateLimitQPS:          1000,
		RateLimitBurst:        1000,
		APIKeys:               []string{"key1"},
	}
	assert.NoError(t, Validate(config))

	config.APIKeys = nil
	assert.Error(t, Validate(config))

	config.APIKeys = []string{""}
	assert.Error(t, Validate(config))
}

func TestDynamicConfig(t *testing.T) {
	configContent := `
{
	"http_port": "8080",
	"grpc_port": "50051",
	"redis_addrs": "localhost:6379",
	"etcd_addrs": "localhost:2379",
	"memcache_addrs": "localhost:11211",
	"cache_type": "redis",
	"cache_ttl": 60,
	"cache_cleanup_interval": 10,
	"service_discovery": "etcd",
	"service_discovery_addrs": "localhost:2379",
	"rate_limit_qps": 1000,
	"rate_limit_burst": 1000,
	"load_balancer": "round-robin",
	"api_keys": ["key1"]
}`
	tmpDir, err := os.MkdirTemp("", "config_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "Config.json")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	dc, err := NewDynamicConfig(configPath, false)
	assert.NoError(t, err)

	cfg := dc.Get()
	assert.Equal(t, []string{"key1"}, cfg.APIKeys)
	assert.Equal(t, 1000, cfg.RateLimitQPS)

	// Update Config file
	newContent := `
{
	"http_port": "8080",
	"grpc_port": "50051",
	"redis_addrs": "localhost:6379",
	"etcd_addrs": "localhost:2379",
	"memcache_addrs": "localhost:11211",
	"cache_type": "redis",
	"cache_ttl": 60,
	"cache_cleanup_interval": 10,
	"service_discovery": "etcd",
	"service_discovery_addrs": "localhost:2379",
	"rate_limit_qps": 2000,
	"rate_limit_burst": 2000,
	"load_balancer": "round-robin",
	"api_keys": ["key2", "key3"]
}`
	err = os.WriteFile(configPath, []byte(newContent), 0644)
	assert.NoError(t, err)

	// Wait for reload
	time.Sleep(100 * time.Millisecond)

	cfg = dc.Get()
	assert.Equal(t, []string{"key2", "key3"}, cfg.APIKeys)
	assert.Equal(t, 2000, cfg.RateLimitQPS)
}
