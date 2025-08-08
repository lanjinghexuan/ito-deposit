package data

import "time"

// Lockers 快递柜表
type Lockers struct {
	Id            int64     `gorm:"column:id;type:bigint;primaryKey;not null;" json:"id"`
	LockerPointId int64     `gorm:"column:locker_point_id;type:bigint;comment:网点ID;not null;" json:"locker_point_id"`
	TypeId        int64     `gorm:"column:type_id;type:bigint;comment:柜型ID;not null;" json:"type_id"`
	Status        int8      `gorm:"column:status;type:tinyint;comment:状态：1-可用，2-使用中，3-故障，4-超时未取，5-维护中;default:1;" json:"status"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"update_time"`
}

func (l *Lockers) TableName() string {
	return "lockers"
}
