# ITO-Deposit 服务注册与负载均衡指南

## 📋 项目现状分析

### ✅ 已具备的条件
- **gRPC服务已配置** - `internal/server/grpc.go` 已实现完整的gRPC服务器
- **依赖已安装** - `go.mod` 中已包含etcd相关依赖：
  - `github.com/go-kratos/kratos/contrib/registry/etcd/v2`
  - `go.etcd.io/etcd/client/v3`
- **服务架构完整** - 使用Kratos框架，支持HTTP和gRPC双协议
- **配置系统完善** - 基于protobuf的配置管理

### 🔧 当前服务配置
```yaml
# configs/config.yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000  # gRPC服务端口
    timeout: 1s
  jwt:
    authkey: a9999
```

### 🎯 注册的服务列表
- `GreeterServer` - 示例服务
- `UserServer` - 用户管理服务
- `HomeServer` - 首页服务
- `DepositServer` - 寄存服务（核心业务）
- `OrderServer` - 订单服务
- `AdminServer` - 管理后台服务

---

## 🚀 etcd服务注册实现方案

### 1. 配置文件扩展

#### 1.1 更新protobuf配置定义
```protobuf
// internal/conf/conf.proto 添加etcd配置
message Server {
  message HTTP { ... }
  message GRPC { ... }
  message Jwt { ... }
  
  // 新增etcd配置
  message Etcd {
    repeated string endpoints = 1;
    google.protobuf.Duration timeout = 2;
    google.protobuf.Duration dial_timeout = 3;
    string username = 4;
    string password = 5;
  }
  
  HTTP http = 1;
  GRPC grpc = 2;
  Jwt jwt = 3;
  Etcd etcd = 4;  // 新增字段
}
```

#### 1.2 更新YAML配置文件
```yaml
# configs/config.yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 1s
  jwt:
    authkey: a9999
  # 新增etcd配置
  etcd:
    endpoints:
      - "127.0.0.1:2379"
      - "127.0.0.1:2380"  # 集群模式
      - "127.0.0.1:2381"
    timeout: 3s
    dial_timeout: 5s
    username: ""          # 如需认证
    password: ""          # 如需认证
```

### 2. 主程序修改

#### 2.1 更新 `cmd/ito-deposit/main.go`
```go
package main

import (
	"flag"
	"os"

	"ito-deposit/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	// 新增etcd相关导入
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"

	_ "go.uber.org/automaxprocs"
)

// 修改newApp函数，支持服务注册
func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, rr registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name("ito-deposit"),           // 服务名称
		kratos.Version(Version),
		kratos.Metadata(map[string]string{
			"env": "production",               // 环境标识
			"region": "cn-hangzhou",          // 区域标识
		}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
		kratos.Registrar(rr),                 // 注册服务发现
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", "ito-deposit",
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	c := config.New(config.WithSource(file.NewSource(flagconf)))
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 创建etcd客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   bc.Server.Etcd.Endpoints,
		DialTimeout: bc.Server.Etcd.DialTimeout.AsDuration(),
		Username:    bc.Server.Etcd.Username,
		Password:    bc.Server.Etcd.Password,
	})
	if err != nil {
		log.Fatalf("Failed to create etcd client: %v", err)
	}
	defer etcdClient.Close()

	// 创建etcd注册器
	registry := etcd.New(etcdClient)

	// 修改wireApp调用，传入registry
	app, cleanup, err := wireApp(bc.Server, bc.Data, registry, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 启动服务
	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

#### 2.2 更新 `cmd/ito-deposit/wire.go`
```go
//go:build wireinject
// +build wireinject

package main

import (
	"ito-deposit/internal/biz"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/data"
	"ito-deposit/internal/server"
	"ito-deposit/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, registry.Registrar, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
```

### 3. 服务健康检查

#### 3.1 添加健康检查接口
```go
// internal/service/health.go
package service

import (
	"context"
	
	pb "ito-deposit/api/helloworld/v1"
)

type HealthService struct {
	pb.UnimplementedHealthServer
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	// 检查数据库连接
	// 检查Redis连接
	// 检查其他依赖服务
	
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}
```

---

## ⚖️ 负载均衡配置方案

### 1. 多实例部署

#### 1.1 启动多个服务实例
```bash
# 实例1 - 默认端口
go run ./cmd/ito-deposit

# 实例2 - 自定义端口
go run ./cmd/ito-deposit -conf ./configs/config.yaml \
  --server.grpc.addr=0.0.0.0:9001 \
  --server.http.addr=0.0.0.0:8001

# 实例3 - 不同机器
go run ./cmd/ito-deposit -conf ./configs/config.yaml \
  --server.grpc.addr=192.168.1.100:9000 \
  --server.http.addr=192.168.1.100:8000
```

#### 1.2 Docker容器化部署
```dockerfile
# Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ito-deposit ./cmd/ito-deposit

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ito-deposit .
COPY --from=builder /app/configs ./configs
CMD ["./ito-deposit"]
```

```yaml
# docker-compose.yml
version: '3.8'
services:
  ito-deposit-1:
    build: .
    ports:
      - "8001:8000"
      - "9001:9000"
    environment:
      - SERVER_GRPC_ADDR=0.0.0.0:9000
      - SERVER_HTTP_ADDR=0.0.0.0:8000
    depends_on:
      - etcd

  ito-deposit-2:
    build: .
    ports:
      - "8002:8000"
      - "9002:9000"
    environment:
      - SERVER_GRPC_ADDR=0.0.0.0:9000
      - SERVER_HTTP_ADDR=0.0.0.0:8000
    depends_on:
      - etcd

  etcd:
    image: quay.io/coreos/etcd:v3.5.0
    ports:
      - "2379:2379"
    command:
      - etcd
      - --advertise-client-urls=http://0.0.0.0:2379
      - --listen-client-urls=http://0.0.0.0:2379
```

### 2. 客户端负载均衡

#### 2.1 gRPC客户端配置
```go
// internal/client/deposit_client.go
package client

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	clientv3 "go.etcd.io/etcd/client/v3"
	
	pb "ito-deposit/api/helloworld/v1"
)

type DepositClient struct {
	client pb.DepositClient
}

func NewDepositClient() (*DepositClient, error) {
	// 创建etcd客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	// 创建服务发现
	discovery := etcd.New(etcdClient)

	// 创建gRPC连接，启用负载均衡
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///ito-deposit"),  // 服务发现地址
		grpc.WithDiscovery(discovery),                   // 服务发现实例
		grpc.WithBalancerName("round_robin"),           // 轮询负载均衡
		grpc.WithTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &DepositClient{
		client: pb.NewDepositClient(conn),
	}, nil
}

func (c *DepositClient) CreateDeposit(ctx context.Context, req *pb.CreateDepositRequest) (*pb.CreateDepositReply, error) {
	return c.client.CreateDeposit(ctx, req)
}
```

#### 2.2 HTTP客户端配置
```go
// internal/client/http_client.go
package client

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewHTTPClient() (*http.Client, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	discovery := etcd.New(etcdClient)

	return http.NewClient(
		context.Background(),
		http.WithEndpoint("discovery:///ito-deposit"),
		http.WithDiscovery(discovery),
		http.WithBalancer("round_robin"),
		http.WithTimeout(10*time.Second),
	)
}
```

---

## 🔧 高级配置选项

### 1. 服务注册配置
```go
// 自定义服务注册信息
registry := etcd.New(etcdClient, etcd.RegisterTTL(30*time.Second))

app := kratos.New(
	kratos.Name("ito-deposit"),
	kratos.Version("v1.0.0"),
	kratos.Metadata(map[string]string{
		"weight": "100",                    // 权重
		"region": "cn-hangzhou",           // 区域
		"zone": "zone-a",                  // 可用区
		"cluster": "production",           // 集群标识
	}),
	kratos.Registrar(registry),
)
```

### 2. 负载均衡策略
```go
// 支持的负载均衡算法
grpc.WithBalancerName("round_robin")     // 轮询（默认）
grpc.WithBalancerName("pick_first")      // 选择第一个
grpc.WithBalancerName("weighted_round_robin") // 加权轮询
```

### 3. 故障转移配置
```go
// gRPC重试配置
import "github.com/grpc-ecosystem/go-grpc-middleware/retry"

conn, err := grpc.DialInsecure(
	context.Background(),
	grpc.WithEndpoint("discovery:///ito-deposit"),
	grpc.WithDiscovery(discovery),
	grpc.WithUnaryInterceptor(
		grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithMax(3),                    // 最大重试3次
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100*time.Millisecond)),
		),
	),
)
```

---

## 📊 监控与运维

### 1. 服务监控
```go
// 添加Prometheus监控
import "github.com/go-kratos/kratos/v2/middleware/metrics"

opts = append(opts, grpc.Middleware(
	metrics.Server(),  // 服务端监控
))
```

### 2. 链路追踪
```go
// 添加Jaeger追踪
import "github.com/go-kratos/kratos/v2/middleware/tracing"

opts = append(opts, grpc.Middleware(
	tracing.Server(),  // 链路追踪
))
```

### 3. 日志聚合
```go
// 结构化日志
logger := log.With(log.NewStdLogger(os.Stdout),
	"service.name", "ito-deposit",
	"service.version", Version,
	"service.id", id,
	"trace.id", tracing.TraceID(),
)
```

---

## 🚨 注意事项

### 1. 生产环境部署
- **etcd集群** - 建议3节点或5节点集群，避免单点故障
- **网络分区** - 配置合适的超时时间和重试策略
- **安全认证** - 启用etcd的TLS和用户认证
- **资源限制** - 设置合适的连接池大小和超时时间

### 2. 性能优化
- **连接复用** - 使用连接池，避免频繁创建连接
- **批量操作** - 对于高频调用，考虑批量处理
- **缓存策略** - 在客户端缓存服务发现结果

### 3. 故障处理
- **熔断器** - 实现熔断机制，防止雪崩效应
- **降级策略** - 准备服务降级方案
- **健康检查** - 实现完善的健康检查机制

---

## 📝 实施计划

### 阶段1：基础配置（1-2天）
1. 更新protobuf配置定义
2. 修改配置文件
3. 更新主程序和wire配置

### 阶段2：服务注册（2-3天）
1. 集成etcd客户端
2. 实现服务注册逻辑
3. 添加健康检查

### 阶段3：负载均衡（3-4天）
1. 实现客户端负载均衡
2. 测试多实例部署
3. 验证故障转移

### 阶段4：监控运维（2-3天）
1. 添加监控指标
2. 实现链路追踪
3. 完善日志系统

---

## 🔗 相关资源

- [Kratos官方文档](https://go-kratos.dev/)
- [etcd官方文档](https://etcd.io/docs/)
- [gRPC负载均衡](https://grpc.io/blog/grpc-load-balancing/)
- [服务发现最佳实践](https://microservices.io/patterns/service-registry.html)

---

**备注**: 此文档为实施指南，当前项目已具备所有必要依赖，可按需实施。建议在开发环境先验证，再部署到生产环境。