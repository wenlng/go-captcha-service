package service_discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

type ZookeeperDiscovery struct {
	conn *zk.Conn
	ttl  int64 // seconds
}

func NewZookeeperDiscovery(addrs string, ttl int64) (*ZookeeperDiscovery, error) {
	conn, _, err := zk.Connect(strings.Split(addrs, ","), 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ZooKeeper: %v", err)
	}
	_, _, err = conn.Children("/")
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("ZooKeeper health check failed: %v", err)
	}
	return &ZookeeperDiscovery{conn: conn, ttl: ttl}, nil
}

func (d *ZookeeperDiscovery) Register(ctx context.Context, serviceName, instanceID, host string, httpPort, grpcPort int) error {
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

	path := fmt.Sprintf("/services/%s/%s", serviceName, instanceID)
	_, err = d.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		return fmt.Errorf("failed to register instance: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				exists, _, err := d.conn.Exists(path)
				if err != nil || !exists {
					d.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
				}
				time.Sleep(time.Duration(d.ttl/2) * time.Second)
			}
		}
	}()
	return nil
}

func (d *ZookeeperDiscovery) Deregister(ctx context.Context, instanceID string) error {
	path := fmt.Sprintf("/services/%s/%s", "go-captcha-service", instanceID)
	return d.conn.Delete(path, -1)
}

func (d *ZookeeperDiscovery) Discover(ctx context.Context, serviceName string) ([]Instance, error) {
	path := fmt.Sprintf("/services/%s", serviceName)
	children, _, err := d.conn.Children(path)
	if err != nil {
		return nil, fmt.Errorf("failed to discover instances: %v", err)
	}
	var instances []Instance
	for _, child := range children {
		data, _, err := d.conn.Get(path + "/" + child)
		if err != nil {
			continue
		}
		var instance Instance
		if err := json.Unmarshal(data, &instance); err != nil {
			continue
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (d *ZookeeperDiscovery) Close() error {
	d.conn.Close()
	return nil
}
