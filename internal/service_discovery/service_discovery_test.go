package service_discovery

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEtcdDiscovery(t *testing.T) {
	discovery, err := NewEtcdDiscovery("localhost:2379", 10)
	assert.NoError(t, err)

	err = discovery.Register(context.Background(), "go-captcha-service", "id1", "127.0.0.1", 8080, 50051)
	assert.NoError(t, err)

	instances, err := discovery.Discover(context.Background(), "go-captcha-service")
	assert.NoError(t, err)
	assert.Empty(t, instances)

	err = discovery.Deregister(context.Background(), "id1")
	assert.NoError(t, err)

	err = discovery.Close()
	assert.NoError(t, err)
}
