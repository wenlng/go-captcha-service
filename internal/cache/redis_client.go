/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient implements the Cache interface for Redis
type RedisClient struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

// NewRedisClient ..
func NewRedisClient(addrs, prefix string, ttl time.Duration, username, password, db string) (*RedisClient, error) {
	var dbNum int
	if db != "" {
		dbNum, _ = strconv.Atoi(db)
	}

	client := redis.NewClient(&redis.Options{
		Addr:            addrs,
		Username:        username,
		Password:        password,
		DB:              dbNum,
		PoolSize:        100,
		MinIdleConns:    10,
		MaxRetries:      3,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     4 * time.Second,
		ConnMaxIdleTime: 30 * time.Minute,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("redis ping failed: %v", err)
	}
	return &RedisClient{client: client, prefix: prefix, ttl: ttl}, nil
}

// GetCache retrieves a value from Redis
func (c *RedisClient) GetCache(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	key = c.prefix + key
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("redis get timeout: %v", err)
		}
		return "", fmt.Errorf("redis get error: %v", err)
	}
	return val, nil
}

// SetCache stores a value in Redis
func (c *RedisClient) SetCache(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	key = c.prefix + key
	err := c.client.Set(ctx, key, value, c.ttl).Err()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("redis set timeout: %v", err)
		}
		return fmt.Errorf("redis set error: %v", err)
	}
	return nil
}

// DeleteCache delete a value in Redis
func (c *RedisClient) DeleteCache(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	key = c.prefix + key
	err := c.client.Del(ctx, key).Err()
	if err != nil && err != redis.Nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("redis delete timeout: %v", err)
		}
		return fmt.Errorf("redis delete error: %v", err)
	}
	return nil
}

// Close ..
func (c *RedisClient) Close() error {
	return c.client.Close()
}
