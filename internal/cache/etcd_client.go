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
func NewEtcdClient(addrs, prefix string, ttl time.Duration, username, password string) (*EtcdClient, error) {
	endpoints := strings.Split(addrs, ",")
	for i := range endpoints {
		endpoints[i] = strings.TrimSpace(endpoints[i])
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:            endpoints,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    30 * time.Second,
		DialKeepAliveTimeout: 5 * time.Second,
		Username:             username,
		Password:             password,
		RejectOldCluster:     false,
		PermitWithoutStream:  true,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdClient{client: client, prefix: prefix, ttl: ttl}, nil
}

// GetCache retrieves a value from etcd
func (c *EtcdClient) GetCache(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	key = c.prefix + key
	resp, err := c.client.Get(ctx, key)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("etcd get timeout: %v", err)
		}
		return "", fmt.Errorf("etcd get error: %v", err)
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

// SetCache stores a value in etcd
func (c *EtcdClient) SetCache(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

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
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("etcd lock timeout: %v", err)
		}
		return fmt.Errorf("failed to acquire etcd lock: %v", err)
	}
	defer mutex.Unlock(ctx)

	_, err = c.client.Put(ctx, key, value, clientv3.WithLease(session.Lease()))
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("etcd put timeout: %v", err)
		}
		return fmt.Errorf("etcd put error: %v", err)
	}
	return nil
}

// DeleteCache ..
func (c *EtcdClient) DeleteCache(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	key = c.prefix + key
	_, err := c.client.Delete(ctx, key)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("etcd delete timeout: %v", err)
		}
		return fmt.Errorf("etcd delete error: %v", err)
	}
	return nil
}

// Close ..
func (c *EtcdClient) Close() error {
	return c.client.Close()
}
