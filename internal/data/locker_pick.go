package models

import "time"

type LockerPickups struct {
	Id         int64     `gorm:"column:id;type:bigint UNSIGNED;comment:临时取出记录ID;primaryKey;" json:"id"`                                  // 临时取出记录ID
	OrderId    int64     `gorm:"column:order_id;type:bigint UNSIGNED;comment:关联主订单ID;not null;" json:"order_id"`                         // 关联主订单ID
	PickupType int64     `gorm:"column:pickup_type;type:tinyint;comment:取出类型：1-用户临时取、2-管理员操作;not null;" json:"pickup_type"`              // 取出类型：1-用户临时取、2-管理员操作
	PickupTime time.Time `gorm:"column:pickup_time;type:datetime;comment:临时取出时间;not null;" json:"pickup_time"`                           // 临时取出时间
	ReturnTime time.Time `gorm:"column:return_time;type:datetime;comment:物品归还时间（为空表示未归还）;" json:"return_time"`                           // 物品归还时间（为空表示未归还）
	Operator   string    `gorm:"column:operator;type:varchar(32);comment:操作人（用户昵称或管理员账号）;" json:"operator"`                              // 操作人（用户昵称或管理员账号）
	PointId    int64     `gorm:"column:point_id;type:int;comment:柜子的ID;not null;" json:"point_id"`                                       // 柜子的ID
	Status     int64     `gorm:"column:status;type:tinyint;comment:状态：1-取出中、2-已归还、3-超期未归;not null;default:1;" json:"status"`             // 状态：1-取出中、2-已归还、3-超期未归
	CreateTime time.Time `gorm:"column:create_time;type:datetime;comment:记录创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"` // 记录创建时间
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;comment:记录更新时间;not null;default:CURRENT_TIMESTAMP;" json:"update_time"` // 记录更新时间
}
