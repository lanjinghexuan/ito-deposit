# ito-deposit

本项目基于 [Kratos](https://go-kratos.dev/) 微服务框架，采用领域驱动设计，支持 gRPC/HTTP 协议，适用于存款、订单、用户等业务场景。

## 项目结构说明

```
ito-deposit/
├── api/                # Protobuf API 定义及生成代码
│   └── helloworld/v1/  # 各业务 proto 及 pb.go 文件
├── cmd/ito-deposit/    # 服务启动入口（main.go）
├── configs/            # 配置文件（config.yaml）
├── internal/           # 内部业务逻辑
│   ├── biz/            # 领域业务对象与用例
│   ├── conf/           # 配置结构体
│   ├── data/           # 数据访问层
│   ├── server/         # 服务注册与启动
│   └── service/        # gRPC/HTTP 服务实现
├── third_party/        # 三方 proto 依赖
├── Dockerfile          # Docker 构建文件
├── go.mod/go.sum       # Go 依赖管理
├── Makefile            # 常用自动化命令
└── README.md           # 项目说明
```

## 使用 Kratos 命令生成文件

### 1. 安装 Kratos 工具

```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
go get github.com/google/wire/cmd/wire
```

### 2. 创建新服务模板

```bash
kratos new ito-deposit
```

### 3. 添加 proto 文件

```bash
kratos proto add api/helloworld/v1/deposit.proto
```

### 4. 生成 proto 相关代码

```bash
kratos proto client api/helloworld/v1/deposit.proto
kratos proto server api/helloworld/v1/deposit.proto -t internal/service
```

### 5. 生成依赖注入代码（wire）

```bash
cd cmd/ito-deposit
wire
```

### 6. 常用 Makefile 命令

| 命令         | 说明                       |
| ------------ | -------------------------- |
| make init    | 初始化/更新依赖            |
| make api     | 生成 API 相关代码          |
| make all     | 生成全部代码               |

---

如需详细了解 Kratos 命令用法，请参考 [Kratos 官方文档](https://go-kratos.dev/docs/getting-started/)。  
如需自定义项目结构或添加新业务 proto，只需重复第 3-4 步即可。

