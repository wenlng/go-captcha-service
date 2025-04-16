package load_balancer

import (
	"fmt"

	"github.com/wenlng/go-captcha-service/internal/service_discovery"
)

// ConsistentHash .
type ConsistentHash struct{}

// NewConsistentHash .
func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{}
}

// Select selects an instance using consistent hashing
func (lb *ConsistentHash) Select(instances []service_discovery.Instance) (service_discovery.Instance, error) {
	if len(instances) == 0 {
		return service_discovery.Instance{}, fmt.Errorf("no instances available")
	}
	// select first instance
	return instances[0], nil
}
