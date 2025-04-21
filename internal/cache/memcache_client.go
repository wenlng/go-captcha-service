/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcacheClient implements the Cache interface for Memcached
type MemcacheClient struct {
	client *memcache.Client
	prefix string
	ttl    time.Duration
}

// NewMemcacheClient ..
func NewMemcacheClient(addrs, prefix string, ttl time.Duration) (*MemcacheClient, error) {
	client := memcache.New(addrs)
	return &MemcacheClient{client: client, prefix: prefix, ttl: ttl}, nil
}

// GetCache retrieves a value from Memcached
func (c *MemcacheClient) GetCache(ctx context.Context, key string) (string, error) {
	key = c.prefix + key
	item, err := c.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

// SetCache stores a value in Memcached
func (c *MemcacheClient) SetCache(ctx context.Context, key, value string) error {
	key = c.prefix + key
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(c.ttl / time.Second),
	}
	return c.client.Set(item)
}

// DeleteCache ..
func (c *MemcacheClient) DeleteCache(ctx context.Context, key string) error {
	key = c.prefix + key
	err := c.client.Delete(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("memcache delete error: %v", err)
	}
	return nil
}

// Close ..
func (c *MemcacheClient) Close() error {
	return nil
}
