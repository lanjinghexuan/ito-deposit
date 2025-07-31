package data

type LocationType struct {
	Id          int32  `gorm:"column:id;type:int;comment:类型ID;primaryKey;not null;" json:"id"`                     // 类型ID
	Name        string `gorm:"column:name;type:varchar(50);comment:类型名称;not null;" json:"name"`                    // 类型名称
	Description string `gorm:"column:description;type:varchar(100);comment:类型描述;default:NULL;" json:"description"` // 类型描述
}

func (l *LocationType) TableName() string {
	return "location_type"
}
