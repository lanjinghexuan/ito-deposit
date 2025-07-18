# Kratos 微服务项目使用说明

本项目基于 [Kratos](https://go-kratos.dev/) 微服务框架，采用多服务（user、order、deposit、home）结构，适合企业级 Go 微服务开发。

## 目录结构

```
ito/
  app/
    user/      # 用户服务
    order/     # 订单服务
    deposit/   # 充值服务
    home/      # 首页服务
  api/         # proto 定义
  internal/    # 公共内部实现
  third_party/ # 三方 proto
  configs/     # 公共配置
```

## 环境准备

- Go 1.18+
- Kratos 工具链
- wire 依赖注入工具
- protoc 及相关插件

### 安装依赖

```bash
go mod tidy
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
go install github.com/google/wire/cmd/wire@latest
```

## 服务开发与运行

以 user 服务为例：

### 1. 运行服务（开发模式）

```bash
cd app/user
kratos run
```
- `kratos run` 会自动查找 `cmd/user` 目录下的 main.go 并运行。
- 你也可以用 `go run ./cmd/user` 直接运行。

### 2. 依赖注入（wire）

```bash
cd app/user/cmd/user
wire
```

### 3. 编译与运行（生产模式）

```bash
make build -C app/user
./app/user/cmd/user/user -conf ./app/user/configs
```

### 4. Docker 构建与运行

```bash
docker build -t kratos-user ./app/user
docker run --rm -p 8000:8000 -p 9000:9000 -v $(pwd)/app/user/configs:/data/conf kratos-user
```

## API 代码与文档生成

### 1. proto 文件操作与代码生成

#### 新增 proto 文件

```bash
# 添加 proto 文件模板（会自动生成基础内容）
kratos proto add api/user/v1/user.proto
```

#### 生成 client 代码

```bash
# 根据 proto 文件生成 client 代码
kratos proto client api/user/v1/user.proto
```

#### 生成 server 代码

```bash
# 根据 proto 文件生成 server 代码（service 层实现模板）
kratos proto server api/user/v1/user.proto -t internal/service
```

### 2. 生成 API 相关代码（Makefile 方式）

```bash
# 编辑/新增 proto 文件
vim api/user/v1/user.proto

# 生成 API 相关代码
make api -C app/user
```

### 3. 生成 Swagger/OpenAPI 文档

```bash
# 生成 swagger.json
make swagger -C app/user
# 或
kratos proto swagger api/user/v1/user.proto

# 生成 openapi 文档
kratos proto openapi api/user/v1/user.proto
```

生成的文档一般在 `api/` 目录下。

## 常用 Makefile 命令

- `make init`：初始化依赖
- `make api`：生成 API 相关代码
- `make all`：生成全部代码
- `make build`：编译服务
- `make swagger`：生成 swagger 文档

## 其他说明

- 配置文件位于各服务的 `configs/` 目录下
- 各服务的入口在 `cmd/<service>` 目录下
- 依赖注入入口为 `wire.go`，自动生成实现为 `wire_gen.go`
- API 定义统一放在 `api/` 目录
- 目前示例文件在 `app/user`
- 数据库初始化文件在 `app/user/internal/data/data.go`
- 数据库使用示例在 `app/user/internal/service/user.go`

---

如需详细定制或补充，请联系项目维护者。

