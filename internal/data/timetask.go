package data

import "time"

type TimerTasks struct {
	Id              uint64    `gorm:"column:id;type:bigint UNSIGNED;comment:任务ID;primaryKey;" json:"id"`                                                      // 任务ID
	TaskCode        string    `gorm:"column:task_code;type:varchar(64);comment:任务唯一标识（如：storage_remind_123）;not null;" json:"task_code"`                      // 任务唯一标识（如：storage_remind_123）
	OrderId         uint64    `gorm:"column:order_id;type:bigint UNSIGNED;comment:关联的订单ID;not null;" json:"order_id"`                                         // 关联的订单ID
	LockerStorageId uint64    `gorm:"column:locker_storage_id;type:bigint UNSIGNED;comment:关联的寄存记录ID（关联locker_storages表）;not null;" json:"locker_storage_id"` // 关联的寄存记录ID（关联locker_storages表）
	TaskType        string    `gorm:"column:task_type;type:varchar(20);" json:"task_type"`
	CronExpression  string    `gorm:"column:cron_expression;type:varchar(64);comment:cron表达式（重复任务用）;" json:"cron_expression"`                   // cron表达式（重复任务用）
	ExecuteTime     time.Time `gorm:"column:execute_time;type:datetime;comment:计划执行时间;not null;" json:"execute_time"`                           // 计划执行时间
	Status          int8      `gorm:"column:status;type:tinyint;comment:任务状态：0-待执行，1-执行中，2-已完成，3-执行失败，4-已取消;not null;default:0;" json:"status"` // 任务状态：0-待执行，1-执行中，2-已完成，3-执行失败，4-已取消
	RetryCount      int8      `gorm:"column:retry_count;type:tinyint;comment:重试次数;not null;default:0;" json:"retry_count"`                      // 重试次数
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"created_at"`       // 创建时间
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`       // 更新时间
	Remark          string    `gorm:"column:remark;type:varchar(255);comment:备注;" json:"remark"`                                                // 备注
}
