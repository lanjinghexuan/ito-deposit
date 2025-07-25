package data

type Lockers struct {
	Id            int32 `gorm:"column:id;type:int;comment:寄存柜ID;primaryKey;" json:"id"`                             // 寄存柜ID
	LockerPointId int32 `gorm:"column:locker_point_id;type:int;comment:所属寄存点ID;not null;" json:"locker_point_id"`   // 所属寄存点ID
	TypeId        int32 `gorm:"column:type_id;type:int;comment:柜型ID;not null;" json:"type_id"`                      // 柜型ID
	Status        int8  `gorm:"column:status;type:tinyint;comment:寄存柜状态:1=可用,2=已占用,3=维护中;default:1;" json:"status"` // 寄存柜状态:1=可用,2=已占用,3=维护中
}
