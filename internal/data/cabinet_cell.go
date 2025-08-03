package data

import "time"

type CabinetCell struct {
	Id             int32     `gorm:"column:id;type:int;comment:格口ID;primaryKey;not null;" json:"id"`                                                                      // 格口ID
	CabinetGroupId int32     `gorm:"column:cabinet_group_id;type:int;comment:所属柜组ID;not null;" json:"cabinet_group_id"`                                                 // 所属柜组ID
	CellNo         int32     `gorm:"column:cell_no;type:int;comment:格口编号;not null;" json:"cell_no"`                                                                     // 格口编号
	Status         string    `gorm:"column:status;type:enum('normal', 'inUse', 'abnormal', 'disabled', 'damaged');comment:格口状态;not null;default:normal;" json:"status"` // 格口状态
	LastOpenTime   time.Time `gorm:"column:last_open_time;type:datetime;comment:最后开启时间;default:NULL;" json:"last_open_time"`                                          // 最后开启时间
	CreateTime     time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`                              // 创建时间
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;not null;default:CURRENT_TIMESTAMP;" json:"update_time"`                              // 更新时间
	DeletedAt      time.Time `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
	CellSize       string    `gorm:"column:cell_size;type:enum('small', 'medium', 'large');comment:格口大小;not null;default:medium;" json:"cell_size"` // 格口大小
}

func (CabinetCell) TableName() string {
	return "cabinet_cells"
}
