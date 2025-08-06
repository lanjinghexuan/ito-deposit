package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "./configs/gateway.yaml", "config file path")
}

type GatewayConfig struct {
	Server struct {
		Http struct {
			Addr string
		}
		Etcd struct {
			Endpoints []string `yaml:"endpoints"`
		} `yaml:"etcd"`
	} `yaml:"server"`
}

func main() {
	flag.Parse()
	fmt.Println("配置文件", flagconf)

	// 1. 读取配置
	c := config.New(config.WithSource(file.NewSource(flagconf)))
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc GatewayConfig
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	fmt.Println(bc)

	// 2. 初始化etcd客户端和注册中心
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: bc.Server.Etcd.Endpoints,
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(cli) // 启用 etcd 监听机制 + 本地缓存)

	// 3. 反向代理 Director 函数
	director := func(req *http.Request) {
		ctx := context.Background()
		services, err := r.GetService(ctx, "ito-deposit")
		if err != nil || len(services) == 0 {
			log.Println("no available instance:", err)
			req.URL.Host = ""
			return
		}

		// 合并所有服务节点
		var allEndpoints []string
		for _, svc := range services {
			for _, ep := range svc.Endpoints {
				if strings.HasPrefix(ep, "http://") {
					allEndpoints = append(allEndpoints, strings.TrimPrefix(ep, "http://"))
				}
			}
		}

		if len(allEndpoints) == 0 {
			log.Println("no http endpoints found for service")
			req.URL.Host = ""
			return
		}

		// 随机选一个节点
		rand.Seed(time.Now().UnixNano())
		target := allEndpoints[rand.Intn(len(allEndpoints))]
		fmt.Println(target)
		req.URL.Scheme = "http"
		req.URL.Host = target
	}

	proxy := &httputil.ReverseProxy{Director: director}

	// 4. 启动 HTTP 服务
	log.Println("Gateway listening on", bc.Server.Http.Addr)
	if err := http.ListenAndServe(bc.Server.Http.Addr, proxy); err != nil {
		log.Fatal(err)
	}
}
