package data

import "time"

type LockerStorages struct {
	Id         int64     `gorm:"column:id;type:bigint UNSIGNED;comment:寄存记录ID;primaryKey;" json:"id"`                                    // 寄存记录ID
	OrderId    int64     `gorm:"column:order_id;type:bigint UNSIGNED;comment:关联主订单ID;not null;" json:"order_id"`                         // 关联主订单ID
	CabinetId  int64     `gorm:"column:cabinet_id;type:int;comment:柜子ID（存放的柜子编号）;not null;" json:"cabinet_id"`                           // 柜子ID（存放的柜子编号）
	Status     int64     `gorm:"column:status;type:tinyint;comment:状态：1-寄存中，2-已取出，3-超期未取;not null;default:1;" json:"status"`             // 状态：1-寄存中，2-已取出，3-超期未取
	CreateTime time.Time `gorm:"column:create_time;type:datetime;comment:寄存时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`   // 寄存时间
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;comment:最后更新时间;not null;default:CURRENT_TIMESTAMP;" json:"update_time"` // 最后更新时间
	UserId     int64     `gorm:"column:user_id;type:int;not null;" json:"user_id"`
}
