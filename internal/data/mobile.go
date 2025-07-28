package data

import "time"

type Users struct {
	Id        int64     `gorm:"column:id;type:bigint;comment:用户ID;primaryKey;not null;" json:"id"`                           // 用户ID
	Username  string    `gorm:"column:username;type:varchar(64);comment:用户名;not null;" json:"username"`                     // 用户名
	Mobile    string    `gorm:"column:mobile;type:char(11);comment:手机号;not null;" json:"mobile"`                            // 手机号
	ImagePath string    `gorm:"column:image_path;type:varchar(200);comment:头像;default:NULL;" json:"image_path"`              // 头像
	Password  string    `gorm:"column:password;type:varchar(128);comment:密码;not null;" json:"password"`                      // 密码
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"` // 创建时间
}

func (u *Users) TableName() string {
	return "users"
}
