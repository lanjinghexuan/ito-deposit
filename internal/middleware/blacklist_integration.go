package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// EnableBlacklistMiddleware 为现有的中间件链添加拉黑检查功能
// 这个函数可以在需要时调用，不会破坏现有的架构
func EnableBlacklistMiddleware(existingMiddlewares []middleware.Middleware, db *gorm.DB, redis *redis.Client) []middleware.Middleware {
	// 在现有中间件的基础上添加拉黑检查
	blacklistMw := BlacklistMiddleware(db, redis)

	// 将拉黑中间件添加到现有中间件链的末尾
	// 这样它会在JWT认证之后执行
	return append(existingMiddlewares, blacklistMw)
}

// WrapWithBlacklist 包装现有的selector，添加拉黑检查
func WrapWithBlacklist(existingSelector middleware.Middleware, db *gorm.DB, redis *redis.Client) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		// 先执行现有的selector（包括JWT认证）
		wrappedHandler := existingSelector(handler)

		// 然后执行拉黑检查
		blacklistHandler := BlacklistMiddleware(db, redis)(wrappedHandler)

		return blacklistHandler
	}
}

// CreateBlacklistSelector 创建一个包含拉黑检查的完整selector
func CreateBlacklistSelector(jwtMiddleware middleware.Middleware, matchFunc selector.MatchFunc, db *gorm.DB, redis *redis.Client) middleware.Middleware {
	return selector.Server(
		jwtMiddleware,
		BlacklistMiddleware(db, redis),
	).Match(matchFunc).Build()
}
