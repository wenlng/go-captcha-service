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

	"github.com/memcachier/mc/v3"
)

// MemcacheClient implements the Cache interface for Memcached
type MemcacheClient struct {
	client *mc.Client
	prefix string
	ttl    time.Duration
}

// NewMemcacheClient ..
func NewMemcacheClient(addrs, prefix string, ttl time.Duration, username, password string) (*MemcacheClient, error) {
	client := mc.NewMC(addrs, username, password)
	return &MemcacheClient{client: client, prefix: prefix, ttl: ttl}, nil
}

// GetCache retrieves a value from Memcached
func (c *MemcacheClient) GetCache(ctx context.Context, key string) (string, error) {
	key = c.prefix + key
	item, _, _, err := c.client.Get(key)
	if err == mc.ErrNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return item, nil
}

// SetCache stores a value in Memcached
func (c *MemcacheClient) SetCache(ctx context.Context, key, value string) error {
	key = c.prefix + key
	_, err := c.client.Set(key, value, uint32(0), uint32(c.ttl/time.Second), uint64(0))
	return err
}

// DeleteCache ..
func (c *MemcacheClient) DeleteCache(ctx context.Context, key string) error {
	key = c.prefix + key
	err := c.client.Del(key)
	if err != nil && err != mc.ErrNotFound {
		return fmt.Errorf("memcache delete error: %v", err)
	}
	return nil
}

// Close ..
func (c *MemcacheClient) Close() error {
	return nil
}
