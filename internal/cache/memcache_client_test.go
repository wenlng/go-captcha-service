package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemcacheClient(t *testing.T) {
	client, err := NewMemcacheClient("localhost:11211", "TEST_KEY:", 60*time.Second)
	assert.NoError(t, err)
	defer client.Close()

	t.Run("SetAndGet", func(t *testing.T) {
		err := client.SetCache(context.Background(), "key1", "value1")
		if err != nil {
			t.Skip("Memcached not running")
		}

		value, err := client.GetCache(context.Background(), "key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value)
	})

	t.Run("GetNonExistent", func(t *testing.T) {
		value, err := client.GetCache(context.Background(), "nonexistent")
		assert.NoError(t, err)
		assert.Equal(t, "", value)
	})
}
