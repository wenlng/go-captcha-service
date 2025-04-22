<div align="center">
<h1 style="margin: 0; padding: 0">GoCaptcha Service</h1>
<br/>
<a href="https://goreportcard.com/report/github.com/wenlng/go-captcha-service"><img src="https://goreportcard.com/badge/github.com/wenlng/go-captcha-service"/></a>
<a href="https://godoc.org/github.com/wenlng/go-captcha-service"><img src="https://godoc.org/github.com/wenlng/go-captcha-service?status.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service/releases"><img src="https://img.shields.io/github/v/release/wenlng/go-captcha-service.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service/blob/LICENSE"><img src="https://img.shields.io/badge/License-Apache2.0-green.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service"><img src="https://img.shields.io/github/stars/wenlng/go-captcha-service.svg"/></a>
<a href="https://github.com/wenlng/go-captcha-service"><img src="https://img.shields.io/github/last-commit/wenlng/go-captcha-service.svg"/></a>
</div>

<br/>

`GoCaptcha Service` is a high-performance behavioral CAPTCHA service developed in Go, based on the **[go-captcha](https://github.com/wenlng/go-captcha)** core library. It supports multiple CAPTCHA modes including click, slide, drag, and rotate. The service provides HTTP and gRPC interfaces, integrates with various service discovery mechanisms (Etcd, Nacos, Zookeeper, Consul), distributed caching (Memory, Redis, Etcd, Memcache), and dynamic configuration. It supports both standalone and distributed deployments, aiming to provide a secure and flexible CAPTCHA solution for web applications.

<br/>

> English | [中文](README_zh.md)
<p> ⭐️ If this project is helpful, please give it a star!</p>

<div align="center">
<img src="https://github.com/wenlng/git-assets/blob/master/go-captcha/go-captcha-v2.jpg?raw=true" alt="Poster">
</div>

<br/>
<hr/>
<br/>

## Related Projects

| Name                                                                         | Description                                                                                         |
|----------------------------------------------------------------------------|--------------------------------------------------------------------------------------------|
| [go-captcha](https://github.com/wenlng/go-captcha)                         | Golang CAPTCHA core library                                                                          |
| [document](http://gocaptcha.wencodes.com)                                  | GoCaptcha documentation                                                                             |
| [online demo](http://gocaptcha.wencodes.com/demo/)                         | GoCaptcha online demo                                                                               |
| [go-captcha-service](https://github.com/wenlng/go-captcha-service)         | GoCaptcha service providing HTTP/gRPC interfaces, <br/>supporting standalone and distributed modes (service discovery, load balancing, dynamic configuration), <br/>deployable via binary or Docker images |
| [go-captcha-service-sdk](https://github.com/wenlng/go-captcha-service-sdk) | GoCaptcha service SDK toolkit, including HTTP/gRPC request interfaces, <br/>supporting static mode, service discovery, and load balancing                               |
| [go-captcha-jslib](https://github.com/wenlng/go-captcha-jslib)             | JavaScript CAPTCHA library                                                                          |
| [go-captcha-vue](https://github.com/wenlng/go-captcha-vue)                 | Vue CAPTCHA                                                                                         |
| [go-captcha-react](https://github.com/wenlng/go-captcha-react)             | React CAPTCHA                                                                                       |
| [go-captcha-angular](https://github.com/wenlng/go-captcha-angular)         | Angular CAPTCHA                                                                                     |
| [go-captcha-svelte](https://github.com/wenlng/go-captcha-svelte)           | Svelte CAPTCHA                                                                                      |
| [go-captcha-solid](https://github.com/wenlng/go-captcha-solid)             | Solid CAPTCHA                                                                                       |
| [go-captcha-uni](https://github.com/wenlng/go-captcha-uni)                 | UniApp CAPTCHA, compatible with APP, mini-programs, and quick apps                                      |
| ...                                                                        |                                                                                            |

<br/>
<br/>

## Features

- **Multiple CAPTCHA Modes**: Supports text/image click, slide, drag, and rotate CAPTCHAs.
- **Dual Protocol Support**: Provides RESTful HTTP and gRPC interfaces.
- **Service Discovery**: Integrates with Etcd, Nacos, Zookeeper, and Consul for distributed service registration and discovery.
- **Distributed Caching**: Supports Memory, Redis, Etcd, and Memcache for optimized CAPTCHA data storage.
- **Dynamic Configuration**: Enables real-time configuration updates via Etcd, Nacos, Zookeeper, or Consul.
- **Highly Configurable**: Supports customization of text, fonts, image resources, CAPTCHA dimensions, and generation rules.
- **High Performance**: Built on Go’s concurrency model, suitable for high-traffic scenarios, with distributed architecture ensuring high availability, performance, and responsiveness.
- **Cross-Platform**: Supports deployment via binary, command line, PM2, Docker, and Docker Compose.

<br/>
<br/>

## Installation and Deployment
`GoCaptcha Service` supports multiple deployment methods, including standalone (binary, command line, PM2, Docker) and distributed deployments (with service discovery, distributed caching, and optional dynamic configuration).

### Prerequisites
- Optional: Docker (for containerized deployment)
- Optional: Service discovery/dynamic configuration middleware (Etcd, Nacos, Zookeeper, Consul)
- Optional: Caching services (Redis, Etcd, Memcache)
- Optional: Node.js and PM2 (for PM2 deployment)
- Optional: gRPC client tools (e.g., `grpcurl`)

### Standalone Deployment

#### Binary Deployment

1. Download the latest binary executable for your platform from [Github Releases](https://github.com/wenlng/go-captcha-service/releases).

    ```bash
    ./go-captcha-service-[xxx]
    ```

2. Optional: Configure the application by copying the `config.json` and `gocaptcha.json` files from the repository to the same directory and specifying them at startup.

   ```bash
    ./go-captcha-service-[xxx] -config config.json -gocaptcha-config gocaptcha.json
    ```

3. Access the HTTP interface (e.g., `http://localhost:8080/api/v1/public/get-data?id=click-default-ch`) or gRPC interface (`localhost:50051`).

<br/>
<br/>

#### PM2 Deployment
PM2 is a Node.js process manager that can manage Go services, providing process monitoring and log management.

1. Install Node.js and PM2:

   ```bash
   npm install -g pm2
   ```

2. Create a PM2 configuration file `ecosystem.config.js`:

   ```javascript
   module.exports = {
     apps: [{
       name: 'go-captcha-service',
       script: './go-captcha-service-[xxx]',
       instances: 1,
       autorestart: true,
       watch: false,
       max_memory_restart: '1G',
       env: {
         CONFIG: 'config.json',
         GO_CAPTCHA_CONFIG: 'gocaptcha.json',
         SERVICE_NAME: 'go-captcha-service',
         CACHE_TYPE: 'redis',
         CACHE_ADDRS: 'localhost:6379',
       },
       env_production: {
         CONFIG: 'config.json',
         GO_CAPTCHA_CONFIG: 'gocaptcha.json',
         SERVICE_NAME: 'go-captcha-service',
         CACHE_TYPE: 'redis',
         CACHE_ADDRS: 'localhost:6379',
       }
     }]
   };
   ```

3. Start the service:

   ```bash
   pm2 start ecosystem.config.js
   ```

4. View logs and status:

   ```bash
   pm2 logs go-captcha-service
   pm2 status
   ```

5. Enable auto-start on boot:

   ```bash
   pm2 startup
   pm2 save
   ```

<br/>
<br/>

#### Golang Source Code + Docker Deployment

1. Create a `Dockerfile` for building from source:

   ```dockerfile
    FROM --platform=$BUILDPLATFORM golang:1.23 AS builder
    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    ARG TARGETOS
    ARG TARGETARCH
    RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -v -a -trimpath -o go-captcha-service ./cmd/go-captcha-service
    
    FROM scratch AS binary
    WORKDIR /app
    
    COPY --from=builder /app/go-captcha-service .
    COPY config.json .
    COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
    
    EXPOSE 8080 50051
    
    HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/go-captcha-service", "--health-check"] || exit 1
    
    CMD ["/app/go-captcha-service"]
   ```

2. Build the image:

   ```bash
   docker build -t go-captcha-service:1.0.0 .
   ```

3. Run the container with mounted configuration files:

   ```bash
   docker run -d -p 8080:8080 -p 50051:50051 \
     -v $(pwd)/config.json:/app/config.json \
     -v $(pwd_MAIN)/gocaptcha.json:/app/gocaptcha.json \
     -v $(pwd)/resource/gocaptcha:/app/resource/gocaptcha \
     --name go-captcha-service go-captcha-service:latest
   ```

<br/>
<br/>

#### Official Docker Image

1. Pull the official image:

   ```bash
   docker pull wenlng/go-captcha-service@latest
   ```

2. Run the container:

   ```bash
   docker run -d -p 8080:8080 -p 50051:50051 \
     -v $(pwd)/config.json:/app/config.json \
     -v $(pwd)/gocaptcha.json:/app/gocaptcha.json \
     -v $(pwd)/resource/gocaptcha:/app/resource/gocaptcha \
     --name go-captcha-service wenlng/go-captcha-service:latest
   ```

<br/>
<br/>

### Distributed Deployment

#### Distributed Caching

1. Configure distributed caching (e.g., Redis) in `config.json`:

   ```json
   {
     "cache_type": "redis",
     "cache_ttl": 1800,
     "cache_key_prefix": "GO_CAPTCHA_DATA:",
     "redis_addrs": "localhost:6379"
   }
   ```

2. Start Redis:

   ```bash
   docker run -d -p 6379:6379 --name redis redis:latest
   ```

<br/>
<br/>

#### Dynamic Configuration
Note: When dynamic configuration is enabled, both `config.json` and `gocaptcha.json` are applied simultaneously.

1. Enable dynamic configuration in `config.json` and select middleware (e.g., Etcd):

   ```json
   {
     "enable_dynamic_config": true,
     "dynamic_config": "etcd",
     "dynamic_config_addrs": "localhost:2379"
   }
   ```

2. Start Etcd:

   ```bash
   docker run -d -p 8848:8848 --name etcd bitnami/etcd::latest
   ```

3. Configuration Synchronization and Retrieval
    - At startup, the service decides whether to push or pull configurations based on the `config_version`. If the local version is higher than the remote (e.g., Etcd) version, the local configuration is pushed to override the remote one; otherwise, the remote configuration is pulled to update the local one (non-file-based update).
    - After startup, the dynamic configuration manager continuously monitors remote configuration changes (e.g., in Etcd). When a remote configuration change occurs, it is fetched and compared with the local version, overriding the local configuration if the remote version is higher.

<br/>
<br/>

#### Service Discovery

1. Enable service discovery in `config.json` and select middleware (e.g., Etcd):

   ```json
   {
     "enable_service_discovery": true,
     "service_discovery": "etcd",
     "service_discovery_addrs": "localhost:2379"
   }
   ```

2. Start Etcd:

   ```bash
   docker run -d -p 8848:8848 --name etcd bitnami/etcd::latest
   ```

3. Service Registration and Discovery
    - At startup, the service automatically registers its instance with the service discovery center (e.g., Etcd).
    - After startup, the service monitors changes in service instances. Refer to [go-captcha-service-sdk](https://github.com/wenlng/go-captcha-service-sdk) for load balancing implementations.

<br/>
<br/>

#### Docker Compose Multi-Instance Deployment

Create a `docker-compose.yml` file including multiple service instances, Consul, Redis, ZooKeeper, and Nacos:

```yaml
version: '3'
services:
  captcha-service-1:
    image: wenlng/go-captcha-service:latest
    ports:
      - "8080:8080"
      - "50051:50051"
    volumes:
      - ./config.json:/app/config.json
      - ./gocaptcha.json:/app/gocaptcha.json
      - ./resources/gocaptcha:/app/resources/gocaptcha
    environment:
      - CONFIG=config.json
      - GO_CAPTCHA_CONFIG=gocaptcha.json
      - SERVICE_NAME=go-captcha-service
      - CACHE_TYPE=redis
      - CACHE_ADDRS=localhost:6379
      - ENABLE_DYNAMIC_CONFIG=true
      - DYNAMIC_CONFIG_TYPE=etcd
      - DYNAMIC_CONFIG_ADDRS=localhost:2379
      - ENABLE_SERVICE_DISCOVERY=true
      - SERVICE_DISCOVERY_TYPE=etcd
      - SERVICE_DISCOVERY_ADDRS=localhost:2379 
    depends_on:
      - etcd
      - redis
    restart: unless-stopped

  captcha-service-2:
    image: wenlng/go-captcha-service:latest
    ports:
      - "8081:8080"
      - "50052:50051"
    volumes:
      - ./config.json:/app/config.json
      - ./gocaptcha.json:/app/gocaptcha.json
      - ./resources/gocaptcha:/app/resources/gocaptcha
    environment:
      - CONFIG=config.json
      - GO_CAPTCHA_CONFIG=gocaptcha.json
      - SERVICE_NAME=go-captcha-service
      - CACHE_TYPE=redis
      - CACHE_ADDRS=localhost:6379
      - ENABLE_DYNAMIC_CONFIG=true
      - DYNAMIC_CONFIG_TYPE=etcd
      - DYNAMIC_CONFIG_ADDRS=localhost:2379
      - ENABLE_SERVICE_DISCOVERY=true
      - SERVICE_DISCOVERY_TYPE=etcd
      - SERVICE_DISCOVERY_ADDRS=localhost:2379
    depends_on:
      - etcd
      - redis
    restart: unless-stopped
       
  etcd:
    image: bitnami/etcd:latest
    ports:
      - "2379:2379"
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
    privileged: true
    restart: unless-stopped

  redis:
    image: redis:latest
    ports:
      - "6379:6379"
    restart: unless-stopped
```

Run:

```bash
docker-compose up -d
```

<br/>
<br/>

## Predefined APIs

* Get CAPTCHA
    ```shell
    curl -H "X-API-Key:my-secret-key-123" http://127.0.0.1:8080/api/v1/public/get-data\?id\=click-default-ch
    ```

* Verify CAPTCHA
    ```shell
    curl -X POST -H "X-API-Key:my-secret-key-123" -H "Content-Type:application/json" -d '{"id":"click-default-ch","captchaKey":"xxxx-xxxxx","value": "x1,y1,x2,y2"}' http://127.0.0.1:8181/api/v1/public/check-data
    ```

* Check Verification Status (`data == "ok"` indicates success)
  ```shell
  curl -H "X-API-Key:my-secret-key-123" http://127.0.0.1:8080/api/v1/public/check-status\?captchaKey\=xxxx-xxxx
  ```

* Get Status Info (not exposed to public networks)
  ```shell
  curl -H "X-API-Key:my-secret-key-123" http://127.0.0.1:8080/api/v1/manage/get-status-info\?captchaKey\=xxxx-xxxx
  ```

* Upload Resources (not exposed to public networks)
  ```shell
  curl -X POST -H "X-API-Key:my-secret-key-123" -F "dirname=imagesdir" -F "files=@/path/to/file1.jpg" -F "files=@/path/to/file2.jpg" http://127.0.0.1:8080/api/v1/manage/upload-resource
  ```

* Delete Resources (not exposed to public networks)
  ```shell
  curl -X DELETE -H "X-API-Key:my-secret-key-123" http://127.0.0.1:8080/api/v1/manage/delete-resource?path=xxxxx.jpg
  ```

* Get Resource File List (not exposed to public networks)
  ```shell
  curl -H "X-API-Key:my-secret-key-123" http://127.0.0.1:8080/api/v1/manage/get-resource-list?path=imagesdir
  ```

* Get CAPTCHA Configuration (not exposed to public networks)
  ```shell
  curl -H "X-API-Key:my-secret-key-123" http://127.0.0.1:8080/api/v1/manage/get-config
  ```

* Update CAPTCHA Configuration (non-file update, not exposed to public networks)
  ```shell
  curl -X POST -H "X-API-Key:my-secret-key-123" -H "Content-Type:application/json" -d '{"config_version":3,"resources":{ ... },"builder": { ... }}' http://127.0.0.1:8080/api/v1/manage/update-hot-config
  ```

For more details and gRPC APIs, refer to [go-captcha-service-sdk](https://github.com/wenlng/go-captcha-service-sdk).

<br/>
<br/>

## API Authentication Configuration
If `api-keys` are configured in `config.json`, all HTTP and gRPC APIs require the `X-API-Key` header for authentication.

The `/api/v1/manage` APIs are not allowed to be exposed to public networks due to security concerns. Only the `/api/v1/public` routes should be publicly accessible. These can be proxied through web servers, reverse proxy servers, or gateway software such as Kong, Envoy, Tomcat, or Nginx.

Example Nginx reverse proxy configuration for public routes:

```text
server {
    listen 80;
    server_name example.com;

    # Proxy requests matching /api/v1/public to the backend
    location ^~ /api/v1/public {
        proxy_pass http://localhost:8080; # Assuming the service runs on port 8080
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # Deny requests matching /api/v1/manage
    location ^~ /api/v1/manage {
        deny all; # Deny all requests, return 403
    }
}
```

<br/>
<br/>

## Configuration Details

### Startup Parameters
Note: Startup parameters correspond to fields in `config.json`. It is recommended to use the configuration file.

* `config`: Specifies the configuration file path, default `config.json`.
* `gocaptcha-config`: Specifies the GoCaptcha configuration file path, default `gocaptcha.json`.
* `service-name`: Sets the service name.
* `http-port`: Sets the HTTP server port.
* `grpc-port`: Sets the gRPC server port.
* `redis-addrs`: Sets Redis cluster addresses, comma-separated.
* `etcd-addrs`: Sets Etcd addresses, comma-separated.
* `memcache-addrs`: Sets Memcached addresses, comma-separated.
* `cache-type`: Sets the cache type, supports `redis`, `memory`, `etcd`, `memcache`.
* `cache-ttl`: Sets cache TTL in seconds.
* `cache-key-prefix`: Sets the cache key prefix, default `GO_CAPTCHA_DATA:`.

* `enable-dynamic-config`: Enables dynamic configuration service, default `false`.
* `dynamic-config-type`: Sets the dynamic configuration service type, supports `etcd`, `zookeeper`, `consul`, `nacos`.
* `dynamic-config-addrs`: Sets the dynamic configuration server addresses, comma-separated.
* `dynamic-config-ttl`: Sets the dynamic configuration service registration TTL in seconds, default `10`.
* `dynamic-config-keep-alive`: Sets the dynamic configuration service keep-alive interval in seconds, default `3`.
* `dynamic-config-max-retries`: Sets the maximum retry attempts for dynamic configuration operations, default `3`.
* `dynamic-config-base-retry-delay`: Sets the base retry delay for dynamic configuration in milliseconds, default `3`.
* `dynamic-config-username`: Sets the dynamic configuration service authentication username.
* `dynamic-config-password`: Sets the dynamic configuration service authentication password.
* `dynamic-config-tls-server-name`: Sets the dynamic configuration service TLS server name.
* `dynamic-config-tls-address`: Sets the dynamic configuration service TLS server address.
* `dynamic-config-tls-cert-file`: Sets the dynamic configuration service TLS certificate file path.
* `dynamic-config-tls-key-file`: Sets the dynamic configuration service TLS key file path.
* `dynamic-config-tls-ca-file`: Sets the dynamic configuration service TLS CA file path.

* `enable-service-discovery`: Enables service discovery, default `false`.
* `service-discovery-type`: Sets the service discovery type, supports `etcd`, `zookeeper`, `consul`, `nacos`.
* `service-discovery-addrs`: Sets the service discovery server addresses, comma-separated.
* `service-discovery-ttl`: Sets the service discovery registration TTL in seconds, default `10`.
* `service-discovery-keep-alive`: Sets the service discovery keep-alive interval in seconds, default `3`.
* `service-discovery-max-retries`: Sets the maximum retry attempts for service discovery operations, default `3`.
* `service-discovery-base-retry-delay`: Sets the base retry delay for service discovery in milliseconds, default `3`.
* `service-discovery-username`: Sets the service discovery authentication username.
* `service-discovery-password`: Sets the service discovery authentication password.
* `service-discovery-tls-server-name`: Sets the service discovery TLS server name.
* `service-discovery-tls-address`: Sets the service discovery TLS server address.
* `service-discovery-tls-cert-file`: Sets the service discovery TLS certificate file path.
* `service-discovery-tls-key-file`: Sets the service discovery TLS key file path.
* `service-discovery-tls-ca-file`: Sets the service discovery TLS CA file path.

* `rate-limit-qps`: Sets the rate limit QPS.
* `rate-limit-burst`: Sets the rate limit burst capacity.
* `api-keys`: Sets the API keys, comma-separated.
* `log-level`: Sets the log level, supports `error`, `debug`, `warn`, `info`.
* `health-check`: Runs a health check and exits, default `false`.
* `enable-cors`: Enables Cross-Origin Resource Sharing, default `false`.

<br/>

### Environment Variables
Basic Configuration:

* `CONFIG`: Main configuration file path for loading application settings.
* `GO_CAPTCHA_CONFIG`: CAPTCHA service configuration file path.
* `SERVICE_NAME`: Service name to identify the service instance.
* `HTTP_PORT`: HTTP service listening port.
* `GRPC_PORT`: gRPC service listening port.
* `API_KEYS`: API keys for authentication or authorization.

Cache Configuration:
* `CACHE_TYPE`: Cache type (e.g., `redis`, `memcached`, `memory`, `etcd`).
* `CACHE_ADDRS`: Cache service address list.
* `CACHE_USERNAME`: Cache service authentication username.
* `CACHE_PASSWORD`: Cache service authentication password.

Dynamic Configuration Service:
* `ENABLE_DYNAMIC_CONFIG`: Enables dynamic configuration (`true` to enable).
* `DYNAMIC_CONFIG_TYPE`: Dynamic configuration type (e.g., `consul`, `zookeeper`, `nacos`, `etcd`).
* `DYNAMIC_CONFIG_ADDRS`: Dynamic configuration service address list.
* `DYNAMIC_CONFIG_USERNAME`: Dynamic configuration service authentication username.
* `DYNAMIC_CONFIG_PASSWORD`: Dynamic configuration service authentication password.

Service Discovery:
* `ENABLE_SERVICE_DISCOVERY`: Enables service discovery (`true` to enable).
* `SERVICE_DISCOVERY_TYPE`: Service discovery type (e.g., `consul`, `zookeeper`, `nacos`, `etcd`).
* `SERVICE_DISCOVERY_ADDRS`: Service discovery service address list.
* `SERVICE_DISCOVERY_USERNAME`: Service discovery service authentication username.
* `SERVICE_DISCOVERY_PASSWORD`: Service discovery service authentication password.

<br/>

### Configuration Files
The service uses two configuration files: `config.json` for service runtime parameters and `gocaptcha.json` for CAPTCHA generation settings.

### config.json

`config.json` defines the basic service configuration.

```json
{
   "config_version": 1,
   "service_name": "go-captcha-service",
   "http_port": "8080",
   "grpc_port": "50051",
   "redis_addrs": "localhost:6379",
   "etcd_addrs": "localhost:2379",
   "memcache_addrs": "localhost:11211",
   "cache_type": "memory",
   "cache_ttl": 1800,
   "cache_key_prefix": "GO_CAPTCHA_DATA:",
  
   "enable_dynamic_config": false,
   "dynamic_config_type": "etcd",
   "dynamic_config_addrs": "localhost:2379",
   "dynamic_config_username": "",
   "dynamic_config_password": "",
   "dynamic_config_ttl": 10,
   "dynamic_config_keep_alive": 3,
   "dynamic_config_max_retries": 3,
   "dynamic_config_base_retry_delay": 500,
   "dynamic_config_tls_server_name": "",
   "dynamic_config_tls_address": "",
   "dynamic_config_tls_cert_file": "",
   "dynamic_config_tls_key_file": "",
   "dynamic_config_tls_ca_file": "",
  
   "enable_service_discovery": false,
   "service_discovery_type": "etcd",
   "service_discovery_addrs": "localhost:2379",
   "service_discovery_username": "",
   "service_discovery_password": "",
   "service_discovery_ttl": 10,
   "service_discovery_keep_alive": 3,
   "service_discovery_max_retries": 3,
   "service_discovery_base_retry_delay": 500,
   "service_discovery_tls_server_name": "",
   "service_discovery_tls_address": "",
   "service_discovery_tls_cert_file": "",
   "service_discovery_tls_key_file": "",
   "service_discovery_tls_ca_file": "",
  
   "rate_limit_qps": 1000,
   "rate_limit_burst": 1000,
   "enable_cors": true,
   "log_level": "info",
   "api_keys": ["my-secret-key-123", "another-key-456", "another-key-789"]
}
```

#### Parameter Descriptions

- `config_version` (integer): Configuration file version for distributed dynamic configuration, default `1`.
- `service_name` (string): Service name, default `go-captcha-service`.
- `http_port` (string): HTTP port, default `8080`.
- `grpc_port` (string): gRPC port, default `50051`.
- `redis_addrs` (string): Redis address, default `localhost:6379`. Used when `cache_type: redis`.
- `etcd_addrs` (string): Etcd address, default `localhost:2379`. Used when `cache_type: etcd` or `service_discovery: etcd`.
- `memcache_addrs` (string): Memcache address, default `localhost:11211`. Used when `cache_type: memcache`.
- `cache_type` (string): Cache type, default `memory`:
    - `memory`: In-memory cache, suitable for standalone deployment.
    - `redis`: Distributed key-value store, suitable for high-availability scenarios.
    - `etcd`: Distributed key-value store, suitable for sharing with service discovery.
    - `memcache`: High-performance distributed cache, suitable for high concurrency.
- `cache_ttl` (integer): Cache expiration time in seconds, default `1800`.
- `cache_key_prefix` (string): Cache key prefix, default `GO_CAPTCHA_DATA:`.

- `enable_dynamic_config` (boolean): Enables dynamic configuration service, default `false`.
- `dynamic_config_type` (string): Dynamic configuration service type, default `etcd`:
    - `etcd`: Suitable for high-consistency scenarios.
    - `nacos`: Suitable for cloud-native environments.
    - `zookeeper`: Suitable for complex distributed systems.
    - `consul`: Lightweight, supports health checks.
- `dynamic_config_addrs` (string): Dynamic configuration service addresses, e.g., Etcd: `localhost:2379`, Nacos: `localhost:8848`.
- `dynamic_config_username` (string): Username, e.g., Nacos default username is `nacos`, default empty.
- `dynamic_config_password` (string): Password, e.g., Nacos default password is `nacos`, default empty.
- `dynamic_config_ttl` (integer): Service lease time in seconds, default `10`.
- `dynamic_config_keep_alive` (integer): Heartbeat interval in seconds, default `3`.
- `dynamic_config_max_retries` (integer): Retry attempts, default `3`.
- `dynamic_config_base_retry_delay` (integer): Retry delay in milliseconds, default `500`.
- `dynamic_config_tls_server_name` (string): TLS server name, default empty.
- `dynamic_config_tls_address` (string): TLS server's address, default empty.
- `dynamic_config_tls_cert_file` (string): TLS certificate file, default empty.
- `dynamic_config_tls_key_file` (string): TLS key file, default empty.
- `dynamic_config_tls_ca_file` (string): TLS CA certificate file, default empty.

- `enable_service_discovery` (boolean): Enables service discovery, default `false`.
- `service_discovery_type` (string): Service discovery type, default `etcd`:
    - `etcd`: Suitable for high-consistency scenarios.
    - `nacos`: Suitable for cloud-native environments.
    - `zookeeper`: Suitable for complex distributed systems.
    - `consul`: Lightweight, supports health checks.
- `service_discovery_addrs` (string): Service discovery addresses, e.g., Etcd: `localhost:2379`, Nacos: `localhost:8848`.
- `service_discovery_username` (string): Username, e.g., Nacos default username is `nacos`, default empty.
- `service_discovery_password` (string): Password, e.g., Nacos default password is `nacos`, default empty.
- `service_discovery_ttl` (integer): Service registration lease time in seconds, default `10`.
- `service_discovery_keep_alive` (integer): Heartbeat interval in seconds, default `3`.
- `service_discovery_max_retries` (integer): Retry attempts, default `3`.
- `service_discovery_base_retry_delay` (integer): Retry delay in milliseconds, default `500`.
- `service_discovery_tls_server_name` (string): TLS server name, default empty.
- `service_discovery_tls_address` (string): TLS server's address, default empty.
- `service_discovery_tls_cert_file` (string): TLS certificate file, default empty.
- `service_discovery_tls_key_file` (string): TLS key file, default empty.
- `service_discovery_tls_ca_file` (string): TLS CA certificate file, default empty.

- `rate_limit_qps` (integer): API requests per second limit, default `1000`.
- `rate_limit_burst` (integer): API burst capacity limit, default `1000`.
- `enable_cors` (boolean): Enables CORS, default `true`.
- `log_level` (string): Log level (`debug`, `info`, `warn`, `error`), default `info`.
- `api_keys` (string array): API authentication keys.

### gocaptcha.json

`gocaptcha.json` defines resources and generation settings for CAPTCHAs.

```json
{
  "config_version": 1,
  "resources": {
    "version": "0.0.1",
    "char": {
      "languages": {
        "chinese": [],
        "english": []
      }
    },
    "font": {
      "type": "load",
      "file_dir": "./gocaptcha/fonts/",
      "file_maps": {
        "yrdzst_bold": "yrdzst-bold.ttf"
      }
    },
    "shape_image": {
      "type": "load",
      "file_dir": "./gocaptcha/shape_images/",
      "file_maps": {
        "shape_01": "shape_01.png",
        "shape_01.png":"c.png"
      }
    },
    "master_image": {
      "type": "load",
      "file_dir": "./gocaptcha/master_images/",
      "file_maps": {
        "image_01": "image_01.jpg",
        "image_02":"image_02.jpg"
      }
    },
    "thumb_image": {
      "type": "load",
      "file_dir": "./gocaptcha/thumb_images/",
      "file_maps": {

      }
    },
    "tile_image": {
      "type": "load",
      "file_dir": "./gocaptcha/tile_images/",
      "file_maps": {
        "tile_01": "tile_01.png",
        "tile_02": "tile_02.png"
      },
      "file_maps_02": {
        "tile_mask_01": "tile_mask_01.png",
        "tile_mask_02": "tile_mask_02.png"
      },
      "file_maps_03": {
        "tile_shadow_01": "tile_shadow_01.png",
        "tile_shadow_02": "tile_shadow_02.png"
      }
    }
  },
  "builder": {
    "click_config_maps": {
      "click-default-ch": {
        "version": "0.0.1",
        "language": "chinese",
        "master": {
          "image_size": { "width": 300, "height": 200 },
          "range_length": { "min": 6, "max": 7 },
          "range_angles": [
            { "min": 20, "max": 35 },
            { "min": 35, "max": 45 },
            { "min": 290, "max": 305 },
            { "min": 305, "max": 325 },
            { "min": 325, "max": 330 }
          ],
          "range_size": { "min": 26, "max": 32 },
          "range_colors": [ "#fde98e", "#60c1ff", "#fcb08e", "#fb88ff", "#b4fed4", "#cbfaa9", "#78d6f8"],
          "display_shadow": true,
          "shadow_color": "#101010",
          "shadow_point": { "x": -1, "y": -1 },
          "image_alpha": 1,
          "use_shape_original_color": true
        },
        "thumb": {
          "image_size": { "width": 150, "height": 40 },
          "range_verify_length": { "min": 2, "max": 4 },
          "disabled_range_verify_length": false,
          "range_text_size": { "min": 22, "max": 28 },
          "range_text_colors": [ "#1f55c4", "#780592", "#2f6b00", "#910000", "#864401", "#675901", "#016e5c"],
          "range_background_colors": ["#1f55c4", "#780592", "#2f6b00", "#910000", "#864401", "#675901", "#016e5c"],
          "is_non_deform_ability": false,
          "background_distort": 4,
          "background_distort_alpha": 1,
          "background_circles_num": 24,
          "background_slim_line_num": 2
        }
      },
      "click-dark-ch": {
        "version": "0.0.1",
        "language": "chinese",
        // Same as above...
      },
      "click-default-en": {
        "version": "0.0.1",
        "language": "english",
        // Same as above...
      },
      "click-dark-en": {
        "version": "0.0.1",
        "language": "english",
        // Same as above...
      }
    },
    "click_shape_config_maps": {
      "click-shape-default":  {
        "version": "0.0.1",
        "master": {
          "image_size": { "width": 300, "height": 200 },
          "range_length": { "min": 6, "max": 7 },
          "range_angles": [
            { "min": 20, "max": 35 },
            { "min": 35, "max": 45 },
            { "min": 290, "max": 305 },
            { "min": 305, "max": 325 },
            { "min": 325, "max": 330 }
          ],
          "range_size": { "min": 26, "max": 32 },
          "range_colors": [ "#fde98e", "#60c1ff", "#fcb08e", "#fb88ff", "#b4fed4", "#cbfaa9", "#78d6f8"],
          "display_shadow": true,
          "shadow_color": "#101010",
          "shadow_point": { "x": -1, "y": -1 },
          "image_alpha": 1,
          "use_shape_original_color": true
        },
        "thumb": {
          "image_size": { "width": 150, "height": 40},
          "range_verify_length": { "min": 2, "max": 4 },
          "disabled_range_verify_length": false,
          "range_text_size": { "min": 22, "max": 28},
          "range_text_colors": [ "#1f55c4", "#780592", "#2f6b00", "#910000", "#864401", "#675901", "#016e5c"],
          "range_background_colors": [ "#1f55c4", "#780592", "#2f6b00", "#910000", "#864401", "#675901", "#016e5c" ],
          "is_non_deform_ability": false,
          "background_distort": 4,
          "background_distort_alpha": 1,
          "background_circles_num": 24,
          "background_slim_line_num": 2
        }
      }
    },
    "slide_config_maps": {
      "slide-default": {
        "version": "0.0.1",
        "master": {
          "image_size": { "width": 300, "height": 200 },
          "image_alpha": 1
        },
        "thumb": {
          "range_graph_size":  { "min": 60, "max": 70 },
          "range_graph_angles": [
            { "min": 20, "max": 35 },
          ],
          "generate_graph_number": 1,
          "enable_graph_vertical_random": false,
          "range_dead_zone_directions": ["left", "right"]
        }
      }
    },
    "drag_config_maps": {
      "drag-default": {
        "version": "0.0.1",
        "master": {
          "image_size": { "width": 300, "height": 200 },
          "image_alpha": 1
        },
        "thumb": {
          "range_graph_size":  { "min": 60, "max": 70 },
          "range_graph_angles": [
            { "min": 0, "max": 0 },
          ],
          "generate_graph_number": 2,
          "enable_graph_vertical_random": true,
          "range_dead_zone_directions": ["left", "right", "top", "bottom"]
        }
      }
    },
    "rotate_config_maps": {
      "rotate-default": {
        "version": "0.0.1",
        "master": {
          "image_square_size": 220,
        },
        "thumb": {
          "range_angles":  [{ "min": 30, "max": 330 }],
          "range_image_square_sizes":  [140, 150, 160, 170],
          "image_alpha":  1
        }
      }
    }
  }
}
```

<br/>

##### Top-Level Fields

- `config_version` (integer): Configuration file version for distributed dynamic configuration management, default `1`.

##### resources Field

- `version` (string): Resource configuration version to control CAPTCHA instance recreation, default `0.0.1`.
- `char.languages.chinese` (string array): Chinese character set for click CAPTCHA text, default empty (uses built-in resources).
- `char.languages.english` (string array): English character set, default empty (uses built-in resources).
- `font.type` (string): Font loading method, fixed as `load` (load from file).
- `font.file_dir` (string): Font file directory, default `./gocaptcha/fonts/`.
- `font.file_maps` (object): Font file mappings, key is the font name, value is the file name.
    - Example: `"yrdzst_bold": "yrdzst-bold.ttf"` uses `yrdzst-bold.ttf` font.
- `shape_image.type` (string): Shape image loading method, fixed as `load`.
- `shape_image.file_dir` (string): Shape image directory, default `./gocaptcha/shape_images/`.
- `shape_image.file_maps` (object): Shape image mappings.
    - Example: `"shape_01": "shape_01.png"` uses `shape_01.png` as a shape.
- `master_image.type` (string): Main image loading method, fixed as `load`.
- `master_image.file_dir` (string): Main image directory, default `./gocaptcha/master_images/`.
- `master_image.file_maps` (object): Main image mappings.
    - Example: `"image_01": "image_01.jpg"` uses `image_01.jpg` as the background.
- `thumb_image.type` (string): Thumbnail image loading method, fixed as `load`.
- `thumb_image.file_dir` (string): Thumbnail image directory, default `./gocaptcha/thumb_images/`.
- `thumb_image.file_maps` (object): Thumbnail image mappings, default empty.
- `tile_image.type` (string): Tile image loading method, fixed as `load`.
- `tile_image.file_dir` (string): Tile image directory, default `./gocaptcha/tile_images/`.
- `tile_image.file_maps` (object): Tile image mappings.
    - Example: `"tile_01": "tile_01.png"`.
- `tile_image.file_maps_02` (object): Tile mask mappings.
    - Example: `"tile_mask_01": "tile_mask_01.png"`.
- `tile_image.file_maps_03` (object): Tile shadow mappings.
    - Example: `"tile_shadow_01": "tile_shadow_01.png"`.

<br/>

##### builder Field

Defines CAPTCHA generation styles, including configurations for click, shape click, slide, drag, and rotate CAPTCHAs.

###### click_config_maps

Defines text click CAPTCHA configurations, supporting Chinese and English with light and dark themes. The key is the ID passed in the CAPTCHA API request, e.g., `api/v1/public/get-data?id=click-default-ch`.

- `click-default-ch` (object): Default Chinese theme configuration.
    - `version` (string): Configuration version to control CAPTCHA instance recreation, default `0.0.1`.
      $
    - `language` (string): Language, matches defined `char.languages`, e.g., `chinese` for Chinese.
    - `master` (object): Main CAPTCHA image configuration.
        - `image_size.width` (integer): Main image width, default `300`.
        - `image_size.height` (integer): Main image height, default `200`.
        - `range_length.min` (integer): Minimum number of CAPTCHA points, default `6`.
        - `range_length.max` (integer): Maximum number of CAPTCHA points, default `7`.
        - `range_angles` (object array): Text rotation angle ranges (degrees).
            - Example: `{"min": 20, "max": 35}` for 20°-35°.
        - `range_size.min` (integer): Minimum text size (pixels), default `26`.
        - `range_size.max` (integer): Maximum text size, default `32`.
        - `range_colors` (string array): Text color list (hexadecimal).
            - Example: `"#fde98e"`.
        - `display_shadow` (boolean): Display text shadow, default `true`.
        - `shadow_color` (string): Shadow color, default `#101010`.
        - `shadow_point.x` (integer): Shadow offset X coordinate, default `-1` (auto-calculated).
        - `shadow_point.y` (integer): Shadow offset Y coordinate, default `-1`.
        - `image_alpha` (float): Image opacity (0-1), default `1`.
        - `use_shape_original_color` (boolean): Use original shape color, default `true`.
    - `thumb` (object): Thumbnail (prompt text) configuration.
        - `image_size.width` (integer): Thumbnail width, default `150`.
        - `image_size.height` (integer): Thumbnail height, default `40`.
        - `range_verify_length.min` (integer): Minimum verification points, default `2`.
        - `range_verify_length.max` (integer): Maximum verification points, default `4`.
        - `disabled_range_verify_length` (boolean): Disable verification point limit, default `false`.
        - `range_text_size.min` (integer): Minimum text size, default `22`.
        - `range_text_size.max` (integer): Maximum text size, default `28`.
        - `range_text_colors` (string array): Text color list.
        - `range_background_colors` (string array): Background color list.
        - `is_non_deform_ability` (boolean): Disable deformation effect, default `false`.
        - `background_distort` (integer): Background distortion level, default `4`.
        - `background_distort_alpha` (float): Background distortion opacity, default `1`.
        - `background_circles_num` (integer): Number of background circle interference points, default `24`.
        - `background_slim_line_num` (integer): Number of background slim line interferences, default `2`.

- `click-dark-ch` (object): Chinese dark theme configuration, similar to `click-default-ch`, but `thumb.range_text_colors` uses brighter colors for dark backgrounds.

- `click-default-en` (object): Default English theme configuration, with `language: english`, larger `master.range_size` and `thumb.range_text_size` (`34-48`) for English characters.

- `click-dark-en` (object): English dark theme configuration, similar to `click-dark-ch`, with `language: english`.

<br/>

###### click_shape_config_maps

Defines shape click CAPTCHA configurations.

- `click-shape-default` (object): Default shape click configuration, similar to `click_config_maps` `master` and `thumb` but for shape images instead of text.

<br/>

###### slide_config_maps

Defines slide CAPTCHA configurations.

- `slide-default` (object):
    - `version` (string): Configuration version to control CAPTCHA instance recreation, default `0.0.1`.
    - `master` (object): Main CAPTCHA image configuration.
        - `image_size.width` (integer): Main image width, default `300`.
        - `image_size.height` (integer): Main image height, default `200`.
        - `image_alpha` (float): Image opacity (0-1), default `1`.
    - `thumb` (object): Slider configuration.
        - `range_graph_size.min` (integer): Minimum slider graphic size (pixels), default `60`.
        - `range_graph_size.max` (integer): Maximum slider graphic size, default `70`.
        - `range_graph_angles` (object array): Slider graphic rotation angle ranges (degrees).
            - Example: `{"min": 20, "max": 35}`.
        - `generate_graph_number` (integer): Number of slider graphics to generate, default `1`.
        - `enable_graph_vertical_random` (boolean): Enable vertical random offset, default `false`.
        - `range_dead_zone_directions` (string array): Slider dead zone directions, default `["left", "right"]`.

<br/>

###### drag_config_maps

Defines drag CAPTCHA configurations.

- `drag-default` (object):
    - `version` (string): Configuration version to control CAPTCHA instance recreation, default `0.0.1`.
    - `master` (object): Main CAPTCHA image configuration.
        - `image_size.width` (integer): Main image width, default `300`.
        - `image_size.height` (integer): Main image height, default `200`.
        - `image_alpha` (float): Image opacity (0-1), default `1`.
    - `thumb` (object): Drag graphic configuration.
        - `range_graph_size.min` (integer): Minimum drag graphic size (pixels), default `60`.
        - `range_graph_size.max` (integer): Maximum drag graphic size, default `70`.
        - `range_graph_angles` (object array): Drag graphic rotation angle ranges (degrees).
            - Example: `{"min": 0, "max": 0}` for no rotation.
        - `generate_graph_number` (integer): Number of drag graphics to generate, default `2`.
        - `enable_graph_vertical_random` (boolean): Enable vertical random offset, default `true`.
        - `range_dead_zone_directions` (string array): Drag dead zone directions, default `["left", "right", "top", "bottom"]`.

<br/>

###### rotate_config_maps

Defines rotate CAPTCHA configurations.

- `rotate-default` (object):
    - `version` (string): Configuration version to control CAPTCHA instance recreation, default `0.0.1`.
    - `master` (object): Main CAPTCHA image configuration.
        - `image_square_size` (integer): Main image square side length (pixels), default `220`.
    - `thumb` (object): Rotate graphic configuration.
        - `range_angles` (object array): Rotation angle ranges (degrees).
            - Example: `{"min": 30, "max": 330}` for 30°-330°.
        - `range_image_square_sizes` (integer array): Rotate image square side length list, default `[140, 150, 160, 170]`.
        - `image_alpha` (float): Image opacity (0-1), default `1`.

<br/>
<br/>

### Configuration Hot Reloading
Hot reloading for `gocaptcha.json` is determined by the `version` field of each configuration item.

Hot-reloadable fields in `config.json` include:
* `cache_type`
* `cache_addrs`
* `cache_ttl`
* `cache_key_prefix`
* `api_keys`
* `log_level`
* `rate_limit_qps`
* `rate_limit_burst`

### Testing

- CAPTCHA Generation: Verify image, shape, and key validity.
- Verification Logic: Test handling of different inputs.
- Service Discovery: Simulate Etcd/Nacos/Zookeeper/Consul.
- Caching: Test Memory/Redis/Etcd/Memcache.
- Dynamic Configuration: Test Nacos configuration updates.

<br/>
<br/>


## LICENSE
Go Captcha Service source code is licensed under the Apache License, Version 2.0 [http://www.apache.org/licenses/LICENSE-2.0.html](http://www.apache.org/licenses/LICENSE-2.0.html)