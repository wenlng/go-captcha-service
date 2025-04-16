package load_balancer

import (
	"fmt"

	"github.com/wenlng/go-captcha-service/internal/service_discovery"
)

// RoundRobin implements round-robin load balancing
type RoundRobin struct {
	index int
}

// NewRoundRobin .
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

// Select selects an instance using round-robin
func (lb *RoundRobin) Select(instances []service_discovery.Instance) (service_discovery.Instance, error) {
	if len(instances) == 0 {
		return service_discovery.Instance{}, fmt.Errorf("no instances available")
	}
	lb.index = (lb.index + 1) % len(instances)
	return instances[lb.index], nil
}
