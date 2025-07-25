package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	pb "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/data"
)

func TestDepositService_CreateDeposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 1. 创建mock对象
	mockDB := data.NewMockDBInterface(ctrl)
	mockRedis := data.NewMockRedisInterface(ctrl)
	mockServerConf := &conf.Server{Jwt: &conf.Server_Jwt{Authkey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM3NTUyNDMsImlkIjoiMTIzIiwieW91cl9jdXN0b21fY2xhaW0iOiJ5b3VyX2N1c3RvbV92YWx1ZSJ9.qcdoe8dSYtfQBZgCP30Yln4r8z9ovPDEF1fNVlviWX4"}}

	service := NewDepositService(&data.Data{DB: mockDB, Redis: mockRedis}, mockServerConf, nil)

	// 2. 设置测试用例
	cases := []struct {
		name         string
		req          *pb.CreateDepositRequest
		mock         func()
		expectedCode int32
		expectedErr  bool
	}{{
		name: "正常创建订单",
		req: &pb.CreateDepositRequest{
			CabinetId:         1,
			LockerType:        1,
			ScheduledDuration: 2,
		},
		mock: func() {
			// 设置价格规则查询mock
			mockDB.EXPECT().Table("locker_pricing_rules").Return(mockDB)
			mockDB.EXPECT().Where("network_id = ? ", 1).Return(mockDB)
			mockDB.EXPECT().Where("status = 1").Return(mockDB)
			mockDB.EXPECT().Where("locker_type = ?", 1).Return(mockDB)
			mockDB.EXPECT().Limit(1).Return(mockDB)
			mockDB.EXPECT().Find(gomock.Any()).SetArg(0, data.LockerPricingRules{
				IsDepositEnabled: 1,
				DepositAmount:    100,
				HourlyRate:       5,
			}).Return(nil)

			// 设置柜子查询mock
			mockDB.EXPECT().Table("lockers").Return(mockDB)
			mockDB.EXPECT().Where("locker_point_id = ?", 1).Return(mockDB)
			mockDB.EXPECT().Select("id").Return(mockDB)
			mockDB.EXPECT().Where("status = 1").Return(mockDB)
			mockDB.EXPECT().Limit(1).Return(mockDB)
			mockDB.EXPECT().Find(gomock.Any()).SetArg(0, data.Locker{Id: 1001}).Return(nil)

			// 设置事务mock
			mockDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f func(tx *gorm.DB) error) error {
				return f(mockDB)
			})

			// 设置更新柜子状态mock
			mockDB.EXPECT().Table("lockers").Return(mockDB)
			mockDB.EXPECT().Where("id = ?", 1001).Return(mockDB)
			mockDB.EXPECT().Update("status", 2).Return(&gorm.DB{RowsAffected: 1})

			// 设置创建订单mock
			mockDB.EXPECT().Table("locker_orders").Return(mockDB)
			mockDB.EXPECT().Create(gomock.Any()).Return(nil)
		},
		expectedCode: 200,
		expectedErr:  false,
	}, {
		name: "柜子数量不足",
		req: &pb.CreateDepositRequest{
			CabinetId:         1,
			LockerType:        1,
			ScheduledDuration: 2,
		},
		mock: func() {
			// 设置价格规则查询mock（同上）
			mockDB.EXPECT().Table("locker_pricing_rules").Return(mockDB)
			mockDB.EXPECT().Where("network_id = ? ", 1).Return(mockDB)
			mockDB.EXPECT().Where("status = 1").Return(mockDB)
			mockDB.EXPECT().Where("locker_type = ?", 1).Return(mockDB)
			mockDB.EXPECT().Limit(1).Return(mockDB)
			mockDB.EXPECT().Find(gomock.Any()).SetArg(0, data.LockerPricingRules{
				IsDepositEnabled: 1,
				DepositAmount:    100,
				HourlyRate:       5,
			}).Return(nil)

			// 设置柜子查询mock
			mockDB.EXPECT().Table("lockers").Return(mockDB)
			mockDB.EXPECT().Where("locker_point_id = ?", 1).Return(mockDB)
			mockDB.EXPECT().Select("id").Return(mockDB)
			mockDB.EXPECT().Where("status = 1").Return(mockDB)
			mockDB.EXPECT().Limit(1).Return(mockDB)
			mockDB.EXPECT().Find(gomock.Any()).SetArg(0, data.Locker{Id: 1001}).Return(nil)

			// 设置事务mock
			mockDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f func(tx *gorm.DB) error) error {
				return f(mockDB)
			})

			// 设置更新柜子状态mock
			mockDB.EXPECT().Table("lockers").Return(mockDB)
			mockDB.EXPECT().Where("id = ?", 1001).Return(mockDB)
			mockDB.EXPECT().Update("status", 2).Return(&gorm.DB{RowsAffected: 1})

			// 设置创建订单mock
			mockDB.EXPECT().Table("locker_orders").Return(mockDB)
			mockDB.EXPECT().Create(gomock.Any()).Return(nil)
		},
		expectedCode: 0,
		expectedErr:  true,
	}}

	// 3. 执行测试
	for _, c := range cases {
		c.mock()
		resp, err := service.CreateDeposit(context.Background(), c.req)
		assert.Equal(t, c.expectedErr, err != nil)
		if !c.expectedErr {
			assert.Equal(t, c.expectedCode, resp.Code)
			assert.NotEmpty(t, resp.Data.OrderNo)
		}
	}
}
