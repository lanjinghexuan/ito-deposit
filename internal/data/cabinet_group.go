package data

import "time"

type CabinetGroup struct {
	Id              string    `gorm:"column:id;type:varchar(32);comment:柜组ID;primaryKey;not null;" json:"id"`                                                          // 柜组ID
	LocationPointId string    `gorm:"column:location_id;type:varchar(32);comment:所属寄存点ID;not null;" json:"location_id"`                                             // 所属寄存点ID
	GroupName       string    `gorm:"column:group_name;type:varchar(64);comment:柜组名称;not null;" json:"group_name"`                                                   // 柜组名称
	GroupCode       string    `gorm:"column:group_code;type:varchar(32);comment:柜组编码(可扫码);not null;" json:"group_code"`                                           // 柜组编码(可扫码)
	GroupType       string    `gorm:"column:group_type;type:enum('standard', 'refrigerated', 'oversize');comment:柜组类型;not null;default:standard;" json:"group_type"` // 柜组类型
	Status          string    `gorm:"column:status;type:enum('normal', 'abnormal', 'disabled', 'damaged');comment:柜组状态;not null;default:normal;" json:"status"`      // 柜组状态
	TotalCells      int32     `gorm:"column:total_cells;type:int;comment:总格口数;not null;" json:"total_cells"`                                                         // 总格口数
	StartNo         int32     `gorm:"column:start_no;type:int;comment:起始编号;not null;" json:"start_no"`                                                               // 起始编号
	EndNo           int32     `gorm:"column:end_no;type:int;comment:结束编号;not null;" json:"end_no"`                                                                   // 结束编号
	InstallTime     time.Time `gorm:"column:install_time;type:datetime;comment:安装时间;default:NULL;" json:"install_time"`                                              // 安装时间
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`                          // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;not null;default:CURRENT_TIMESTAMP;" json:"update_time"`                          // 更新时间
}
