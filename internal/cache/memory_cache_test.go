package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache("TEST_KEY:", 1*time.Second, 500*time.Millisecond)
	defer cache.Stop()

	t.Run("SetAndGet", func(t *testing.T) {
		err := cache.SetCache(context.Background(), "key1", "value1")
		assert.NoError(t, err)

		value, err := cache.GetCache(context.Background(), "key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value)
	})

	t.Run("Expiration", func(t *testing.T) {
		err := cache.SetCache(context.Background(), "key2", "value2")
		assert.NoError(t, err)

		time.Sleep(1100 * time.Millisecond)

		value, err := cache.GetCache(context.Background(), "key2")
		assert.NoError(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("Cleanup", func(t *testing.T) {
		err := cache.SetCache(context.Background(), "key3", "value3")
		assert.NoError(t, err)

		time.Sleep(600 * time.Millisecond)

		cache.mu.RLock()
		_, exists := cache.items["key3"]
		cache.mu.RUnlock()
		assert.False(t, exists)
	})
}
