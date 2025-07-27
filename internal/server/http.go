package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	minio1 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	v1 "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/service"
	"mime/multipart"
	http2 "net/http"
	"time"
)

func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, order *service.OrderService, user *service.UserService,
	home *service.HomeService, deposit *service.DepositService, admin *service.AdminService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Filter(corsFilter),
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	opts = append(opts, http.Middleware(
		selector.Server(
			jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
				return []byte(c.Jwt.Authkey), nil
			}, jwt.WithSigningMethod(jwtv5.SigningMethodHS256), jwt.WithClaims(func() jwtv5.Claims {
				return &jwtv5.MapClaims{}
			})),
		).
			Match(NewWhiteListMatcher()).
			Build(),
	))

	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterUserHTTPServer(srv, user)
	v1.RegisterHomeHTTPServer(srv, home)
	v1.RegisterDepositHTTPServer(srv, deposit)
	v1.RegisterOrderHTTPServer(srv, order)
	v1.RegisterAdminHTTPServer(srv, admin)

	srv.Route("/").POST("/upload", admin.DownloadFile)

	return srv
}

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	// 添加不需要 JWT 验证的接口到白名单
	whiteList["/api.helloworld.v1.Deposit/ReturnToken"] = struct{}{}
	whiteList["/api.helloworld.v1.User/SendSms"] = struct{}{}
	whiteList["/api.helloworld.v1.User/Register"] = struct{}{}
	whiteList["/api.helloworld.v1.User/Login"] = struct{}{}
	whiteList["/api.helloworld.v1.User/List"] = struct{}{}
	whiteList["/api.helloworld.v1.User/Admin"] = struct{}{}
	whiteList["/api.helloworld.v1.User/GetUser"] = struct{}{}
	whiteList["/api.helloworld.v1.Order/ListOrder"] = struct{}{}
	whiteList["/api.helloworld.v1.Order/ShowOrder"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// corsFilter 手动实现 CORS
func corsFilter(next http2.Handler) http2.Handler {
	return http2.HandlerFunc(func(w http2.ResponseWriter, r *http2.Request) {
		// 允许的域名、方法、头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		// 预检请求直接返回 204
		if r.Method == http2.MethodOptions {
			w.WriteHeader(http2.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func downloadFile(ctx kratosHttp.Context) error {
	req := ctx.Request() // 拿到 *http.Request
	file, header, err := req.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	//url, err := UploadFile("aaa", file, *conf.Data)
	var url string
	// 接下来就能用 file、header.Filename、header.Size
	fmt.Println("upload:", header.Filename, header.Size)

	return ctx.Result(200, map[string]string{"url": url})
}

func UploadFile(objectName string, file *multipart.FileHeader, c *conf.Data) (string, error) {
	//方法为实现参数自行编写
	var addr string
	addr = c.Minio.Endpoint
	// 初始化 MinIO 客户端
	minioClient, err := minio1.New(addr, &minio1.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKeyId, c.Minio.AccessKeySecret, ""),
		Secure: c.Minio.UseSsl,
	})
	if err != nil {
		return "", fmt.Errorf("failed to initialize MinIO client: %v", err)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// 获取文件信息
	fileStat, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to get file stats: %v", err)
	}
	defer fileStat.Close()
	objectName = fmt.Sprintf("%s/%s", time.Now().Format("2006-01-02"), objectName)
	// 使用 PutObject 上传文件
	_, err = minioClient.PutObject(
		context.Background(),
		c.Minio.BucketName,
		objectName,
		src,
		file.Size,
		minio1.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %v", err)
	}

	fmt.Printf("Successfully uploaded file %s to bucket %s as object %s\n",
		file.Filename, c.Minio.BucketName, objectName)
	return addr + "/" + c.Minio.BucketName + "/" + objectName, nil
}
