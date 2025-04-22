/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package cache

import (
	"context"
	"sync"
	"time"
)

// MemoryCache is an in-memory cache with TTL and cleanup
type MemoryCache struct {
	items    map[string]cacheItem
	mu       sync.RWMutex
	ttl      time.Duration
	prefix   string
	stop     chan struct{}
	cleanInt time.Duration
}

type cacheItem struct {
	value      string
	expiration int64
}

// NewMemoryCache ..
func NewMemoryCache(prefix string, ttl, cleanupInterval time.Duration) *MemoryCache {
	cache := &MemoryCache{
		items:    make(map[string]cacheItem),
		ttl:      ttl,
		prefix:   prefix,
		stop:     make(chan struct{}),
		cleanInt: cleanupInterval,
	}
	go cache.startCleanup()
	return cache
}

// startCleanup runs periodic cleanup of expired items
func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(c.cleanInt)
	defer ticker.Stop()
	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			c.mu.Lock()
			for key, item := range c.items {
				if item.expiration <= time.Now().UnixNano() {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

// GetCache retrieves a value from memory cache
func (c *MemoryCache) GetCache(ctx context.Context, key string) (string, error) {
	key = c.prefix + key
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.items[key]
	if !exists || item.expiration <= time.Now().UnixNano() {
		return "", nil
	}
	return item.value, nil
}

// SetCache stores a value in memory cache
func (c *MemoryCache) SetCache(ctx context.Context, key, value string) error {
	key = c.prefix + key
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl).UnixNano(),
	}
	return nil
}

// DeleteCache delete a value in memory cache
func (c *MemoryCache) DeleteCache(ctx context.Context, key string) error {
	key = c.prefix + key
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

// Close ..
func (c *MemoryCache) Close() error {
	c.Stop()
	return nil
}

// Stop ..
func (c *MemoryCache) Stop() {
	close(c.stop)
}
