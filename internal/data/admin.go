package data

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ito-deposit/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type AdminRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewAdminRepo(data *Data, logger log.Logger) biz.AdminRepo {
	return &AdminRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (s *AdminRepo) SetPriceRule(ctx context.Context, networkId int32, data []*biz.LockerPricingRules) error {
	// 1. 开启事务
	tx := s.data.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 2. 停用旧规则（软删除）
	if err := tx.Model(&biz.LockerPricingRules{}).
		Where("network_id = ? AND status = 1", networkId).
		Update("status", 0).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.Internal, "停用旧规则失败: %v", err)
	}
	err := tx.Model(&biz.LockerPricingRules{}).Where("network_id = ? AND status = 1", networkId).
		Update("status", 0).Error
	if err != nil {
		tx.Rollback()
		return status.Errorf(codes.Internal, "停用旧规则失败: %v", err)
	}

	// 3. 创建新规则
	err = tx.Create(data).Error
	if err != nil {
		tx.Rollback()
	}

	// 4. 提交事务
	if err := tx.Commit().Error; err != nil {
		return status.Errorf(codes.Internal, "提交事务失败: %v", err)
	}

	return nil
}
