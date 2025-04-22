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

`GoCaptcha Service` 是一个基于 Go 语言开发的高性能行为验证码服务，基于 **[go-captcha](https://github.com/wenlng/go-captcha)** 行为验证码基本库，支持点击、滑动、拖拽和旋转等多种验证码模式。它提供 HTTP 和 gRPC 接口，集成多种服务发现（Etcd、Nacos、Zookeeper、Consul）、分布式缓存（Memory、Redis、Etcd、Memcache）和分布式动态配置，支持单机和分布式部署，旨在为 Web 应用提供安全、灵活的验证码解决方案。

<br/>

> English | [中文](README_zh.md)
<p> ⭐️ 如果能帮助到你，请随手给点一个star</p>
<p>QQ交流群：178498936</p>

<div align="center">
<img src="https://github.com/wenlng/git-assets/blob/master/go-captcha/go-captcha-v2.jpg?raw=true" alt="Poster">
</div>

<br/>
<hr/>
<br/>

## 项目生态

| 名称                                                                         | 描述                                                                                          |
|----------------------------------------------------------------------------|---------------------------------------------------------------------------------------------|
| [document](http://gocaptcha.wencodes.com)                                  | GoCaptcha 文档                                                                                |
| [online demo](http://gocaptcha.wencodes.com/demo/)                         | GoCaptcha 在线演示                                                                              |
| [go-captcha-example](https://github.com/wenlng/go-captcha-example)         | Golang + 前端 + APP实例                                                                         |
| [go-captcha-assets](https://github.com/wenlng/go-captcha-assets)           | Golang 内嵌素材资源                                                                               |
| [go-captcha](https://github.com/wenlng/go-captcha)                         | Golang 验证码                                                                                  |
| [go-captcha-jslib](https://github.com/wenlng/go-captcha-jslib)             | Javascript 验证码                                                                              |
| [go-captcha-vue](https://github.com/wenlng/go-captcha-vue)                 | Vue 验证码                                                                                     |
| [go-captcha-react](https://github.com/wenlng/go-captcha-react)             | React 验证码                                                                                   |
| [go-captcha-angular](https://github.com/wenlng/go-captcha-angular)         | Angular 验证码                                                                                 |
| [go-captcha-svelte](https://github.com/wenlng/go-captcha-svelte)           | Svelte 验证码                                                                                  |
| [go-captcha-solid](https://github.com/wenlng/go-captcha-solid)             | Solid 验证码                                                                                   |
| [go-captcha-uni](https://github.com/wenlng/go-captcha-uni)                 | UniApp 验证码，兼容 Android、IOS、小程序、快应用等                                                          |
| [go-captcha-service](https://github.com/wenlng/go-captcha-service)         | GoCaptcha 服务，支持二进制、Docker镜像等方式部署，<br/> 提供 HTTP/GRPC 方式访问接口，<br/>可用单机模式和分布式（服务发现、负载均衡、动态配置等） |
| [go-captcha-service-sdk](https://github.com/wenlng/go-captcha-service-sdk) | GoCaptcha 服务SDK工具包，包含 HTTP/GRPC 请求服务接口，<br/>支持静态模式、服务发现、负载均衡                                |
| ...                                                                        |                                                                                             |

<br/>
<br/>

## 功能特性

- **多种验证码模式**：支持文本/图形点击、滑动、拖拽和旋转验证码。
- **双协议支持**：提供 RESTful HTTP 和 gRPC 接口。
- **服务发现**：支持 Etcd、Nacos、Zookeeper 和 Consul，实现分布式服务注册与发现。
- **分布式缓存**：支持 Memory、Redis、Etcd 和 Memcache，优化验证码数据存储。
- **分布式动态配置**：通过 Etcd、Nacos、Zookeeper 或 Consul 实时更新配置。
- **高可配置性**：支持自定义文本、字体、图片资源、验证码尺寸、生成规则等配置。
- **高性能**：基于 Go 的并发模型，适合高流量场景，同时结合分布式架构部署，确保服务处于高可用、高性能、高响应的状态。
- **跨平台**：支持二进制、命令行、PM2、Docker 和 Docker Compose 部署。

<br/>
<br/>

## 安装与部署
`GoCaptcha Service` 支持多种部署方式，包括单机部署（二进制、命令行、PM2、Docker）和分布式部署（结合服务发现和分布式缓存，分布式动态配置可选）。

### 前置条件
- 可选：Docker（用于容器化部署）
- 可选：服务发现/动态配置中间件（Etcd、Nacos、Zookeeper、Consul）
- 可选：缓存服务（Redis、Etcd、Memcache）
- 可选：Node.js 和 PM2（用于 PM2 部署）
- 可选：gRPC 客户端工具（如 `grpcurl`）

### 单机部署

#### 二进制方式

1. 从 [Github Releases](https://github.com/wenlng/go-captcha-service/releases) 或 [Gitee Releases](https://gitee.com/wenlng/go-captcha-service/releases) 发布页下载最新版本相对应平台的二进制可执行文件。
   
    ```bash
    ./go-captcha-service-[xxx]
    ```

2. 可选：配置应用：可复制仓库下的 config.json 应用配置和 gocaptcha.json 验证码配置文件放在同级目录下，在启动时指定配置文件。
    
   ```bash
    ./go-captcha-service-[xxx] -config config.json -gocaptcha-config gocaptcha.json
    ```

3. 访问 HTTP 接口（如 `http://localhost:8080/api/v1/get-data?id=click-default-ch`）或 gRPC 接口（`localhost:50051`）。


<br/>
<br/>

#### PM2 部署
PM2 是 Node.js 进程管理工具，可用于管理 Go 服务，提供进程守护和日志管理。
1. 安装 Node.js 和 PM2：

   ```bash
   npm install -g pm2
   ```

2. 创建 PM2 配置文件 `ecosystem.config.js`：

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
         CAPTCHA_HTTP_PORT: '8080',
         CAPTCHA_GRPC_PORT: '50051',
         CAPTCHA_CACHE_TYPE: 'memory'
       }
     }]
   };
   ```

3. 启动服务：

   ```bash
   pm2 start ecosystem.config.js
   ```

4. 查看日志和状态：

   ```bash
   pm2 logs go-captcha-service
   pm2 status
   ```

5. 设置开机自启：

   ```bash
   pm2 startup
   pm2 save
   ```

<br/>
<br/>

#### Docker 部署

1. 创建 `Dockerfile`：

   ```dockerfile
   FROM golang:1.18
   
   WORKDIR /app
   
   COPY . .
   
   RUN go mod download
   RUN go build -o go-captcha-service
   
   EXPOSE 8080 50051
   
   CMD ["./go-captcha-service"]
   ```

2. 构建镜像：

   ```bash
   docker build -t go-captcha-service:1.0.0 .
   ```

3. 运行容器，挂载配置文件：

   ```bash
   docker run -d -p 8080:8080 -p 50051:50051 \
     -v $(pwd)/config.json:/app/config.json \
     -v $(pwd)/gocaptcha.json:/app/gocaptcha.json \
     -v $(pwd)/gocaptcha:/app/gocaptcha \
     --name go-captcha-service go-captcha-service:latest
   ```

<br/>
<br/>


#### 使用官方 Docker 镜像

1. 拉取官方镜像：

   ```bash
   docker pull wenlng/go-captcha-service
   ```

2. 运行容器：

   ```bash
   docker run -d -p 8080:8080 -p 50051:50051 \
     -v $(pwd)/config.json:/app/config.json \
     -v $(pwd)/gocaptcha.json:/app/gocaptcha.json \
     -v $(pwd)/gocaptcha:/app/gocaptcha \
     --name go-captcha-service wenlng/go-captcha-service:latest
   ```

<br/>
<br/>

#### 配置分布式缓存

1. 在 `config.json` 中选择分布式缓存（如 Redis）：

   ```json
   {
     "cache_type": "redis",
     "redis_addrs": "localhost:6379",
     "cache_ttl": 1800,
     "cache_key_prefix": "GO_CAPTCHA_DATA:"
   }
   ```

2. 启动 Redis：

   ```bash
   docker run -d -p 6379:6379 --name redis redis:latest
   ```

<br/>

#### 分布式动态配置
注意：当开启分布式动态配置功能后，`config.json` 和 `gocaptcha.json` 会同时作用

1. 在 `config.json` 中启用动态配置，选择中间件（如 Etcd）：

   ```json
   {
     "enable_dynamic_config": true,
     "service_discovery": "etcd",
     "service_discovery_addrs": "localhost:2379"
   }
   ```

2. 启动 Etcd：

   ```bash
   docker run -d -p 8848:8848 --name etcd bitnami/etcd::latest
   ```

3. 例如在 gocaptcha.json 配置文件中，修改配置：
    
   ```json
   {
     "builder": {
       
     }
   }
   ```

4. 配置文件同步与拉取
* 服务在启动时会根据 `config_version` 版本决定推送与拉取，当本地版本大于远程（如 Etcd）的配置版本时会将本地配置推送覆盖，反之自动拉取更新到本地（非文件式更新）。
* 在服务启动后，动态配置管理器会实时监听远程（如 Etcd）的配置，当远程配置发生变更后，会摘取到本地进行版本比较，当大于本地版本时都会覆盖本地的配置。


<br/>
<br/>


#### 分布式服务发现
1. 在 `config.json` 中启用动态配置，选择中间件（如 Etcd）：

   ```json
   {
     "enable_service_discovery": true,
     "service_discovery": "etcd",
     "service_discovery_addrs": "localhost:2379"
   }
   ```

2. 启动 Etcd：

   ```bash
   docker run -d -p 8848:8848 --name etcd bitnami/etcd::latest
   ```
3. 服务注册与发现
* 服务在启动时会自动向（Etcd | xxx）的中心注册服务实例。
* 在服务启动后，同时会进行服务实例的变化监听，可参考在 [go-captcha-service-sdk](https://github.com/wenlng/go-captcha-service-sdk) 中的负载均衡应用。

<br/>
<br/>

#### Docker Compose 分布式部署

创建 `docker-compose.yml`，包含多个服务实例、Consul、Redis、ZooKeeper 和 Nacos：

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
    depends_on:
      - consul
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
    depends_on:
      - consul
      - redis
    restart: unless-stopped
       
  consul:
    image: consul:latest
    ports:
      - "8500:8500"
    command: agent -server -bootstrap -ui -client=0.0.0.0
    restart: unless-stopped

  redis:
    image: redis:latest
    ports:
      - "6379:6379"
    restart: unless-stopped
```

运行：

```bash
docker-compose up -d
```

<br/>
<br/>


## 配置说明

### 启动参数
注：启动参数与 `config.json` 文件中有相对应，注意名称格式（**推荐使用配置文件方式**）
* config：指定配置文件路径，默认 "config.json"。
* gocaptcha-config：指定 GoCaptcha 配置文件路径，默认 "gocaptcha.json"。
* service-name：设置服务名称。
* http-port：设置 HTTP 服务器端口。
* grpc-port：设置 gRPC 服务器端口。
* redis-addrs：设置 Redis 集群地址，逗号分隔。
* etcd-addrs：设置 etcd 地址，逗号分隔。
* memcache-addrs：设置 Memcached 地址，逗号分隔。
* cache-type：设置缓存类型，支持 redis、memory、etcd、memcache。
* cache-ttl：设置缓存 TTL，单位秒。
* cache-key-prefix：设置缓存键前缀，默认 "GO_CAPTCHA_DATA:"。
* service-discovery：设置服务发现类型，支持 etcd、zookeeper、consul、nacos。
* service-discovery-addrs：设置服务发现服务器地址，逗号分隔。
* service-discovery-ttl：设置服务发现注册存活时间，单位秒，默认 10。
* service-discovery-keep-alive：设置服务发现保活间隔，单位秒，默认 3。
* service-discovery-max-retries：设置服务发现操作最大重试次数，默认 3。
* service-discovery-base-retry-delay：设置服务发现重试基础延迟，单位毫秒，默认 3。
* service-discovery-username：设置服务发现认证用户名。
* service-discovery-password：设置服务发现认证密码。
* service-discovery-tls-server-name：设置服务发现 TLS 服务器名称。
* service-discovery-tls-address：设置服务发现 TLS 服务器地址。
* service-discovery-tls-cert-file：设置服务发现 TLS 证书文件路径。
* service-discovery-tls-key-file：设置服务发现 TLS 密钥文件路径。
* service-discovery-tls-ca-file：设置服务发现 TLS CA 文件路径。
* rate-limit-qps：设置速率限制 QPS。
* rate-limit-burst：设置速率限制突发量。
* api-keys：设置 API 密钥，逗号分隔。
* log-level：设置日志级别，支持 error、debug、warn、info。
* enable-service-discovery：启用服务发现，默认 false。
* enable-dynamic-config：启用动态配置，默认 false。
* health-check：运行健康检查并退出，默认 false。
* enable-cors：启用跨域资源共享，默认 false。


### 配置文件
服务使用两个配置文件：`config.json` 和 `gocaptcha.json`，分别定义服务运行参数和验证码生成的配置.

### config.json

`config.json` 定义服务的基础配置。

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
   "enable_service_discovery": false,
   "service_discovery": "etcd",
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

#### 参数说明

- `config_version` (整数)：配置文件版本号，用于分布式动态配置控制，默认 `1`。
- `service_name` (字符串)：服务名称，默认 `go-captcha-service`。
- `http_port` (字符串)：HTTP 端口，默认 `8080`。
- `grpc_port` (字符串)：gRPC 端口，默认 `50051`。
- `redis_addrs` (字符串)：Redis 地址，默认 `localhost:6379`。用于 `cache_type: redis`。
- `etcd_addrs` (字符串)：Etcd 地址，默认 `localhost:2379`。用于 `cache_type: etcd` 或 `service_discovery: etcd`.
- `memcache_addrs` (字符串)：Memcache 地址，默认 `localhost:11211`。用于 `cache_type: memcache`.
- `cache_type` (字符串)：缓存类型，默认 `memory`：
   - `memory`：内存缓存，适合单机部署。
   - `redis`：分布式键值存储，适合高可用场景。
   - `etcd`：分布式键值存储，适合与服务发现共用 Etcd。
   - `memcache`：高性能分布式缓存，适合高并发。
- `cache_ttl` (整数)：缓存有效期（秒），默认 `1800`.
- `cache_key_prefix` (字符串)：缓存键前缀，默认 `GO_CAPTCHA_DATA:`。
- `enable_dynamic_config` (布尔)：启用动态配置，默认 `false`。
- `enable_service_discovery` (布尔)：启用服务发现，默认 `false`。
- `service_discovery` (字符串)：服务发现类型，默认 `etcd`：
   - `etcd`：适合一致性要求高的场景。
   - `nacos`：适合云原生环境。
   - `zookeeper`：适合复杂分布式系统。
   - `consul`：轻量级，支持健康检查。
- `service_discovery_addrs` (字符串)：服务发现地址，如 Etcd 为 `localhost:2379`，Nacos 为 `localhost:8848`。
- `service_discovery_username` (字符串)：用户名，例如 Nacos 的默认用户名为`nacos`，默认空。
- `service_discovery_password` (字符串)：密码，例如 Nacos 的默认用户密码为`nacos`，默认空。
- `service_discovery_ttl` (整数)：服务注册租约时间（秒），默认 `10`。
- `service_discovery_keep_alive` (整数)：心跳间隔（秒），默认 `3`。
- `service_discovery_max_retries` (整数)：重试次数，默认 `3`。
- `service_discovery_base_retry_delay` (整数)：重试延迟（毫秒），默认 `500`。
- `service_discovery_tls_server_name` (字符串)：TLS 服务器名称，默认空。
- `service_discovery_tls_address` (字符串)：TLS 地址，默认空。
- `service_discovery_tls_cert_file` (字符串)：TLS 证书文件，默认空。
- `service_discovery_tls_key_file` (字符串)：TLS 密钥文件，默认空。
- `service_discovery_tls_ca_file` (字符串)：TLS CA 证书文件，默认空。
- `rate_limit_qps` (整数)：API 每秒请求限流，默认 `1000`。
- `rate_limit_burst` (整数)：API 限流突发容量，默认 `1000`。
- `enable_cors` (布尔)：启用 CORS，默认 `true`。
- `log_level` (字符串)：日志级别（`debug`、`info`、`warn`、`error`），默认 `info`。
- `api_keys` (字符串数组)：API 认证密钥，默认包含示例密钥。

### gocaptcha.json

`gocaptcha.json` 定义验证码的资源和生成配置。

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
        // ...
      },
      "click-default-en": {
        "version": "0.0.1",
        "language": "english",
        // ...
      },
      "click-dark-en": {
        "version": "0.0.1",
        "language": "english",
        // ...
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

##### 顶级字段

- `config_version` (整数)：配置文件版本号，用于分布动态配置管理，默认 `1`。


##### resources 字段

- `version` (字符串)：资源配置版本号，用于控制重新创建新的验证码实例，默认 `0.0.1`。
- `char.languages.chinese` (字符串数组)：中文字符集，用于点击验证码的文本内容，默认空（默认取内置的资源）。
- `char.languages.english` (字符串数组)：英文字符集，默认空（默认取内置的资源）。
- `font.type` (字符串)：字体加载方式，固定为 `load`（从文件加载）。
- `font.file_dir` (字符串)：字体文件目录，默认 `./gocaptcha/fonts/`。
- `font.file_maps` (对象)：字体文件映射，键为字体名称，值为文件名。
    - 示例：`"yrdzst_bold": "yrdzst-bold.ttf"` 表示使用 `yrdzst-bold.ttf` 字体。
- `shape_image.type` (字符串)：形状图片加载方式，固定为 `load`。
- `shape_image.file_dir` (字符串)：形状图片目录，默认 `./gocaptcha/shape_images/`。
- `shape_image.file_maps` (对象)：形状图片映射。
    - 示例：`"shape_01": "shape_01.png"` 表示使用 `shape_01.png` 作为形状。
- `master_image.type` (字符串)：主图片加载方式，固定为 `load`。
- `master_image.file_dir` (字符串)：主图片目录，默认 `./gocaptcha/master_images/`。
- `master_image.file_maps` (对象)：主图片映射。
    - 示例：`"image_01": "image_01.jpg"` 表示使用 `image_01.jpg` 作为背景。
- `thumb_image.type` (字符串)：缩略图加载方式，固定为 `load`。
- `thumb_image.file_dir` (字符串)：缩略图目录，默认 `./gocaptcha/thumb_images/`。
- `thumb_image.file_maps` (对象)：缩略图映射，默认空。
- `tile_image.type` (字符串)：拼图图片加载方式，固定为 `load`。
- `tile_image.file_dir` (字符串)：拼图图片目录，默认 `./gocaptcha/tile_images/`。
- `tile_image.file_maps` (对象)：拼图图片映射。
    - 示例：`"tile_01": "tile_01.png"`。
- `tile_image.file_maps_02` (对象)：拼图蒙版映射。
    - 示例：`"tile_mask_01": "tile_mask_01.png"`。
- `tile_image.file_maps_03` (对象)：拼图阴影映射。
    - 示例：`"tile_shadow_01": "tile_shadow_01.png"`。

<br/>

##### builder 字段

定义验证码生成样式，包含点击、形状点击、滑动、拖拽和旋转验证码的配置。


###### click_config_maps

定义文本点击验证码的配置，支持中英文和明暗主题，key为ID，在请求时传递，例如：`api/v1/get-data?id=click-default-ch`。

- `click-default-ch` (对象)：中文默认主题配置。
    - `version` (字符串)：配置版本号，用于控制重新创建新的验证码实例，默认 `0.0.1`。
    - `language` (字符串)：语言，可配置 `char.languages` 中定义的语言名称，例如中文： `chinese`。
    - `master` (对象)：主验证码图片配置。
        - `image_size.width` (整数)：主图片宽度，默认 `300`。
        - `image_size.height` (整数)：主图片高度，默认 `200`。
        - `range_length.min` (整数)：验证码点数最小值，默认 `6`。
        - `range_length.max` (整数)：验证码点数最大值，默认 `7`。
        - `range_angles` (对象数组)：文本旋转角度范围（度）。
            - 示例：`{"min": 20, "max": 35}` 表示角度范围 20°-35°。
        - `range_size.min` (整数)：文本大小最小值（像素），默认 `26`。
        - `range_size.max` (整数)：文本大小最大值，默认 `32`。
        - `range_colors` (字符串数组)：文本颜色列表（十六进制）。
            - 示例：`"#fde98e"`。
        - `display_shadow` (布尔)：是否显示文本阴影，默认 `true`。
        - `shadow_color` (字符串)：阴影颜色，默认 `#101010`。
        - `shadow_point.x` (整数)：阴影偏移 X 坐标，默认 `-1`（自动计算）。
        - `shadow_point.y` (整数)：阴影偏移 Y 坐标，默认 `-1`。
        - `image_alpha` (浮点数)：图片透明度（0-1），默认 `1`。
        - `use_shape_original_color` (布尔)：是否使用形状原始颜色，默认 `true`。
    - `thumb` (对象)：缩略图（提示文本）配置。
        - `image_size.width` (整数)：缩略图宽度，默认 `150`。
        - `image_size.height` (整数)：缩略图高度，默认 `40`。
        - `range_verify_length.min` (整数)：验证点数最小值，默认 `2`。
        - `range_verify_length.max` (整数)：验证点数最大值，默认 `4`。
        - `disabled_range_verify_length` (布尔)：是否禁用验证点数限制，默认 `false`。
        - `range_text_size.min` (整数)：文本大小最小值，默认 `22`。
        - `range_text_size.max` (整数)：文本大小最大值，默认 `28`。
        - `range_text_colors` (字符串数组)：文本颜色列表。
        - `range_background_colors` (字符串数组)：背景颜色列表。
        - `is_non_deform_ability` (布尔)：是否禁用变形效果，默认 `false`。
        - `background_distort` (整数)：背景扭曲程度，默认 `4`。
        - `background_distort_alpha` (浮点数)：背景扭曲透明度，默认 `1`。
        - `background_circles_num` (整数)：背景圆形干扰点数量，默认 `24`。
        - `background_slim_line_num` (整数)：背景细线干扰数量，默认 `2`。
    
- `click-dark-ch` (对象)：中文暗色主题配置，参数与 `click-default-ch` 类似，区别在于 `thumb.range_text_colors` 使用更亮的颜色以适配暗色背景。

- `click-default-en` (对象)：英文默认主题配置，`language: english` 、`master.range_size` 和 `thumb.range_text_size` 更大（`34-48`），适配英文字符。

- `click-dark-en` (对象)：英文暗色主题配置，类似 `click-dark-ch`, 注意区别字段 `language: english`。

<br/>

###### click_shape_config_maps

定义形状点击验证码的配置。

- `click-shape-default` (对象)：默认形状点击配置，参数与 `click_config_maps` 的 `master` 和 `thumb` 类似，但针对形状图片而非文本。

<br/>

###### slide_config_maps

定义滑动验证码配置。

- `slide-default` (对象)：
    - `version` (字符串)：配置版本号，用于控制重新创建新的验证码实例，默认 `0.0.1`。
    - `master` (对象)：主验证码图片配置。
        - `image_size.width` (整数)：主图片宽度，默认 `300`。
        - `image_size.height` (整数)：主图片高度，默认 `200`。
        - `image_alpha` (浮点数)：图片透明度（0-1），默认 `1`。
    - `thumb` (对象)：滑块配置。
        - `range_graph_size.min` (整数)：滑块图形大小最小值（像素），默认 `60`。
        - `range_graph_size.max` (整数)：滑块图形大小最大值，默认 `70`。
        - `range_graph_angles` (对象数组)：滑块图形旋转角度范围（度）。
            - 示例：`{"min": 20, "max": 35}`。
        - `generate_graph_number` (整数)：生成滑块图形数量，默认 `1`。
        - `enable_graph_vertical_random` (布尔)：是否启用垂直方向随机偏移，默认 `false`。
        - `range_dead_zone_directions` (字符串数组)：滑块禁区方向，默认 `["left", "right"]`。

<br/>

###### drag_config_maps

定义拖拽验证码配置。

- `drag-default` (对象)：
    - `version` (字符串)：配置版本号，用于控制重新创建新的验证码实例，默认 `0.0.1`。
    - `master` (对象)：主验证码图片配置。
        - `image_size.width` (整数)：主图片宽度，默认 `300`。
        - `image_size.height` (整数)：主图片高度，默认 `200`。
        - `image_alpha` (浮点数)：图片透明度（0-1），默认 `1`。
    - `thumb` (对象)：拖拽图形配置。
        - `range_graph_size.min` (整数)：拖拽图形大小最小值（像素），默认 `60`。
        - `range_graph_size.max` (整数)：拖拽图形大小最大值，默认 `70`。
        - `range_graph_angles` (对象数组)：拖拽图形旋转角度范围（度）。
            - 示例：`{"min": 0, "max": 0}` 表示无旋转。
        - `generate_graph_number` (整数)：生成拖拽图形数量，默认 `2`。
        - `enable_graph_vertical_random` (布尔)：是否启用垂直方向随机偏移，默认 `true`。
        - `range_dead_zone_directions` (字符串数组)：拖拽禁区方向，默认 `["left", "right", "top", "bottom"]`。

<br/>

###### rotate_config_maps

定义旋转验证码配置。

- `rotate-default` (对象)：
    - `version` (字符串)：配置版本号，用于控制重新创建新的验证码实例，默认 `0.0.1`。
    - `master` (对象)：主验证码图片配置。
        - `image_square_size` (整数)：主图片正方形边长（像素），默认 `220`。
    - `thumb` (对象)：旋转图形配置。
        - `range_angles` (对象数组)：旋转角度范围（度）。
            - 示例：`{"min": 30, "max": 330}` 表示旋转范围 30°-330°。
        - `range_image_square_sizes` (整数数组)：旋转图片正方形边长列表，默认 `[140, 150, 160, 170]`。
        - `image_alpha` (浮点数)：图片透明度（0-1），默认 `1`。



<br/>
<br/>


### 配置热重载说明
`gocaptcha.json` 热重载以每个配置顶的 version 字段决定是否生效。

`config.json` 热重载有效的字段如下：
* redis_addrs
* etcd_addrs
* memcache_addrs
* cache_type
* cache_ttl
* cache_key_prefix
* api_keys
* log_level
* rate_limit_qps
* rate_limit_burst



### 测试：

- 验证码生成：验证图片、形状和密钥有效性。
- 验证逻辑：测试不同输入的处理。
- 服务发现：模拟 Etcd/Nacos/Zookeeper/Consul。
- 缓存：测试 Memory/Redis/Etcd/Memcache。
- 动态配置：测试 Nacos 配置更新。


### 压力测试

测试 HTTP 接口：

```bash
wrk -t12 -c400 -d30s http://127.0.0.1:8080/api/v1/get-data?id=click-default-ch
```

测试 gRPC 接口：

```bash
grpcurl -plaintext -d '{"id":"click-default-ch"}' localhost:50051 gocaptcha.GoCaptchaService/GetData
```
