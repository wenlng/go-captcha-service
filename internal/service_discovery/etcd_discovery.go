package service_discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdDiscovery struct {
	client *clientv3.Client
	ttl    int64 // seconds
}

func NewEtcdDiscovery(addrs string, ttl int64) (*EtcdDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(addrs, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = client.Status(ctx, addrs)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("etcd health check failed: %v", err)
	}
	return &EtcdDiscovery{client: client, ttl: ttl}, nil
}

func (d *EtcdDiscovery) Register(ctx context.Context, serviceName, instanceID, host string, httpPort, grpcPort int) error {
	lease, err := d.client.Grant(ctx, d.ttl)
	if err != nil {
		return fmt.Errorf("failed to create lease: %v", err)
	}

	instance := Instance{
		InstanceID: instanceID,
		Host:       host,
		HTTPPort:   httpPort,
		GRPCPort:   grpcPort,
	}
	data, err := json.Marshal(instance)
	if err != nil {
		return fmt.Errorf("failed to marshal instance: %v", err)
	}

	key := fmt.Sprintf("/services/%s/%s", serviceName, instanceID)
	_, err = d.client.Put(ctx, key, string(data), clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("failed to register instance: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, err := d.client.KeepAliveOnce(context.Background(), lease.ID)
				if err != nil {
					return
				}
				time.Sleep(time.Duration(d.ttl/2) * time.Second)
			}
		}
	}()
	return nil
}

func (d *EtcdDiscovery) Deregister(ctx context.Context, instanceID string) error {
	key := fmt.Sprintf("/services/%s/%s", "go-captcha-service", instanceID)
	_, err := d.client.Delete(ctx, key)
	return err
}

func (d *EtcdDiscovery) Discover(ctx context.Context, serviceName string) ([]Instance, error) {
	resp, err := d.client.Get(ctx, fmt.Sprintf("/services/%s/", serviceName), clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to discover instances: %v", err)
	}
	var instances []Instance
	for _, kv := range resp.Kvs {
		var instance Instance
		if err := json.Unmarshal(kv.Value, &instance); err != nil {
			continue
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (d *EtcdDiscovery) Close() error {
	return d.client.Close()
}
