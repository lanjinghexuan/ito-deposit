# ito-deposit

本项目基于 [Kratos](https://go-kratos.dev/) 微服务框架，采用领域驱动设计，支持 gRPC/HTTP 协议，适用于存款、订单、用户等业务场景。

---

## 目录结构

```
ito-deposit/
├── api/                # Protobuf API 定义及生成代码
│   └── helloworld/v1/  # 各业务 proto 及 pb.go 文件
├── cmd/ito-deposit/    # 服务启动入口（main.go、wire.go）
├── configs/            # 配置文件（config.yaml）
├── internal/           # 内部业务逻辑
│   ├── biz/            # 领域业务对象与用例
│   ├── conf/           # 配置结构体
│   ├── data/           # 数据访问层（如数据库、缓存等）
│   ├── server/         # 服务注册与启动（HTTP/gRPC）
│   └── service/        # gRPC/HTTP 服务实现
├── third_party/        # 三方 proto 依赖
├── Dockerfile          # Docker 构建文件
├── go.mod/go.sum       # Go 依赖管理
├── Makefile            # 常用自动化命令
└── README.md           # 项目说明
```

---

## 主要功能模块

- 用户服务（注册、登录、发送短信等）
- 存款服务
- 订单服务
- 首页服务
- 管理后台服务
- 示例 Greeter 服务

---

## 快速开始

### 1. 安装依赖

```bash
make init
```

### 2. 生成代码

```bash
make api      # 生成 API 相关代码（pb.go, http, grpc, swagger 等）
make all      # 生成全部代码（包括 wire 依赖注入等）
```

### 3. 启动服务（支持热更新）

在项目根目录下运行：

```bash
kratos run
```

如需指定配置文件：

```bash
kratos run -conf ./configs/config.yaml
```

或手动编译运行：

```bash
go build -o bin/ito-deposit ./cmd/ito-deposit
./bin/ito-deposit -conf ./configs/config.yaml
```

---

## 使用 Kratos 命令生成文件

### 安装 Kratos 工具

```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
go get github.com/google/wire/cmd/wire
```

### 创建新服务模板

```bash
kratos new ito-deposit
```

### 添加 proto 文件

```bash
kratos proto add api/helloworld/v1/deposit.proto
```

### 生成 proto 相关代码

```bash
kratos proto client api/helloworld/v1/deposit.proto
kratos proto server api/helloworld/v1/deposit.proto -t internal/service
```

### 生成依赖注入代码（wire）

```bash
cd cmd/ito-deposit
wire
```

---

## JWT 白名单接口添加说明

本项目默认对 HTTP 接口启用了 JWT 认证。如果有接口不需要认证（如登录、注册、获取 token 等），需将其添加到白名单。

**操作步骤：**

1. 打开 `internal/server/http.go` 文件。
2. 找到 `NewWhiteListMatcher` 方法。
3. 在 `whiteList` 变量中，按如下格式添加不需要认证的接口路径：

```go
whiteList["/api.helloworld.v1.Deposit/ReturnToken"] = struct{}{}
whiteList["/shop.interface.v1.ShopInterface/Register"] = struct{}{}
```

接口路径格式为：`/包名.服务名/方法名`，可在 proto 文件中查找。

4. 保存文件，重新编译并启动服务。

---

## 常用 Makefile 命令

| 命令         | 说明                       |
| ------------ | -------------------------- |
| make init    | 初始化/更新依赖            |
| make api     | 生成 API 相关代码          |
| make all     | 生成全部代码               |
| make build   | 编译可执行文件到 bin/      |

---

## 依赖工具

- Go 1.21+
- Kratos v2
- Wire
- Protobuf
- Docker（可选）

---

## 贡献

欢迎提交 issue 和 PR！

---

如需详细了解 Kratos 命令用法，请参考 [Kratos 官方文档](https://go-kratos.dev/docs/getting-started/)。  
如需自定义项目结构或添加新业务 proto，只需重复相关步骤即可。

---

如需进一步定制（如添加 API 示例、数据库说明、部署文档等），请告知你的具体需求！
