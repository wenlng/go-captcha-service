package service_discovery

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosDiscovery struct {
	client naming_client.INamingClient
}

func NewNacosDiscovery(addrs string, ttl int64) (*NacosDiscovery, error) {
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
	)
	serverConfigs := []constant.ServerConfig{}
	for _, addr := range strings.Split(addrs, ",") {
		hostPort := strings.Split(addr, ":")
		host := hostPort[0]
		port, _ := strconv.Atoi(hostPort[1])
		serverConfigs = append(serverConfigs, *constant.NewServerConfig(host, uint64(port)))
	}
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Nacos: %v", err)
	}
	return &NacosDiscovery{client: namingClient}, nil
}

func (d *NacosDiscovery) Register(ctx context.Context, serviceName, instanceID, host string, httpPort, grpcPort int) error {
	_, err := d.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          host,
		Port:        uint64(httpPort),
		ServiceName: serviceName,
		GroupName:   instanceID,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"grpc_port": fmt.Sprintf("%d", grpcPort),
		},
	})

	return err
}

func (d *NacosDiscovery) Deregister(ctx context.Context, instanceID string) error {
	_, err := d.client.DeregisterInstance(vo.DeregisterInstanceParam{
		GroupName:   instanceID,
		ServiceName: "go-captcha-service",
	})
	return err
}

func (d *NacosDiscovery) Discover(ctx context.Context, serviceName string) ([]Instance, error) {
	instances, err := d.client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to discover instances: %v", err)
	}
	var result []Instance
	for _, inst := range instances.Hosts {
		grpcPort, _ := strconv.Atoi(inst.Metadata["grpc_port"])
		result = append(result, Instance{
			InstanceID: inst.InstanceId,
			Host:       inst.Ip,
			HTTPPort:   int(inst.Port),
			GRPCPort:   grpcPort,
			Metadata:   inst.Metadata,
		})
	}
	return result, nil
}

func (d *NacosDiscovery) Close() error {
	return nil
}
