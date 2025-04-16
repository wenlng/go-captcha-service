package service_discovery

import (
	"context"
)

// ServiceDiscovery defines the interface for service discovery
type ServiceDiscovery interface {
	Register(ctx context.Context, serviceName, instanceID, host string, httpPort, grpcPort int) error
	Deregister(ctx context.Context, instanceID string) error
	Discover(ctx context.Context, serviceName string) ([]Instance, error)
	Close() error
}

// Instance represents a service instance
type Instance struct {
	InstanceID string
	Host       string
	HTTPPort   int
	GRPCPort   int
	Metadata   map[string]string
}
