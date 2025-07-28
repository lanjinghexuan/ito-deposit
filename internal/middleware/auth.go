package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
	"ito-deposit/internal/conf"
	"strings"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(conf *conf.Server_Jwt) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				// 从请求头中获取令牌
				tokenString := tr.RequestHeader().Get("Authorization")
				if tokenString == "" {
					return nil, errors.Unauthorized("UNAUTHORIZED", "JWT token is missing")
				}

				// 移除Bearer前缀
				tokenString = strings.TrimPrefix(tokenString, "Bearer ")

				// 解析JWT令牌
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					// 验证签名算法
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.Unauthorized("UNAUTHORIZED", "Invalid token signing method")
					}
					return []byte(conf.Authkey), nil
				})

				if err != nil {
					return nil, errors.Unauthorized("UNAUTHORIZED", "Invalid or expired token")
				}

				// 验证令牌有效性
				if !token.Valid {
					return nil, errors.Unauthorized("UNAUTHORIZED", "Invalid token")
				}

				// 获取令牌中的Claims
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					return nil, errors.Unauthorized("UNAUTHORIZED", "Invalid token claims")
				}

				// 将Claims添加到上下文中
				ctx = context.WithValue(ctx, "claims", claims)
			}
			return handler(ctx, req)
		}
	}
}