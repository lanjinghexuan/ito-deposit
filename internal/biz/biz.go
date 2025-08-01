package biz

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"time"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase, NewAdminUsecase, NewCityUsecase, NewNearbyUsecase, NewCabinetCellUsecase)

type Locker struct {
	Id            int32 `gorm:"column:id;type:int;comment:寄存柜ID;primaryKey;not null;" json:"id"`                    // 寄存柜ID
	LockerPointId int32 `gorm:"column:locker_point_id;type:int;comment:所属寄存点ID;not null;" json:"locker_point_id"`   // 所属寄存点ID
	TypeId        int32 `gorm:"column:type_id;type:int;comment:柜型ID;not null;" json:"type_id"`                      // 柜型ID
	Status        int8  `gorm:"column:status;type:tinyint;comment:寄存柜状态:1=可用,2=已占用,3=维护中;default:1;" json:"status"` // 寄存柜状态:1=可用,2=已占用,3=维护中
}

type LockerOrders struct {
	Id                  int32          `gorm:"column:id;type:int;comment:订单ID;primaryKey;not null;" json:"id"`                                              // 订单ID
	OrderNumber         string         `gorm:"column:order_number;type:varchar(200);comment:业务订单号（唯一标识）;not null;" json:"order_number"`                     // 业务订单号（唯一标识）
	UserId              uint64         `gorm:"column:user_id;type:bigint UNSIGNED;comment:用户ID（关联用户表）;not null;" json:"user_id"`                            // 用户ID（关联用户表）
	StartTime           time.Time      `gorm:"column:start_time;type:datetime;comment:寄存开始时间;default:CURRENT_TIMESTAMP;" json:"start_time"`                 // 寄存开始时间
	ScheduledDuration   int32          `gorm:"column:scheduled_duration;type:int;comment:计划寄存时长（小时）;not null;default:0;" json:"scheduled_duration"`         // 计划寄存时长（小时）
	ActualDuration      int32          `gorm:"column:actual_duration;type:int;comment:实际寄存时长（小时）;default:NULL;" json:"actual_duration"`                     // 实际寄存时长（小时）
	Price               float64        `gorm:"column:price;type:decimal(10, 2);comment:基础费用;default:0.00;" json:"price"`                                    // 基础费用
	Discount            float64        `gorm:"column:discount;type:decimal(10, 2);comment:优惠金额;default:0.00;" json:"discount"`                              // 优惠金额
	AmountPaid          float64        `gorm:"column:amount_paid;type:decimal(10, 2);comment:实付金额;not null;default:0.00;" json:"amount_paid"`               // 实付金额
	StorageLocationName string         `gorm:"column:storage_location_name;type:varchar(40);comment:寄存网点名称;not null;" json:"storage_location_name"`         // 寄存网点名称
	CabinetId           int32          `gorm:"column:cabinet_id;type:int;comment:柜子ID;not null;" json:"cabinet_id"`                                         // 柜子ID
	Status              int8           `gorm:"column:status;type:tinyint;comment:订单状态：1-待支付、2-寄存中、3-已完成、4-已取消、5-超时、6-异常;not null;default:1;" json:"status"` // 订单状态：1-待支付、2-寄存中、3-已完成、4-已取消、5-超时、6-异常
	CreateTime          time.Time      `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`        // 创建时间
	UpdateTime          time.Time      `gorm:"column:update_time;type:datetime;comment:更新时间;not null;default:CURRENT_TIMESTAMP;" json:"update_time"`        // 更新时间
	DepositStatus       int8           `gorm:"column:deposit_status;type:tinyint;comment:押金状态：1-已支付、2-已退还、3-已扣除;not null;" json:"deposit_status"`           // 押金状态：1-已支付、2-已退还、3-已扣除
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;default:NULL;" json:"deleted_at"`
	LockerPointId       int32          `gorm:"column:locker_point_id;type:int;comment:寄存点id;default:NULL;" json:"locker_point_id"` // 寄存点id
	Title               string         `gorm:"column:title;type:varchar(50);not null;" json:"title"`
}

type LockerPoint struct {
	Id              int32   `gorm:"column:id;type:int;comment:寄存点ID;primaryKey;not null;" json:"id"`                       // 寄存点ID
	LocationId      int32   `gorm:"column:location_id;type:int;comment:所属地点ID;not null;" json:"location_id"`               // 所属地点ID
	CabinetGroupId  int32   `gorm:"column:cabinet_group_id;type:int;comment:寄存柜组ID;default:NULL;" json:"cabinet_group_id"` // 寄存柜组ID
	Name            string  `gorm:"column:name;type:varchar(30);comment:寄存点名称;not null;" json:"name"`                      // 寄存点名称
	Address         string  `gorm:"column:address;type:varchar(50);comment:详细地址;default:NULL;" json:"address"`             // 详细地址
	Latitude        float64 `gorm:"column:latitude;type:decimal(11, 6);comment:纬度;not null;" json:"latitude"`              // 纬度
	Longitude       float64 `gorm:"column:longitude;type:decimal(11, 6);comment:经度;not null;" json:"longitude"`            // 经度
	AvailableLarge  int32   `gorm:"column:available_large;type:int;comment:可用大柜数量;default:0;" json:"available_large"`      // 可用大柜数量
	AvailableMedium int32   `gorm:"column:available_medium;type:int;comment:可用中柜数量;default:0;" json:"available_medium"`    // 可用中柜数量
	AvailableSmall  int32   `gorm:"column:available_small;type:int;comment:可用小柜数量;default:0;" json:"available_small"`      // 可用小柜数量
	OpenTime        string  `gorm:"column:open_time;type:varchar(30);comment:营业时间;default:NULL;" json:"open_time"`         // 营业时间
	Mobile          string  `gorm:"column:mobile;type:varchar(20);comment:联系电话;default:NULL;" json:"mobile"`               // 联系电话
	AdminId         int32   `gorm:"column:admin_id;type:int;comment:管理员id;default:NULL;" json:"admin_id"`                  // 管理员id
	PointImage      string  `gorm:"column:point_image;type:varchar(200);comment:网点图片;default:NULL;" json:"point_image"`    // 网点图片
	Status          string  `gorm:"column:status;type:varchar(30);comment:状态;default:NULL;" json:"status"`                 // 状态
	PointType       string  `gorm:"column:point_type;type:varchar(20);comment:网点类型;default:NULL;" json:"point_type"`       // 网点类型
}

func (l LockerPoint) TableName() string {
	return "locker_point"
}

type LockerPricingRules struct {
	Id               int64     `gorm:"column:id;type:bigint;primaryKey;not null;" json:"id"`
	NetworkId        int64     `gorm:"column:network_id;type:bigint;comment:网点ID;not null;" json:"network_id"`                              // 网点ID
	RuleName         string    `gorm:"column:rule_name;type:varchar(50);comment:规则名称;default:默认规则;" json:"rule_name"`                       // 规则名称
	FeeType          int8      `gorm:"column:fee_type;type:tinyint;comment:1-计时收费 2-按日收费;not null;" json:"fee_type"`                        // 1-计时收费 2-按日收费
	LockerType       int8      `gorm:"column:locker_type;type:tinyint;comment:1-小柜子 2-大柜子;not null;" json:"locker_type"`                    // 1-小柜子 2-大柜子
	FreeDuration     float64   `gorm:"column:free_duration;type:decimal(5, 1);comment:免费时长(小时);not null;default:0.0;" json:"free_duration"` // 免费时长(小时)
	IsDepositEnabled int8      `gorm:"column:is_deposit_enabled;type:tinyint;comment:是否启用押金;default:0;" json:"is_deposit_enabled"`          // 是否启用押金
	IsAdvancePay     int8      `gorm:"column:is_advance_pay;type:tinyint;comment:是否启用预付;default:0;" json:"is_advance_pay"`                  // 是否启用预付
	HourlyRate       float64   `gorm:"column:hourly_rate;type:decimal(10, 2);comment:每小时费用;default:0.00;" json:"hourly_rate"`               // 每小时费用
	DailyCap         float64   `gorm:"column:daily_cap;type:decimal(10, 2);comment:24小时封顶价;default:NULL;" json:"daily_cap"`                 // 24小时封顶价
	DailyRate        float64   `gorm:"column:daily_rate;type:decimal(10, 2);comment:每日费用;default:0.00;" json:"daily_rate"`                  // 每日费用
	AdvanceAmount    float64   `gorm:"column:advance_amount;type:decimal(10, 2);comment:预付金额;default:0.00;" json:"advance_amount"`          // 预付金额
	DepositAmount    float64   `gorm:"column:deposit_amount;type:decimal(10, 2);comment:押金金额;default:0.00;" json:"deposit_amount"`          // 押金金额
	Status           int8      `gorm:"column:status;type:tinyint;comment:1-生效 0-失效;default:1;" json:"status"`                               // 1-生效 0-失效
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;" json:"updated_at"`
}
