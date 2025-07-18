package data

type Users struct {
	Id     int32  `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	Name   string `gorm:"column:name;type:varchar(50);default:NULL;" json:"name"`
	Mobile string `gorm:"column:mobile;type:char(11);default:NULL;" json:"mobile"`
}
