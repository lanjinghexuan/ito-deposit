package service

import (
	"context"
	"ito-deposit/internal/data"
	"github.com/go-kratos/kratos/v2/log"
)

// LockerPointService 寄存点服务
type LockerPointService struct {
	data *data.Data
	log  *log.Helper
}

// NewLockerPointService 创建寄存点服务
func NewLockerPointService(data *data.Data, logger log.Logger) *LockerPointService {
	return &LockerPointService{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// GetAllLockerPoints 获取所有寄存点
func (s *LockerPointService) GetAllLockerPoints(ctx context.Context) ([]*data.LockerPoint, error) {
	var lockerPoints []*data.LockerPoint
	
	// 查询所有寄存点
	err := s.data.DB.Find(&lockerPoints).Error
	if err != nil {
		s.log.Errorf("查询寄存点失败: %v", err)
		return nil, err
	}
	
	s.log.Infof("查询到 %d 个寄存点", len(lockerPoints))
	return lockerPoints, nil
}