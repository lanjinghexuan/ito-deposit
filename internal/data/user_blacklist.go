package data

import "time"

type UserBlacklist struct {
	Id         int64     `gorm:"column:id;type:bigint;comment:主键ID;primaryKey;not null;" json:"id"`                                        // 主键ID
	UserId     int64     `gorm:"column:user_id;type:bigint;comment:用户ID;not null;" json:"user_id"`                                         // 用户ID
	IsActive   int8      `gorm:"column:is_active;type:tinyint(1);comment:是否有效;not null;default:1;" json:"is_active"`                     // 是否有效
	BanLevel   int8      `gorm:"column:ban_level;type:tinyint;comment:封禁级别(1:部分限制 2:完全封禁);not null;default:1;" json:"ban_level"` // 封禁级别(1:部分限制 2:完全封禁)
	Reason     string    `gorm:"column:reason;type:varchar(200);comment:封禁原因;not null;" json:"reason"`                                   // 封禁原因
	OperatorId int64     `gorm:"column:operator_id;type:bigint;comment:操作人ID;default:NULL;" json:"operator_id"`                           // 操作人ID
	StartTime  time.Time `gorm:"column:start_time;type:timestamp;comment:开始时间;not null;default:CURRENT_TIMESTAMP;" json:"start_time"`    // 开始时间
	EndTime    time.Time `gorm:"column:end_time;type:timestamp;comment:结束时间(空表示永久);default:NULL;" json:"end_time"`                  // 结束时间(空表示永久)
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"created_at"`    // 创建时间
}

func (UserBlacklist) TableName() string {
	return "user_blacklist"
}
