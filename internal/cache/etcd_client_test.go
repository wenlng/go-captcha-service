package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/server/v3/embed"
)

func TestEtcdClient(t *testing.T) {
	cfg := embed.NewConfig()
	cfg.Dir = t.TempDir()
	etcd, err := embed.StartEtcd(cfg)
	assert.NoError(t, err)
	defer etcd.Close()

	client, err := NewEtcdClient("localhost:2379", "TEST_KEY:", 60*time.Second)
	assert.NoError(t, err)
	defer client.Close()

	t.Run("SetAndGet", func(t *testing.T) {
		err := client.SetCache(context.Background(), "key1", "value1")
		assert.NoError(t, err)

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
