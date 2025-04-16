package cache

import (
	"context"
)

// Cache defines the interface for cache operations
type Cache interface {
	GetCache(ctx context.Context, key string) (string, error)
	SetCache(ctx context.Context, key, value string) error
	DeleteCache(ctx context.Context, key string) error
	Close() error
}

type CaptCacheData struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}
