/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type CacheType string

// CacheType .
const (
	CacheTypeRedis    CacheType = "redis"
	CacheTypeMemory             = "memory"
	CacheTypeEtcd               = "etcd"
	CacheTypeMemcache           = "memcache"
)

// Cache defines the interface for cache operations
type Cache interface {
	GetCache(ctx context.Context, key string) (string, error)
	SetCache(ctx context.Context, key, value string) error
	DeleteCache(ctx context.Context, key string) error
	Close() error
}

// CaptCacheData ..
type CaptCacheData struct {
	Data   interface{} `json:"data"`
	Type   int         `json:"type"`
	Status int         `json:"status"`
}

// CacheManager ..
type CacheManager struct {
	cache      Cache
	mu         sync.RWMutex
	cType      CacheType
	cAddress   string
	cUsername  string
	cPassword  string
	cKeyPrefix string
	cTtl       time.Duration
	cCleanInt  time.Duration
}

// CacheMgrParams ..
type CacheMgrParams struct {
	Type          CacheType
	CacheAddrs    string
	CacheUsername string
	CachePassword string
	CacheDB       string
	KeyPrefix     string
	Ttl           time.Duration
	CleanInt      time.Duration
}

// NewCacheManager ..
func NewCacheManager(arg *CacheMgrParams) (*CacheManager, error) {
	cm := &CacheManager{}
	err := cm.Setup(arg)
	return cm, err
}

// GetCache ..
func (cm *CacheManager) GetCache() Cache {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cache
}

// Setup initialize the cache
func (cm *CacheManager) Setup(arg *CacheMgrParams) error {
	var curCache Cache
	var err error
	curAddrs := arg.CacheAddrs

	switch arg.Type {
	case CacheTypeRedis:
		if cm.cAddress == curAddrs &&
			cm.cKeyPrefix == arg.KeyPrefix &&
			cm.cTtl == arg.Ttl &&
			cm.cUsername == arg.CacheUsername &&
			cm.cPassword == arg.CachePassword {
			return nil
		}
		curCache, err = NewRedisClient(arg.CacheAddrs, arg.KeyPrefix, arg.Ttl, arg.CacheUsername, arg.CachePassword, arg.CacheDB)
		if err != nil {
			return fmt.Errorf("failed to initialize Redis: %v", err)
		}
	case CacheTypeMemory:
		if cm.cKeyPrefix == arg.KeyPrefix &&
			cm.cTtl == arg.Ttl &&
			cm.cCleanInt == arg.CleanInt &&
			cm.cUsername == arg.CacheUsername &&
			cm.cPassword == arg.CachePassword {
			return nil
		}
		curCache = NewMemoryCache(arg.KeyPrefix, arg.Ttl, arg.CleanInt)
	case CacheTypeEtcd:
		if cm.cAddress == curAddrs &&
			cm.cKeyPrefix == arg.KeyPrefix &&
			cm.cTtl == arg.Ttl &&
			cm.cUsername == arg.CacheUsername &&
			cm.cPassword == arg.CachePassword {
			return nil
		}
		curCache, err = NewEtcdClient(arg.CacheAddrs, arg.KeyPrefix, arg.Ttl, arg.CacheUsername, arg.CachePassword)
		if err != nil {
			return fmt.Errorf("failed to initialize Etcd: %v", err)
		}
	case CacheTypeMemcache:
		if cm.cAddress == curAddrs &&
			cm.cKeyPrefix == arg.KeyPrefix &&
			cm.cTtl == arg.Ttl &&
			cm.cUsername == arg.CacheUsername &&
			cm.cPassword == arg.CachePassword {
			return nil
		}
		curCache, err = NewMemcacheClient(arg.CacheAddrs, arg.KeyPrefix, arg.Ttl, arg.CacheUsername, arg.CachePassword)
		if err != nil {
			return fmt.Errorf("failed to initialize Memcached: %v", err)
		}
	default:
		return fmt.Errorf("invalid cache type: %v", arg.Type)
	}

	cm.cType = arg.Type
	cm.cAddress = curAddrs
	cm.cKeyPrefix = arg.KeyPrefix
	cm.cTtl = arg.Ttl
	cm.cCleanInt = arg.CleanInt
	cm.cUsername = arg.CacheUsername
	cm.cPassword = arg.CacheAddrs

	cm.mu.Lock()
	cm.cache = curCache
	cm.mu.Unlock()

	return nil
}

// Close ..
func (cm *CacheManager) Close() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if redisClient, ok := cm.cache.(*RedisClient); ok {
		if err := redisClient.Close(); err != nil {
			return fmt.Errorf("redis client close error: %v", err)
		}
		return nil
	}

	if memoryCache, ok := cm.cache.(*MemoryCache); ok {
		memoryCache.Stop()
		return nil
	}

	if etcdClient, ok := cm.cache.(*EtcdClient); ok {
		if err := etcdClient.Close(); err != nil {
			return fmt.Errorf("etcd client close error: %v", err)
		}
		return nil
	}

	if memcacheClient, ok := cm.cache.(*MemcacheClient); ok {
		if err := memcacheClient.Close(); err != nil {
			return fmt.Errorf("memcached client close error: %v", err)
		}
	}

	return nil
}
