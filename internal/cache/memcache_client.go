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
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	key = c.prefix + key
	resultChan := make(chan struct {
		item string
		err  error
	}, 1)

	go func() {
		item, _, _, err := c.client.Get(key)
		if err == mc.ErrNotFound {
			resultChan <- struct {
				item string
				err  error
			}{item: "", err: nil}
			return
		}
		if err != nil {
			resultChan <- struct {
				item string
				err  error
			}{item: "", err: err}
			return
		}
		resultChan <- struct {
			item string
			err  error
		}{item: item, err: nil}
	}()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("memcache get timeout: %v", ctx.Err())
	case result := <-resultChan:
		if result.err != nil {
			return "", fmt.Errorf("memcache get error: %v", result.err)
		}
		return result.item, nil
	}
}

// SetCache stores a value in Memcached
func (c *MemcacheClient) SetCache(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	key = c.prefix + key
	errChan := make(chan error, 1)

	go func() {
		_, err := c.client.Set(key, value, uint32(0), uint32(c.ttl/time.Second), uint64(0))
		errChan <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("memcache set timeout: %v", ctx.Err())
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("memcache set error: %v", err)
		}
		return nil
	}
}

// DeleteCache ..
func (c *MemcacheClient) DeleteCache(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	key = c.prefix + key
	errChan := make(chan error, 1)

	go func() {
		err := c.client.Del(key)
		errChan <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("memcache delete timeout: %v", ctx.Err())
	case err := <-errChan:
		if err != nil && err != mc.ErrNotFound {
			return fmt.Errorf("memcache delete error: %v", err)
		}
		return nil
	}
}

// Close ..
func (c *MemcacheClient) Close() error {
	return nil
}
