package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	pb "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/data"

	jwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// createTestContext 创建一个包含JWT token的测试上下文
// 参数：userID - 用户ID，会被放入JWT token的claims中
// 返回：包含JWT token的context，用于模拟已登录用户的请求
func createTestContext(userID string) context.Context {
	// 创建JWT claims（声明），包含用户ID和过期时间
	claims := jwtv5.MapClaims{
		"id":  userID,                           // 用户ID
		"exp": time.Now().Add(time.Hour).Unix(), // 过期时间：1小时后
	}
	// 将JWT token包装到context中，模拟Kratos框架的JWT中间件行为
	return jwt.NewContext(context.Background(), &claims)
}

// createEmptyContext 创建一个空的context（不包含JWT token）
// 用于测试未登录用户的请求场景
func createEmptyContext() context.Context {
	return context.Background()
}

// TestDepositService_CreateDeposit 测试创建寄存订单的功能
// 这是最核心的测试函数，测试了创建订单的各种场景
func TestDepositService_CreateDeposit(t *testing.T) {
	// 创建gomock控制器，用于管理所有的mock对象
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 测试结束时清理mock对象

	// 创建mock数据库接口和Redis接口
	// 这些mock对象会模拟真实的数据库和Redis操作
	mockDB := data.NewMockDBInterface(ctrl)
	mockRedis := data.NewMockRedisInterface(ctrl)

	// 创建服务器配置（这里使用空配置）
	serverConf := &conf.Server{}

	// 创建服务数据结构，注入mock对象
	serviceData := &data.Data{
		DBI:    data.GetDBInterface(mockDB),       // 注入mock数据库接口
		RedisI: data.GetRedisInterface(mockRedis), // 注入mock Redis接口
	}

	// 创建被测试的服务实例
	service := NewDepositService(serviceData, serverConf, nil)

	// 测试场景1：无效token的情况
	// 当用户未登录或token无效时，应该返回401错误
	t.Run("无效token", func(t *testing.T) {
		// 构造请求参数
		req := &pb.CreateDepositRequest{
			CabinetId:         1, // 柜子网点ID
			LockerType:        1, // 柜子类型（1=小柜，2=大柜等）
			ScheduledDuration: 2, // 预计寄存时长（小时）
		}

		// 使用空的context（不包含JWT token）调用服务
		resp, err := service.CreateDeposit(createEmptyContext(), req)

		// 验证结果：应该没有系统错误，但业务返回401
		assert.NoError(t, err)                           // 没有系统级错误
		assert.NotNil(t, resp)                           // 响应不为空
		assert.Equal(t, int32(401), resp.Code)           // 业务错误码为401
		assert.Equal(t, "token不正确或者未传", resp.Msg) // 错误消息正确
	})

	// 测试场景2：价格规则查询失败的情况
	// 模拟数据库查询价格规则时发生错误
	t.Run("价格规则查询失败", func(t *testing.T) {
		// 设置mock期望：模拟数据库查询链式调用
		// 这些调用对应业务代码中查询价格规则的SQL操作
		mockDB.EXPECT().Table("locker_pricing_rules").Return(mockDB)            // 选择表
		mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockDB)        // WHERE network_id = ?
		mockDB.EXPECT().Where("status = 1").Return(mockDB)                      // WHERE status = 1
		mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockDB)        // WHERE locker_type = ?
		mockDB.EXPECT().Limit(1).Return(mockDB)                                 // LIMIT 1
		mockDB.EXPECT().Find(gomock.Any()).Return(errors.New("数据库查询失败")) // 模拟查询失败

		// 构造请求参数
		req := &pb.CreateDepositRequest{
			CabinetId:         1,
			LockerType:        1,
			ScheduledDuration: 2,
		}

		// 调用服务方法，使用有效的token
		_, err := service.CreateDeposit(createTestContext("123"), req)

		// 验证结果：应该返回系统错误
		assert.Error(t, err)                              // 有系统级错误
		assert.Contains(t, err.Error(), "数据库查询失败") // 错误消息包含预期内容
	})

	// 测试场景3：成功创建订单的情况
	// 这是最复杂的测试场景，模拟完整的订单创建流程
	t.Run("成功创建订单", func(t *testing.T) {
		// 第一步：设置价格规则查询的mock期望
		// 模拟查询价格规则表，返回有效的价格配置
		mockDB.EXPECT().Table("locker_pricing_rules").Return(mockDB)     // 选择价格规则表
		mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockDB) // WHERE network_id = ?
		mockDB.EXPECT().Where("status = 1").Return(mockDB)               // WHERE status = 1 (有效状态)
		mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockDB) // WHERE locker_type = ?
		mockDB.EXPECT().Limit(1).Return(mockDB)                          // LIMIT 1
		mockDB.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest interface{}) error {
			// 模拟查询成功，返回价格规则数据
			rules := dest.(*data.LockerPricingRules)
			*rules = data.LockerPricingRules{
				IsDepositEnabled: 1,     // 启用押金
				DepositAmount:    100.0, // 押金100元
				HourlyRate:       5.0,   // 每小时5元
			}
			return nil
		})

		// 第二步：设置数据库事务的mock期望
		// 在事务中完成：查找可用柜子 -> 更新柜子状态 -> 创建订单
		mockDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fn func(tx data.DBInterface) error) error {
			// 创建事务内的mock数据库接口
			mockTx := data.NewMockDBInterface(ctrl)

			// 2.1 在事务内查询可用柜子
			mockTx.EXPECT().Table("lockers").Return(mockTx)                  // 选择柜子表
			mockTx.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockTx) // WHERE locker_point_id = ?
			mockTx.EXPECT().Select("id").Return(mockTx)                      // SELECT id
			mockTx.EXPECT().Where("status = 1").Return(mockTx)               // WHERE status = 1 (可用状态)
			mockTx.EXPECT().Limit(1).Return(mockTx)                          // LIMIT 1
			mockTx.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest interface{}) error {
				// 模拟找到可用柜子
				locker := dest.(*data.Locker)
				*locker = data.Locker{Id: 1001} // 返回柜子ID为1001
				return nil
			})

			// 2.2 更新柜子状态为已占用
			mockTx.EXPECT().Table("lockers").Return(mockTx)                                      // 选择柜子表
			mockTx.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockTx)                     // WHERE id = 1001
			mockTx.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&gorm.DB{RowsAffected: 1}) // UPDATE status = 2, 影响1行

			// 2.3 创建订单记录
			mockTx.EXPECT().Table("locker_orders").Return(mockTx) // 选择订单表
			mockTx.EXPECT().Create(gomock.Any()).Return(nil)      // INSERT 订单数据

			// 执行事务函数
			return fn(mockTx)
		})

		// 构造请求参数
		req := &pb.CreateDepositRequest{
			CabinetId:         1, // 网点ID
			LockerType:        1, // 柜子类型
			ScheduledDuration: 2, // 预计寄存2小时
		}

		// 调用服务方法
		resp, err := service.CreateDeposit(createTestContext("123"), req)

		// 验证结果：应该成功创建订单
		assert.NoError(t, err)                          // 没有系统错误
		assert.NotNil(t, resp)                          // 响应不为空
		assert.Equal(t, int32(200), resp.Code)          // 成功状态码
		assert.Equal(t, "添加寄存订单成功", resp.Msg)   // 成功消息
		assert.NotNil(t, resp.Data)                     // 返回数据不为空
		assert.NotEmpty(t, resp.Data.OrderNo)           // 订单号不为空
		assert.Greater(t, resp.Data.LockerId, int32(0)) // 柜子ID大于0
	})

	// 测试场景4：无可用柜子的情况
	// 模拟价格规则查询成功，但找不到可用柜子的场景
	t.Run("无可用柜子", func(t *testing.T) {
		// 第一步：设置价格规则查询成功
		mockDB.EXPECT().Table("locker_pricing_rules").Return(mockDB)     // 选择价格规则表
		mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockDB) // WHERE network_id = ?
		mockDB.EXPECT().Where("status = 1").Return(mockDB)               // WHERE status = 1
		mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockDB) // WHERE locker_type = ?
		mockDB.EXPECT().Limit(1).Return(mockDB)                          // LIMIT 1
		mockDB.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest interface{}) error {
			// 模拟查询价格规则成功
			rules := dest.(*data.LockerPricingRules)
			*rules = data.LockerPricingRules{
				IsDepositEnabled: 1,     // 启用押金
				DepositAmount:    100.0, // 押金100元
				HourlyRate:       5.0,   // 每小时5元
			}
			return nil
		})

		// 第二步：设置事务处理 - 模拟无可用柜子的情况
		mockDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fn func(tx data.DBInterface) error) error {
			mockTx := data.NewMockDBInterface(ctrl)

			// 在事务内查询柜子，但找不到可用的
			mockTx.EXPECT().Table("lockers").Return(mockTx)                  // 选择柜子表
			mockTx.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockTx) // WHERE locker_point_id = ?
			mockTx.EXPECT().Select("id").Return(mockTx)                      // SELECT id
			mockTx.EXPECT().Where("status = 1").Return(mockTx)               // WHERE status = 1 (可用状态)
			mockTx.EXPECT().Limit(1).Return(mockTx)                          // LIMIT 1
			mockTx.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest interface{}) error {
				// 返回空的柜子结构，表示没有找到可用柜子
				// 注意：这里不设置locker.Id，保持为0，表示无可用柜子
				return nil
			})

			// 执行事务函数
			return fn(mockTx)
		})

		// 构造请求参数
		req := &pb.CreateDepositRequest{
			CabinetId:         1,
			LockerType:        1,
			ScheduledDuration: 2,
		}

		// 调用服务方法
		_, err := service.CreateDeposit(createTestContext("123"), req)

		// 验证结果：应该返回"寄存柜可用数量不足"的错误
		assert.Error(t, err)                                  // 有系统级错误
		assert.Contains(t, err.Error(), "寄存柜可用数量不足") // 错误消息包含预期内容
	})
}

// TestDepositService_ReturnToken 测试JWT token生成功能
// 这个方法用于生成JWT token，通常用于用户登录后获取访问令牌
func TestDepositService_ReturnToken(t *testing.T) {
	// 创建服务器配置，包含JWT密钥
	serverConf := &conf.Server{
		Jwt: &conf.Server_Jwt{
			Authkey: "test_secret_key", // JWT签名密钥
		},
	}

	// 创建服务实例，注入配置
	service := NewDepositService(&data.Data{}, serverConf, nil)

	// 测试成功生成token的场景
	t.Run("成功生成token", func(t *testing.T) {
		// 调用生成token的方法
		resp, err := service.ReturnToken(context.Background(), &pb.ReturnTokenReq{})

		// 验证结果
		assert.NoError(t, err)                     // 没有系统错误
		assert.NotNil(t, resp)                     // 响应不为空
		assert.Equal(t, int32(200), resp.Coe)      // 成功状态码
		assert.Equal(t, "生成token成功", resp.Msg) // 成功消息
		assert.NotEmpty(t, resp.Token)             // token不为空

		// 验证JWT token格式（标准JWT由三部分组成：header.payload.signature）
		parts := len(strings.Split(resp.Token, "."))
		assert.Equal(t, 3, parts, "JWT应该有三个部分")
	})
}

// TestDepositService_DecodeToken 测试JWT token解析功能
// 这个方法用于解析JWT token，提取其中的用户信息
func TestDepositService_DecodeToken(t *testing.T) {
	// 创建服务实例（不需要特殊配置）
	service := NewDepositService(&data.Data{}, &conf.Server{}, nil)

	// 测试场景1：成功解析有效token
	t.Run("成功解析token", func(t *testing.T) {
		// 创建包含用户ID为"123"的测试context
		ctx := createTestContext("123")

		// 调用token解析方法
		resp, err := service.DecodeToken(ctx, &pb.ReturnTokenReq{})

		// 验证结果：应该成功解析出用户ID
		assert.NoError(t, err)                  // 没有系统错误
		assert.NotNil(t, resp)                  // 响应不为空
		assert.Equal(t, int32(200), resp.Coe)   // 成功状态码
		assert.Equal(t, "token内容 ", resp.Msg) // 成功消息
		assert.Equal(t, "123", resp.Token)      // 返回的用户ID正确
	})

	// 测试场景2：解析无效token（空context）
	t.Run("无效token", func(t *testing.T) {
		// 使用空的context（不包含JWT token）
		ctx := createEmptyContext()

		// 调用token解析方法
		resp, err := service.DecodeToken(ctx, &pb.ReturnTokenReq{})

		// 验证结果：应该返回401错误
		assert.NoError(t, err)                                                       // 没有系统级错误
		assert.NotNil(t, resp)                                                       // 响应不为空
		assert.Equal(t, int32(401), resp.Coe)                                        // 错误状态码401
		assert.Equal(t, "未找到有效的 JWT Token（可能未登录或 Token 无效）", resp.Msg) // 错误消息正确
	})
}

// TestDepositService_UnimplementedMethods 测试未实现的方法
// 这些方法在当前版本中只是返回空结构体，但我们需要确保它们不会崩溃
func TestDepositService_UnimplementedMethods(t *testing.T) {
	// 创建服务实例（使用最简配置）
	service := NewDepositService(&data.Data{}, &conf.Server{}, nil)

	// 测试更新寄存订单方法（目前未实现具体逻辑）
	t.Run("UpdateDeposit", func(t *testing.T) {
		resp, err := service.UpdateDeposit(context.Background(), &pb.UpdateDepositRequest{})
		assert.NoError(t, err) // 不应该有错误
		assert.NotNil(t, resp) // 应该返回空的响应结构体
	})

	// 测试删除寄存订单方法（目前未实现具体逻辑）
	t.Run("DeleteDeposit", func(t *testing.T) {
		resp, err := service.DeleteDeposit(context.Background(), &pb.DeleteDepositRequest{})
		assert.NoError(t, err) // 不应该有错误
		assert.NotNil(t, resp) // 应该返回空的响应结构体
	})

	// 测试获取单个寄存订单方法（目前未实现具体逻辑）
	t.Run("GetDeposit", func(t *testing.T) {
		resp, err := service.GetDeposit(context.Background(), &pb.GetDepositRequest{})
		assert.NoError(t, err) // 不应该有错误
		assert.NotNil(t, resp) // 应该返回空的响应结构体
	})

	// 测试获取寄存订单列表方法（目前未实现具体逻辑）
	t.Run("ListDeposit", func(t *testing.T) {
		resp, err := service.ListDeposit(context.Background(), &pb.ListDepositRequest{})
		assert.NoError(t, err) // 不应该有错误
		assert.NotNil(t, resp) // 应该返回空的响应结构体
	})
}

// TestNewDepositService 测试服务构造函数
// 验证服务实例能够正确创建，并且依赖注入正常工作
func TestNewDepositService(t *testing.T) {
	// 创建测试用的依赖对象
	data := &data.Data{}     // 数据层对象
	server := &conf.Server{} // 服务器配置对象

	// 调用构造函数创建服务实例
	service := NewDepositService(data, server, nil)

	// 验证服务实例创建成功
	assert.NotNil(t, service)               // 服务实例不为空
	assert.Equal(t, data, service.data)     // 数据层依赖注入正确
	assert.Equal(t, server, service.server) // 服务器配置注入正确
}

// TestDepositService_EdgeCases 测试边界情况和特殊场景
// 这些测试用于验证系统在特殊条件下的行为
func TestDepositService_EdgeCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 测试场景：无押金情况的价格计算
	// 验证当价格规则中禁用押金时，订单价格计算是否正确
	t.Run("CreateDeposit_无押金情况", func(t *testing.T) {
		// 创建mock对象和服务实例
		mockDB := data.NewMockDBInterface(ctrl)
		serviceData := &data.Data{DBI: data.GetDBInterface(mockDB)}
		service := NewDepositService(serviceData, &conf.Server{}, nil)

		// 第一步：设置价格规则查询 - 返回无押金的配置
		mockDB.EXPECT().Table("locker_pricing_rules").Return(mockDB)     // 选择价格规则表
		mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockDB) // WHERE network_id = ?
		mockDB.EXPECT().Where("status = 1").Return(mockDB)               // WHERE status = 1
		mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockDB) // WHERE locker_type = ?
		mockDB.EXPECT().Limit(1).Return(mockDB)                          // LIMIT 1
		mockDB.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest interface{}) error {
			// 模拟返回无押金的价格规则
			rules := dest.(*data.LockerPricingRules)
			*rules = data.LockerPricingRules{
				IsDepositEnabled: 0,   // 禁用押金（关键测试点）
				DepositAmount:    0.0, // 押金金额为0
				HourlyRate:       5.0, // 每小时5元
			}
			return nil
		})

		// 第二步：设置事务处理
		mockDB.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fn func(tx data.DBInterface) error) error {
			mockTx := data.NewMockDBInterface(ctrl)

			// 查询可用柜子
			mockTx.EXPECT().Table("lockers").Return(mockTx)                  // 选择柜子表
			mockTx.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockTx) // WHERE locker_point_id = ?
			mockTx.EXPECT().Select("id").Return(mockTx)                      // SELECT id
			mockTx.EXPECT().Where("status = 1").Return(mockTx)               // WHERE status = 1
			mockTx.EXPECT().Limit(1).Return(mockTx)                          // LIMIT 1
			mockTx.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest interface{}) error {
				// 模拟找到可用柜子
				locker := dest.(*data.Locker)
				*locker = data.Locker{Id: 1001}
				return nil
			})

			// 更新柜子状态
			mockTx.EXPECT().Table("lockers").Return(mockTx)                                      // 选择柜子表
			mockTx.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockTx)                     // WHERE id = 1001
			mockTx.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&gorm.DB{RowsAffected: 1}) // UPDATE status = 2

			// 创建订单 - 这里是关键的验证点
			mockTx.EXPECT().Table("locker_orders").Return(mockTx) // 选择订单表
			mockTx.EXPECT().Create(gomock.Any()).DoAndReturn(func(order interface{}) error {
				// 验证订单价格计算是否正确（无押金情况）
				lockerOrder := order.(*data.LockerOrders)
				expectedPrice := 5.0 * 2 // 5元/小时 * 2小时 = 10元（不包含押金）
				assert.Equal(t, expectedPrice, lockerOrder.Price, "无押金情况下价格计算应该正确")
				return nil
			})

			return fn(mockTx)
		})

		// 构造请求参数
		req := &pb.CreateDepositRequest{
			CabinetId:         1, // 网点ID
			LockerType:        1, // 柜子类型
			ScheduledDuration: 2, // 预计寄存2小时
		}

		// 调用服务方法
		resp, err := service.CreateDeposit(createTestContext("123"), req)

		// 验证结果：应该成功创建订单
		assert.NoError(t, err)                 // 没有系统错误
		assert.Equal(t, int32(200), resp.Code) // 成功状态码
	})
}

// BenchmarkDepositService_CreateDeposit 性能基准测试
// 用于测试CreateDeposit方法的性能表现，可以帮助发现性能瓶颈
func BenchmarkDepositService_CreateDeposit(b *testing.B) {
	// 创建gomock控制器
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// 创建mock对象和服务实例
	mockDB := data.NewMockDBInterface(ctrl)
	serviceData := &data.Data{DBI: data.GetDBInterface(mockDB)}
	service := NewDepositService(serviceData, &conf.Server{}, nil)

	// 设置基本的mock期望，使用AnyTimes()允许无限次调用
	// 这里故意让查询失败，避免复杂的事务逻辑影响性能测试的准确性
	mockDB.EXPECT().Table(gomock.Any()).Return(mockDB).AnyTimes()                      // 允许任意表名
	mockDB.EXPECT().Where(gomock.Any(), gomock.Any()).Return(mockDB).AnyTimes()        // 允许任意WHERE条件（带参数）
	mockDB.EXPECT().Where(gomock.Any()).Return(mockDB).AnyTimes()                      // 允许任意WHERE条件（不带参数）
	mockDB.EXPECT().Limit(gomock.Any()).Return(mockDB).AnyTimes()                      // 允许任意LIMIT值
	mockDB.EXPECT().Find(gomock.Any()).Return(errors.New("benchmark test")).AnyTimes() // 模拟查询失败，快速返回

	// 构造测试请求
	req := &pb.CreateDepositRequest{
		CabinetId:         1,
		LockerType:        1,
		ScheduledDuration: 2,
	}

	// 重置计时器，排除准备工作的时间
	b.ResetTimer()

	// 执行N次测试，Go会自动调整N的值来获得稳定的测试结果
	for i := 0; i < b.N; i++ {
		service.CreateDeposit(createTestContext("123"), req)
	}
}

// createMockLockerOrder 创建模拟的寄存订单数据
// 用于测试中需要订单数据的场景，返回一个预设的订单结构体
func createMockLockerOrder() *data.LockerOrders {
	return &data.LockerOrders{
		Id:                1,                   // 订单ID
		OrderNumber:       "20240125123456123", // 订单号（格式：日期时间+用户ID）
		UserId:            123,                 // 用户ID
		StartTime:         time.Now(),          // 寄存开始时间
		ScheduledDuration: 2,                   // 计划寄存时长（2小时）
		ActualDuration:    0,                   // 实际寄存时长（初始为0）
		Price:             110.0,               // 订单总价格（押金100+时长费10）
		Discount:          0.0,                 // 优惠金额
		AmountPaid:        110.0,               // 实际支付金额
		CabinetId:         1001,                // 分配的柜子ID
		Status:            1,                   // 订单状态（1=待支付）
		CreateTime:        time.Now(),          // 创建时间
		UpdateTime:        time.Now(),          // 更新时间
		DepositStatus:     0,                   // 押金状态（0=未处理）
		LockerPointId:     1,                   // 寄存点ID
	}
}

// createMockLocker 创建模拟的柜子数据
// 用于测试中需要柜子信息的场景，返回一个可用的柜子结构体
func createMockLocker() *data.Lockers {
	return &data.Lockers{
		Id:            1001, // 柜子ID
		LockerPointId: 1,    // 所属寄存点ID
		TypeId:        1,    // 柜子类型ID（1=小柜）
		Status:        1,    // 柜子状态（1=可用）
	}
}

// createMockPricingRules 创建模拟的价格规则数据
// 用于测试中需要价格配置的场景，返回一个标准的价格规则结构体
func createMockPricingRules() *data.LockerPricingRules {
	return &data.LockerPricingRules{
		Id:               1,          // 规则ID
		NetworkId:        1,          // 网点ID
		RuleName:         "默认规则", // 规则名称
		FeeType:          1,          // 收费类型（1=计时收费）
		LockerType:       1,          // 柜子类型（1=小柜）
		FreeDuration:     0.5,        // 免费时长（0.5小时）
		IsDepositEnabled: 1,          // 是否启用押金（1=启用）
		HourlyRate:       5.0,        // 每小时费用（5元）
		DepositAmount:    100.0,      // 押金金额（100元）
		Status:           1,          // 规则状态（1=生效）
	}
}

// TestHelperFunctions 测试辅助工具函数
// 验证测试中使用的辅助函数是否正常工作
func TestHelperFunctions(t *testing.T) {
	// 测试createTestContext函数：验证能否正确创建包含JWT token的context
	t.Run("createTestContext", func(t *testing.T) {
		// 创建包含用户ID为"123"的测试context
		ctx := createTestContext("123")

		// 从context中提取JWT token
		token, ok := jwt.FromContext(ctx)
		assert.True(t, ok, "应该能从context中提取到JWT token")

		// 将token转换为MapClaims类型并验证内容
		claims, ok := token.(*jwtv5.MapClaims)
		assert.True(t, ok, "token应该是MapClaims类型")
		assert.Equal(t, "123", (*claims)["id"], "JWT claims中的用户ID应该正确")
	})

	// 测试createEmptyContext函数：验证能否正确创建空的context
	t.Run("createEmptyContext", func(t *testing.T) {
		// 创建空的context
		ctx := createEmptyContext()

		// 尝试从空context中提取JWT token，应该失败
		_, ok := jwt.FromContext(ctx)
		assert.False(t, ok, "空context中不应该包含JWT token")
	})

	// 测试mock数据创建函数：验证能否正确创建测试用的模拟数据
	t.Run("createMockData", func(t *testing.T) {
		// 测试创建模拟订单数据
		order := createMockLockerOrder()
		assert.NotNil(t, order, "模拟订单数据不应该为空")
		assert.Greater(t, order.Id, int32(0), "订单ID应该大于0")

		// 测试创建模拟柜子数据
		locker := createMockLocker()
		assert.NotNil(t, locker, "模拟柜子数据不应该为空")
		assert.Greater(t, locker.Id, int32(0), "柜子ID应该大于0")

		// 测试创建模拟价格规则数据
		rules := createMockPricingRules()
		assert.NotNil(t, rules, "模拟价格规则数据不应该为空")
		assert.Greater(t, rules.HourlyRate, 0.0, "每小时费用应该大于0")
	})
}
