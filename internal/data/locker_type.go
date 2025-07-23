package data

type LockerType struct {
	Id          int32  `gorm:"column:id;type:int;comment:类型ID;primaryKey;not null;" json:"id"`                       // 类型ID
	Name        string `gorm:"column:name;type:varchar(50);comment:类型名称;not null;" json:"name"`                    // 类型名称
	Size        string `gorm:"column:size;type:varchar(50);comment:尺寸规格;not null;" json:"size"`                    // 尺寸规格
	Description string `gorm:"column:description;type:varchar(100);comment:类型描述;default:NULL;" json:"description"` // 类型描述
}

func (u *LockerType) TableName() string { return "locker_type" }
