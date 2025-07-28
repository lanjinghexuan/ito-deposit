package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"ito-deposit/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	// 智能检测配置文件路径
	defaultConf := findConfigPath()
	flag.StringVar(&flagconf, "conf", defaultConf, "config path, eg: -conf config.yaml")
}

// findConfigPath 智能查找配置文件路径
func findConfigPath() string {
	// 可能的配置路径
	paths := []string{
		"./configs",           // 从项目根目录运行
		"../../configs",       // 从cmd/ito-deposit目录运行
		"../../../configs",    // 其他可能的路径
	}
	
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
		// 也检查config.yaml文件
		configFile := filepath.Join(path, "config.yaml")
		if _, err := os.Stat(configFile); err == nil {
			return path
		}
	}
	
	// 如果都找不到，返回默认路径
	return "./configs"
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()
	wd, _ := os.Getwd()
	fmt.Println("cwd:", wd)
	fmt.Println("config path:", flagconf)
	
	// 检查配置路径是否存在
	if _, err := os.Stat(flagconf); os.IsNotExist(err) {
		fmt.Printf("Config path does not exist: %s\n", flagconf)
		// 尝试查找配置文件
		if _, err := os.Stat(filepath.Join(flagconf, "config.yaml")); os.IsNotExist(err) {
			fmt.Printf("Config file does not exist: %s/config.yaml\n", flagconf)
		}
	}
	
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
