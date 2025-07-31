package data

import (
	"time"
)

// City 城市数据模型，用于存储城市基本信息
type City struct {
	ID        int32     `gorm:"column:id;type:int;primaryKey;autoIncrement;comment:城市ID;not null;" json:"id"`                                                    // 城市ID
	Name      string    `gorm:"column:name;type:varchar(50);comment:城市名称;not null;" json:"name"`                                                                 // 城市名称
	Code      string    `gorm:"column:code;type:varchar(20);comment:城市编码;not null;uniqueIndex;" json:"code"`                                                     // 城市编码
	Latitude  float64   `gorm:"column:latitude;type:decimal(11, 6);comment:纬度;not null;" json:"latitude"`                                                        // 纬度
	Longitude float64   `gorm:"column:longitude;type:decimal(11, 6);comment:经度;not null;" json:"longitude"`                                                      // 经度
	Status    int8      `gorm:"column:status;type:tinyint(1);comment:状态(1:启用,0:禁用);default:1;" json:"status"`                                                   // 状态(1:启用,0:禁用)
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;comment:创建时间;default:CURRENT_TIMESTAMP;not null;" json:"created_at"`                             // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;not null;" json:"updated_at"` // 更新时间
}

// TableName 设置表名
func (c *City) TableName() string {
	return "city"
}
