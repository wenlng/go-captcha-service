package load_balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wenlng/go-captcha-service/internal/service_discovery"
)

func TestRoundRobin(t *testing.T) {
	lb := NewRoundRobin()
	instances := []service_discovery.Instance{
		{InstanceID: "1", Host: "127.0.0.1", HTTPPort: 8080, GRPCPort: 50051},
		{InstanceID: "2", Host: "127.0.0.2", HTTPPort: 8081, GRPCPort: 50052},
	}

	instance, err := lb.Select(instances)
	assert.NoError(t, err)
	assert.Equal(t, "1", instance.InstanceID)

	instance, err = lb.Select(instances)
	assert.NoError(t, err)
	assert.Equal(t, "2", instance.InstanceID)

	instance, err = lb.Select(instances)
	assert.NoError(t, err)
	assert.Equal(t, "1", instance.InstanceID)
}

func TestConsistentHash(t *testing.T) {
	lb := NewConsistentHash()
	instances := []service_discovery.Instance{
		{InstanceID: "1", Host: "127.0.0.1", HTTPPort: 8080, GRPCPort: 50051},
	}

	instance, err := lb.Select(instances)
	assert.NoError(t, err)
	assert.Equal(t, "1", instance.InstanceID)

	_, err = lb.Select([]service_discovery.Instance{})
	assert.Error(t, err)
}
