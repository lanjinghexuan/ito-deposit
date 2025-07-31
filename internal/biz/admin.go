package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	pb "ito-deposit/api/helloworld/v1"
	"strconv"
)

type AdminRepo interface {
	SetPriceRule(context.Context, int32, []*LockerPricingRules) error
	AddPointAddPoint(ctx context.Context, point *LockerPoint) (*pb.AddPointRes, error)
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

func (u *AdminUsecase) AddPointAddPoint(ctx context.Context, req *pb.AddPointReq, userId string) (*pb.AddPointRes, error) {
	var points *LockerPoint
	intUserId, _ := strconv.ParseInt(userId, 10, 64)
	points = &LockerPoint{
		LocationId:      req.LocationId,
		Name:            req.Name,
		Address:         req.Address,
		Latitude:        float64(req.Latitude),
		Longitude:       float64(req.Longitude),
		AvailableLarge:  req.AvailableLarge,
		AvailableMedium: req.AvailableMedium,
		AvailableSmall:  req.AvailableSmall,
		OpenTime:        req.OpenTime,
		Mobile:          req.Mobile,
		AdminId:         int32(intUserId),
	}
	res, err := u.repo.AddPointAddPoint(ctx, points)
	if err != nil {
		return nil, err
	}
	return res, nil

}
