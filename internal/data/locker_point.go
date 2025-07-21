package data

type LockerPoint struct {
	Id              int32   `gorm:"column:id;type:int;comment:寄存点ID;primaryKey;not null;" json:"id"`                       // 寄存点ID
	LocationId      int32   `gorm:"column:location_id;type:int;comment:所属地点ID;not null;" json:"location_id"`              // 所属地点ID
	Name            string  `gorm:"column:name;type:varchar(30);comment:寄存点名称;not null;" json:"name"`                    // 寄存点名称
	Address         string  `gorm:"column:address;type:varchar(50);comment:详细地址;default:NULL;" json:"address"`            // 详细地址
	Latitude        float64 `gorm:"column:latitude;type:decimal(11, 6);comment:纬度;not null;" json:"latitude"`               // 纬度
	Longitude       float64 `gorm:"column:longitude;type:decimal(11, 6);comment:经度;not null;" json:"longitude"`             // 经度
	AvailableLarge  int32   `gorm:"column:available_large;type:int;comment:可用大柜数量;default:0;" json:"available_large"`   // 可用大柜数量
	AvailableMedium int32   `gorm:"column:available_medium;type:int;comment:可用中柜数量;default:0;" json:"available_medium"` // 可用中柜数量
	AvailableSmall  int32   `gorm:"column:available_small;type:int;comment:可用小柜数量;default:0;" json:"available_small"`   // 可用小柜数量
	OpenTime        string  `gorm:"column:open_time;type:varchar(30);comment:营业时间;default:NULL;" json:"open_time"`        // 营业时间
	Mobile          string  `gorm:"column:mobile;type:varchar(20);comment:联系电话;default:NULL;" json:"mobile"`              // 联系电话
}

func (u *LockerPoint) TableName() string {
	return "locker_point"
}
