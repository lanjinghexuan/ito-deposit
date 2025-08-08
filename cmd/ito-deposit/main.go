package main

import (
	"context"
	"flag"
	etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2" // etcd注册中心插件
	"github.com/go-kratos/kratos/v2/registry"                   // kratos服务注册相关接口
	clientv3 "go.etcd.io/etcd/client/v3"                        // etcd客户端v3版本
	"ito-deposit/internal/basic/pkg/job"
	"ito-deposit/internal/conf"
	"os"

	"ito-deposit/internal/basic/pkg"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"             // kratos配置加载包
	"github.com/go-kratos/kratos/v2/config/file"        // 文件配置源
	"github.com/go-kratos/kratos/v2/log"                // kratos日志包
	"github.com/go-kratos/kratos/v2/middleware/tracing" // 分布式追踪中间件
	"github.com/go-kratos/kratos/v2/transport/grpc"     // grpc服务支持
	"github.com/go-kratos/kratos/v2/transport/http"     // http服务支持

	"go.uber.org/zap"

	_ "go.uber.org/automaxprocs" // 自动调整GOMAXPROCS与CPU配合
	_ "net/http/pprof"
)

// 通过 -ldflags 方式注入的版本信息
var (
	Name     string // 服务名称
	Version  string // 服务版本号
	flagconf string // 配置文件路径参数

	id, _ = os.Hostname() // 主机名，作为服务实例ID部分
)

// 初始化命令行参数，允许通过 -conf 参数指定配置路径
func init() {
	flag.StringVar(&flagconf, "conf", "./configs/config.yaml", "config path, eg: -conf config.yaml")
}

// 构造 kratos 应用实例
func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, reg registry.Registrar, c *conf.Server, scheduler *job.Scheduler) *kratos.App {
	Name = c.Etcd.Name    // 读取配置中的服务名赋值
	id = id + c.Grpc.Addr // 用 主机名+grpc地址 作为服务唯一id
	return kratos.New(
		kratos.ID(id),                        // 服务实例ID
		kratos.Name(Name),                    // 服务名
		kratos.Version(Version),              // 服务版本
		kratos.Metadata(map[string]string{}), // 可添加元信息
		kratos.Logger(logger),                // 日志组件
		kratos.Server(
			gs, // grpc服务
			hs, // http服务
		),
		kratos.Registrar(reg), // 服务注册中心
		kratos.BeforeStart(func(ctx context.Context) error {
			scheduler.Start()
			return nil
		}),
		kratos.BeforeStop(func(ctx context.Context) error {
			return scheduler.Stop(ctx)
		}),
	)
}

func main() {
	flag.Parse() // 解析命令行参数

	// 创建一个标准输出日志组件，带时间戳、调用者、服务id等字段
	flag.Parse()

	// 初始化自定义zap日志
	if err := pkg.InitLogger(); err != nil {
		panic(err)
	}
	defer pkg.Sync()

	// 记录应用启动日志
	pkg.LogInfo("Application starting",
		zap.String("service.id", id),
		zap.String("service.name", Name),
		zap.String("service.version", Version),
	)

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	// 1. 创建配置实例，从文件加载配置

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf), // 从指定路径读取配置文件
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		pkg.LogError("Failed to load config", zap.Error(err))
		panic(err)
	}

	// 2. 将加载的配置反序列化到Bootstrap结构体
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		pkg.LogError("Failed to scan config", zap.Error(err))
		panic(err)
	}

	// 3. 通过wireApp（依赖注入初始化）构造整个应用，传入服务器配置、数据配置和日志
	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		pkg.LogError("Failed to initialize app", zap.Error(err))
		panic(err)
	}
	defer cleanup()

	// 4. 启动应用（启动grpc和http服务，注册服务到etcd）
	// start and wait for stop signal
	pkg.LogInfo("Application server started successfully")
	if err := app.Run(); err != nil {
		pkg.LogError("Application run failed", zap.Error(err))
		panic(err)
	}
	pkg.LogInfo("Application stopped gracefully")
}

// 创建 etcd 客户端，传入配置中的etcd节点地址和超时时间
func NewEtcdClient(c *conf.Server) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   c.Etcd.Endpoints,                // etcd集群地址列表
		DialTimeout: c.Etcd.DialTimeout.AsDuration(), // 连接超时时间
	})
}

// 基于etcd客户端创建kratos的注册器（注册服务实例到etcd）
func NewRegistrar(cli *clientv3.Client) registry.Registrar {
	return etcd.New(cli)
}

func NewContext() context.Context {
	return context.Background()
}
