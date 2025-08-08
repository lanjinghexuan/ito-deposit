package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"ito-deposit/internal/data"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// BlacklistMiddleware 用户拉黑检查中间件
func BlacklistMiddleware(db *gorm.DB, redis *redis.Client) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 从上下文中获取JWT claims（由认证中间件设置）
			claims, ok := ctx.Value("claims").(jwt.MapClaims)
			if !ok {
				// 如果没有claims，说明可能是白名单接口，直接通过
				return handler(ctx, req)
			}

			// 从claims中提取用户ID
			userID, err := extractUserID(claims)
			if err != nil {
				log.Errorf("提取用户ID失败: %v", err)
				return handler(ctx, req) // 提取失败时继续执行，避免影响现有功能
			}

			// 检查用户是否被拉黑
			blacklistInfo, err := checkUserBlacklist(ctx, db, redis, userID)
			if err != nil {
				log.Errorf("检查用户拉黑状态失败: %v", err)
				return handler(ctx, req) // 检查失败时继续执行，避免影响现有功能
			}

			// 如果用户被拉黑，返回错误
			if blacklistInfo != nil {
				return nil, handleBlacklistedUser(blacklistInfo)
			}

			// 用户未被拉黑，继续执行
			return handler(ctx, req)
		}
	}
}

// extractUserID 从JWT claims中提取用户ID
func extractUserID(claims jwt.MapClaims) (int64, error) {
	// 尝试从不同的字段中提取用户ID
	if id, ok := claims["id"]; ok {
		switch v := id.(type) {
		case string:
			return strconv.ParseInt(v, 10, 64)
		case float64:
			return int64(v), nil
		case int64:
			return v, nil
		case int:
			return int64(v), nil
		}
	}

	// 尝试从user_id字段提取
	if userID, ok := claims["user_id"]; ok {
		switch v := userID.(type) {
		case string:
			return strconv.ParseInt(v, 10, 64)
		case float64:
			return int64(v), nil
		case int64:
			return v, nil
		case int:
			return int64(v), nil
		}
	}

	return 0, fmt.Errorf("无法从JWT claims中提取用户ID")
}

// checkUserBlacklist 检查用户是否被拉黑
func checkUserBlacklist(ctx context.Context, db *gorm.DB, redis *redis.Client, userID int64) (*data.UserBlacklist, error) {
	// 先从Redis缓存中查询
	cacheKey := fmt.Sprintf("user_blacklist:%d", userID)

	// 尝试从缓存获取
	if redis != nil {
		cached := redis.Get(ctx, cacheKey)
		if cached.Err() == nil {
			// 如果缓存中存在且值为"0"，表示用户未被拉黑
			if cached.Val() == "0" {
				return nil, nil
			}
			// 如果缓存中存在其他值，需要从数据库重新查询详细信息
		}
	}

	// 从数据库查询用户拉黑信息
	var blacklist data.UserBlacklist
	err := db.Where("user_id = ? AND is_active = 1", userID).
		Order("created_at DESC").
		First(&blacklist).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户未被拉黑，缓存结果
			if redis != nil {
				redis.Set(ctx, cacheKey, "0", 5*time.Minute)
			}
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户拉黑状态失败: %w", err)
	}

	// 检查拉黑是否已过期
	now := time.Now()
	if !blacklist.EndTime.IsZero() && now.After(blacklist.EndTime) {
		// 拉黑已过期，缓存结果
		if redis != nil {
			redis.Set(ctx, cacheKey, "0", 5*time.Minute)
		}
		return nil, nil
	}

	// 检查拉黑是否已生效
	if now.Before(blacklist.StartTime) {
		// 拉黑尚未生效，缓存结果
		if redis != nil {
			redis.Set(ctx, cacheKey, "0", 1*time.Minute) // 较短的缓存时间
		}
		return nil, nil
	}

	// 用户被拉黑，缓存结果
	if redis != nil {
		cacheTime := 5 * time.Minute
		if !blacklist.EndTime.IsZero() {
			// 如果有结束时间，缓存到结束时间
			remaining := blacklist.EndTime.Sub(now)
			if remaining > 0 && remaining < cacheTime {
				cacheTime = remaining
			}
		}
		redis.Set(ctx, cacheKey, "1", cacheTime)
	}

	return &blacklist, nil
}

// handleBlacklistedUser 处理被拉黑的用户
func handleBlacklistedUser(blacklist *data.UserBlacklist) error {
	// 根据封禁级别返回不同的错误信息
	switch blacklist.BanLevel {
	case 1: // 部分限制
		return errors.Forbidden("USER_RESTRICTED", fmt.Sprintf("用户访问受限: %s", blacklist.Reason))
	case 2: // 完全封禁
		return errors.Forbidden("USER_BANNED", fmt.Sprintf("用户已被封禁: %s", blacklist.Reason))
	default:
		return errors.Forbidden("USER_BLACKLISTED", fmt.Sprintf("用户已被拉黑: %s", blacklist.Reason))
	}
}

// GetBlacklistInfo 获取用户拉黑信息（供其他服务使用）
func GetBlacklistInfo(ctx context.Context, db *gorm.DB, redis *redis.Client, userID int64) (*data.UserBlacklist, error) {
	return checkUserBlacklist(ctx, db, redis, userID)
}

// ClearBlacklistCache 清除用户拉黑缓存（当拉黑状态变更时调用）
func ClearBlacklistCache(ctx context.Context, redis *redis.Client, userID int64) error {
	if redis == nil {
		return nil
	}

	cacheKey := fmt.Sprintf("user_blacklist:%d", userID)
	return redis.Del(ctx, cacheKey).Err()
}
