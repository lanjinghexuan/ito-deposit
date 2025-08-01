package data

import "time"

type CabinetCell struct {
	Id           string    `gorm:"column:id;type:varchar(32);comment:格口ID;primaryKey;not null;" json:"id"`                                                              // 格口ID
	GroupId      string    `gorm:"column:group_id;type:varchar(32);comment:所属柜组ID;not null;" json:"group_id"`                                                         // 所属柜组ID
	CellNo       int32     `gorm:"column:cell_no;type:int;comment:格口编号;not null;" json:"cell_no"`                                                                     // 格口编号
	Status       string    `gorm:"column:status;type:enum('normal', 'inUse', 'abnormal', 'disabled', 'damaged');comment:格口状态;not null;default:normal;" json:"status"` // 格口状态
	LastOpenTime time.Time `gorm:"column:last_open_time;type:datetime;comment:最后开启时间;default:NULL;" json:"last_open_time"`                                          // 最后开启时间
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`                              // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;not null;default:CURRENT_TIMESTAMP;" json:"update_time"`                              // 更新时间
}
