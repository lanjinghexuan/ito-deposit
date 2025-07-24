package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	pb "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/conf" // 添加缺失的配置包导入
	"ito-deposit/internal/data"

	jwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// 用于测试的jwt context key，变量名与kratos jwt包内部一致
var contextKey = struct{}{}

// 覆盖 Kratos jwt1.FromContext，保证测试用例和业务一致
func FromContext(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(contextKey)
	if v == nil {
		return nil, false
	}
	return v, true
}

func TestDepositService_CreateDeposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 1. 创建mock对象
	mockDB := data.NewMockDBInterface(ctrl)
	mockRedis := data.NewMockRedisInterface(ctrl)

	// 创建模拟的服务器配置
	serverConf := &conf.Server{}

	// 使用数据层提供的导出方法获取接口实例
	serviceData := &data.Data{
		DBI:    data.GetDBInterface(mockDB),
		RedisI: data.GetRedisInterface(mockRedis),
	}

	// 初始化服务
	service := NewDepositService(serviceData, serverConf, nil)

	// 构造带token的context
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM3NTUyNDMsImlkIjoiMTIzIiwieW91cl9jdXN0b21fY2xhaW0iOiJ5b3VyX2N1c3RvbV92YWx1ZSJ9.qcdoe8dSYtfQBZgCP30Yln4r8z9ovPDEF1fNVlviWX4"
	claims := jwtv5.MapClaims{}
	_, _ = jwtv5.ParseWithClaims(tokenStr, &claims, func(token *jwtv5.Token) (interface{}, error) {
		return []byte("a9999"), nil // 密钥与配置一致
	})
	ctx := jwt.NewContext(context.Background(), &claims)

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
		// 使用闭包捕获当前测试用例的req变量
		mock: func(req *pb.CreateDepositRequest) func() {
		    return func() {
		        // 设置价格规则查询mock
		        mockDB.EXPECT().WithContext(ctx).Return(mockDB)
		        mockDB.EXPECT().Table("locker_pricing_rules").Return(mockDB)
		        // 现在可以正确访问req变量
		        mockDB.EXPECT().Where("network_id = ?", req.CabinetId).Return(mockDB)
		        mockDB.EXPECT().Where("status = ?", 1).Return(mockDB)
		        mockDB.EXPECT().Where("locker_type = ?", req.LockerType).Return(mockDB)
		        mockDB.EXPECT().Limit(1).Return(mockDB)
		        mockDB.EXPECT().Find(gomock.Any()).SetArg(0, data.LockerPricingRules{
		            IsDepositEnabled: 1,
		            DepositAmount:    100,
		            HourlyRate:       5,
		        }).Return(nil)

			// 设置柜子查询mock
			mockDB.EXPECT().Table("lockers").Return(mockDB)
			// 修正参数名称以匹配实际代码
			mockDB.EXPECT().Where("locker_point_id = ?", req.CabinetId).Return(mockDB)
			mockDB.EXPECT().Select("id").Return(mockDB)
			mockDB.EXPECT().Where("status = ?", 1).Return(mockDB)
			mockDB.EXPECT().Limit(1).Return(mockDB)
			mockDB.EXPECT().Find(gomock.Any()).SetArg(0, data.Locker{Id: 1001}).Return(nil)

			// 设置事务mock
			mockDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f func(tx data.DBInterface) error) error {
				return f(mockDB)
			})

			// 设置更新柜子状态mock
			mockDB.EXPECT().Table("lockers").Return(mockDB)
			mockDB.EXPECT().Where("id = ?", 1001).Return(mockDB)
			mockDB.EXPECT().Update("status", 2).Return(&gorm.DB{RowsAffected: 1})

			// 设置创建订单mock
			mockDB.EXPECT().Table("locker_orders").Return(mockDB)
			mockDB.EXPECT().Create(gomock.Any()).Return(nil)
		}, // 添加缺少的逗号
		// 修复结构体字段间缺少的逗号
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
			// 设置价格规则查询mock
			mockDB.EXPECT().WithContext(ctx).Return(mockDB)
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
			mockDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(f func(tx data.DBInterface) error) error {
				return f(mockDB)
			})

			// 设置更新柜子状态mock
			mockDB.EXPECT().Table("lockers").Return(mockDB)
			mockDB.EXPECT().Where("id = ?", 1001).Return(mockDB)
			mockDB.EXPECT().Update("status", 2).Return(&gorm.DB{RowsAffected: 1})

			// 设置创建订单mock
			mockDB.EXPECT().Table("locker_orders").Return(mockDB)
			mockDB.EXPECT().Create(gomock.Any()).Return(nil)
		}, // 添加缺少的逗号
		// 修复结构体字段间缺少的逗号
		expectedCode: 0,
		expectedErr:  true,
	}}

	// 3. 执行测试
	for _, c := range cases {
		c.mock()
		resp, err := service.CreateDeposit(ctx, c.req)
		assert.Equal(t, c.expectedErr, err != nil)
		if !c.expectedErr {
			assert.Equal(t, c.expectedCode, resp.Code)
			assert.NotEmpty(t, resp.Data.OrderNo)
		}
	}
}
