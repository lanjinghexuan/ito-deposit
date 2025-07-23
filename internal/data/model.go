package data

import "time"

type Users struct {
	Id        int64     `gorm:"column:id;type:bigint;comment:用户ID;primaryKey;not null;" json:"id"`                         // 用户ID
	Username  string    `gorm:"column:username;type:varchar(64);comment:用户名;not null;" json:"username"`                    // 用户名
	Mobile    string    `gorm:"column:mobile;type:char(11);comment:手机号;not null;" json:"mobile"`                           // 手机号
	ImagePath string    `gorm:"column:image_path;type:varchar(200);comment:头像;default:NULL;" json:"image_path"`            // 头像
	Password  string    `gorm:"column:password;type:varchar(128);comment:密码;not null;" json:"password"`                    // 密码
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"` // 创建时间
}

func (u *Users) TableName() string {
	return "users"
}

type LockerOrders struct {
	Id                  int32     `gorm:"column:id;type:int;comment:订单ID;primaryKey;not null;" json:"id"`                                              // 订单ID
	OrderNumber         string    `gorm:"column:order_number;type:varchar(200);comment:业务订单号（唯一标识）;not null;" json:"order_number"`                     // 业务订单号（唯一标识）
	UserId              uint64    `gorm:"column:user_id;type:bigint UNSIGNED;comment:用户ID（关联用户表）;not null;" json:"user_id"`                            // 用户ID（关联用户表）
	StartTime           time.Time `gorm:"column:start_time;type:datetime;comment:寄存开始时间;default:CURRENT_TIMESTAMP;" json:"start_time"`                 // 寄存开始时间
	ScheduledDuration   int32     `gorm:"column:scheduled_duration;type:int;comment:计划寄存时长（小时）;not null;default:0;" json:"scheduled_duration"`         // 计划寄存时长（小时）
	ActualDuration      int32     `gorm:"column:actual_duration;type:int;comment:实际寄存时长（小时）;default:NULL;" json:"actual_duration"`                     // 实际寄存时长（小时）
	Price               float64   `gorm:"column:price;type:decimal(10, 2);comment:基础费用;default:0.00;" json:"price"`                                    // 基础费用
	Discount            float64   `gorm:"column:discount;type:decimal(10, 2);comment:优惠金额;default:0.00;" json:"discount"`                              // 优惠金额
	AmountPaid          float64   `gorm:"column:amount_paid;type:decimal(10, 2);comment:实付金额;not null;default:0.00;" json:"amount_paid"`               // 实付金额
	StorageLocationName string    `gorm:"column:storage_location_name;type:varchar(40);comment:寄存网点名称;not null;" json:"storage_location_name"`         // 寄存网点名称
	CabinetId           string    `gorm:"column:cabinet_id;type:varchar(20);comment:柜子ID;not null;" json:"cabinet_id"`                                 // 柜子ID
	Status              int8      `gorm:"column:status;type:tinyint;comment:订单状态：1-待支付、2-寄存中、3-已完成、4-已取消、5-超时、6-异常;not null;default:1;" json:"status"` // 订单状态：1-待支付、2-寄存中、3-已完成、4-已取消、5-超时、6-异常
	CreateTime          time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`        // 创建时间
	UpdateTime          time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;not null;default:CURRENT_TIMESTAMP;" json:"update_time"`        // 更新时间
	DepositStatus       int8      `gorm:"column:deposit_status;type:tinyint;comment:押金状态：1-已支付、2-已退还、3-已扣除;not null;" json:"deposit_status"`           // 押金状态：1-已支付、2-已退还、3-已扣除
	DeletedAt           time.Time `gorm:"column:deleted_at;type:datetime;default:NULL;" json:"deleted_at"`
	LockerPointId       int32     `gorm:"column:locker_point_id;type:int;comment:寄存点id;default:NULL;" json:"locker_point_id"` // 寄存点id
}

func (lo *LockerOrders) TableName() string {
	return "locker_orders"
}
