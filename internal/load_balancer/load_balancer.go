package load_balancer

import (
	"github.com/wenlng/go-captcha-service/internal/service_discovery"
)

// LoadBalancer .
type LoadBalancer interface {
	Select(instances []service_discovery.Instance) (service_discovery.Instance, error)
}
