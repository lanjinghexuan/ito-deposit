package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	pb "ito-deposit/api/helloworld/v1"
	"strconv"
)

type AdminRepo interface {
	SetPriceRule(context.Context, int32, []*LockerPricingRules) error
	AddPointAddPoint(ctx context.Context, point *LockerPoint) (*pb.AddPointRes, error)
	UpdatePoint(ctx context.Context, point *LockerPoint, intUserId int32) (*pb.UpdatePointRes, error)
	FindPoint(ctx context.Context, id int32, userId int32) (*LockerPoint, error)
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
		PointType:       req.PointType,
		PointImage:      req.PointImage,
		Status:          "1",
	}
	res, err := u.repo.AddPointAddPoint(ctx, points)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (u *AdminUsecase) UpdatePoint(ctx context.Context,
	req *pb.UpdatePointReq, userId string) (*pb.UpdatePointRes, error) {

	// 参数合法性
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid userId: %w", err)
	}

	// 查询原始记录
	original, err := u.repo.FindPoint(ctx, req.Point.Id, int32(uid))
	if err != nil {
		return nil, fmt.Errorf("failed to find point: %w", err)
	}

	// 创建更新结构体（以原始数据为基础）
	updated := *original // 值拷贝

	// 标记是否有变化
	changed := false

	// 以下手动比对并仅更新变化字段
	if req.Point.LocationId != updated.LocationId {
		updated.LocationId = req.Point.LocationId
		changed = true
	}
	if req.Point.Name != updated.Name {
		updated.Name = req.Point.Name
		changed = true
	}
	if req.Point.Address != updated.Address {
		updated.Address = req.Point.Address
		changed = true
	}
	if float64(req.Point.Latitude) != updated.Latitude {
		updated.Latitude = float64(req.Point.Latitude)
		changed = true
	}
	if float64(req.Point.Longitude) != updated.Longitude {
		updated.Longitude = float64(req.Point.Longitude)
		changed = true
	}
	if req.Point.AvailableLarge != updated.AvailableLarge {
		updated.AvailableLarge = req.Point.AvailableLarge
		changed = true
	}
	if req.Point.AvailableMedium != updated.AvailableMedium {
		updated.AvailableMedium = req.Point.AvailableMedium
		changed = true
	}
	if req.Point.AvailableSmall != updated.AvailableSmall {
		updated.AvailableSmall = req.Point.AvailableSmall
		changed = true
	}
	if req.Point.OpenTime != updated.OpenTime {
		updated.OpenTime = req.Point.OpenTime
		changed = true
	}
	if req.Point.Mobile != updated.Mobile {
		updated.Mobile = req.Point.Mobile
		changed = true
	}
	if req.Point.PointType != updated.PointType {
		updated.PointType = req.Point.PointType
		changed = true
	}
	if req.Point.PointImage != updated.PointImage {
		updated.PointImage = req.Point.PointImage
		changed = true
	}
	if req.Point.Status != updated.Status {
		updated.Status = req.Point.Status
		changed = true
	}

	// 无变化直接返回
	if !changed {
		return &pb.UpdatePointRes{
			Code: 200,
			Msg:  "无数据更新",
		}, nil
	}

	// 执行更新
	res, err := u.repo.UpdatePoint(ctx, &updated, int32(uid))
	if err != nil {
		return nil, fmt.Errorf("failed to update point: %w", err)
	}

	return res, nil
}
