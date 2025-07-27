package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
)

type AdminRepo interface {
	SetPriceRule(context.Context, int32, []*LockerPricingRules) error
}

type AdminUsecase struct {
	repo AdminRepo
	log  *log.Helper
}

func NewAdminUsecase(repo AdminRepo, logger log.Logger) *AdminUsecase {
	return &AdminUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (u *AdminUsecase) SetPriceRule(ctx context.Context, networkId int32, data []*LockerPricingRules) error {
	if networkId == 0 {
		return errors.New("网点id为空")
	}

	return u.repo.SetPriceRule(ctx, networkId, data)
}
