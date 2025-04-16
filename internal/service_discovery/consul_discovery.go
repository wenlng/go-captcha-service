package service_discovery

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
)

type ConsulDiscovery struct {
	client *api.Client
	ttl    int
}

func NewConsulDiscovery(addrs string, ttl int) (*ConsulDiscovery, error) {
	cfg := api.DefaultConfig()
	cfg.Address = strings.Split(addrs, ",")[0]
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Consul: %v", err)
	}
	_, err = client.Status().Leader()
	if err != nil {
		return nil, fmt.Errorf("Consul health check failed: %v", err)
	}
	return &ConsulDiscovery{client: client, ttl: ttl}, nil
}

func (d *ConsulDiscovery) Register(ctx context.Context, serviceName, instanceID, host string, httpPort, grpcPort int) error {
	reg := &api.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: host,
		Port:    httpPort,
		Tags:    []string{"http", "grpc"},
		Meta: map[string]string{
			"grpc_port": fmt.Sprintf("%d", grpcPort),
		},
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/hello", host, httpPort),
			Interval: fmt.Sprintf("%ds", d.ttl),
			Timeout:  "5s",
			TTL:      fmt.Sprintf("%ds", d.ttl),
		},
	}
	return d.client.Agent().ServiceRegister(reg)
}

func (d *ConsulDiscovery) Deregister(ctx context.Context, instanceID string) error {
	return d.client.Agent().ServiceDeregister(instanceID)
}

func (d *ConsulDiscovery) Discover(ctx context.Context, serviceName string) ([]Instance, error) {
	services, _, err := d.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover instances: %v", err)
	}
	var instances []Instance
	for _, entry := range services {
		grpcPort, _ := strconv.Atoi(entry.Service.Meta["grpc_port"])
		instances = append(instances, Instance{
			InstanceID: entry.Service.ID,
			Host:       entry.Service.Address,
			HTTPPort:   entry.Service.Port,
			GRPCPort:   grpcPort,
			Metadata:   entry.Service.Meta,
		})
	}
	return instances, nil
}

func (d *ConsulDiscovery) Close() error {
	return nil
}
