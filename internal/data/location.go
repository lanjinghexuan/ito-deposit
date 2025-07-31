package data

type Location struct {
	Id         int32   `gorm:"column:id;type:int;comment:地点ID;primaryKey;not null;" json:"id"`                     // 地点ID
	Name       string  `gorm:"column:name;type:varchar(100);comment:地点名称;not null;" json:"name"`                   // 地点名称
	CityId     int32   `gorm:"column:city_id;type:int;comment:所属城市ID;not null;" json:"city_id"`                    // 所属城市ID
	TypeId     int32   `gorm:"column:type_id;type:int;comment:地点类型ID;default:NULL;" json:"type_id"`                // 地点类型ID
	Address    string  `gorm:"column:address;type:varchar(50);comment:详细地址;default:NULL;" json:"address"`          // 详细地址
	Latitude   float64 `gorm:"column:latitude;type:decimal(10, 6);comment:纬度;not null;" json:"latitude"`           // 纬度
	Longitude  float64 `gorm:"column:longitude;type:decimal(10, 6);comment:经度;not null;" json:"longitude"`         // 经度
	IsHot      int8    `gorm:"column:is_hot;type:tinyint(1);comment:是否热门地点，1=是，0=否;default:0;" json:"is_hot"`      // 是否热门地点，1=是，0=否
	CustomType string  `gorm:"column:custom_type;type:varchar(50);comment:自定义类型;default:NULL;" json:"custom_type"` // 自定义类型
}

func (l *Location) TableName() string {
	return "location"
}
