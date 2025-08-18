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

方案 1：使用 Chocolatey（推荐）
Chocolatey 是 Windows 的包管理器，可快速安装make。

安装 Chocolatey
以管理员身份打开 PowerShell，执行：
powershell
```
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```



安装 make
在 PowerShell 中执行：
powershell
```
choco install make
```



验证安装
重启 PowerShell 或 Git Bash，输入：
bash
```
make --version
```


如果要修改conf.proto 文件
修改后根目录输入命令 
```
make config
```

官方地址
https://go-kratos.dev/docs/component/config#%E4%BD%BF%E7%94%A8-1


//git 解决冲突
```
git fetch origin
git merge main
// 解决本地冲突后  
git add .
git commit -m '解决冲突'
git push
git checkout main
git merge wxy   // 自己的分支
git  push 

```

---

## JWT 白名单接口添加说明

本项目默认对 HTTP 接口启用了 JWT 认证。如果有接口不需要认证（如登录、注册、获取 token 等），需将其添加到白名单。

### 步骤：
1. 打开 `internal/server/http.go` 文件。
2. 找到 `NewWhiteListMatcher` 方法。
3. 在 `whiteList` 变量中，按如下格式添加不需要认证的接口路径：

```go
whiteList["/api.helloworld.v1.Deposit/ReturnToken"] = struct{}{}
whiteList["/shop.interface.v1.ShopInterface/Register"] = struct{}{}
// 继续添加...
```

接口路径格式为：`/包名.服务名/方法名`，可在 proto 文件中查找。

4. 保存文件，重新编译并启动服务。

---

## 启动服务

在项目根目录下，使用 Kratos 命令一键启动服务（支持热更新）：

```bash
kratos run
```

如需指定配置文件，可加参数：

```bash
kratos run -conf ./configs/config.yaml
```

---
项目必须启动etcd
```bash
etcdctl --endpoints=http://localhost:2379 member list
```

## 定时任务方式实现
项目中定时任务是采用官方的依赖注入形式进行实现的。
主要代码在  
`internal/basic/pkg/job/scheduler.go`
`cmd/ito-deposit/wire.go`
`cmd/ito-deposit/main.go`

scheduler.go 代码主要是注入代码，并且定义定时任务
wire是对代码进行依赖注入
main是启动定时任务