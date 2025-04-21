/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// EtcdClient implements the Cache interface for etcd
type EtcdClient struct {
	client *clientv3.Client
	prefix string
	ttl    time.Duration
}

// NewEtcdClient ..
func NewEtcdClient(addrs, prefix string, ttl time.Duration) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{addrs},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdClient{client: client, prefix: prefix, ttl: ttl}, nil
}

// GetCache retrieves a value from etcd
func (c *EtcdClient) GetCache(ctx context.Context, key string) (string, error) {
	key = c.prefix + key
	resp, err := c.client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

// SetCache stores a value in etcd
func (c *EtcdClient) SetCache(ctx context.Context, key, value string) error {
	key = c.prefix + key
	session, err := concurrency.NewSession(c.client, concurrency.WithTTL(int(c.ttl/time.Second)))
	if err != nil {
		return fmt.Errorf("failed to create etcd session: %v", err)
	}
	defer session.Close()

	prefix := "http"
	if strings.Contains(key, ":grpc:") {
		prefix = "grpc"
	}
	mutex := concurrency.NewMutex(session, "/go-captcha-cache-lock/"+prefix)
	if err = mutex.Lock(ctx); err != nil {
		return fmt.Errorf("failed to acquire etcd lock: %v", err)
	}
	defer mutex.Unlock(ctx)

	_, err = c.client.Put(ctx, key, value, clientv3.WithLease(session.Lease()))
	if err != nil {
		return err
	}
	return nil
}

// DeleteCache ..
func (c *EtcdClient) DeleteCache(ctx context.Context, key string) error {
	key = c.prefix + key
	_, err := c.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("etcd delete error: %v", err)
	}
	return nil
}

// Close ..
func (c *EtcdClient) Close() error {
	return c.client.Close()
}
