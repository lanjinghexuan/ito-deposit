package data

import "time"

type Admin struct {
	Id            int32     `gorm:"column:id;type:int;comment:管理员ID;primaryKey;not null;" json:"id"`                             // 管理员ID
	Username      string    `gorm:"column:username;type:varchar(20);comment:管理员用户名;not null;" json:"username"`                   // 管理员用户名
	Password      string    `gorm:"column:password;type:varchar(100);comment:管理员密码;not null;" json:"password"`                   // 管理员密码
	RealName      string    `gorm:"column:real_name;type:varchar(20);comment:管理员真实姓名;default:NULL;" json:"real_name"`            // 管理员真实姓名
	Mobile        string    `gorm:"column:mobile;type:char(18);comment:管理员手机号;default:NULL;" json:"mobile"`                      // 管理员手机号
	Email         string    `gorm:"column:email;type:varchar(50);comment:管理员邮箱;default:NULL;" json:"email"`                      // 管理员邮箱
	Role          int8      `gorm:"column:role;type:tinyint;comment:管理员角色：1-超级管理员，2-普通管理员;not null;" json:"role"`                // 管理员角色：1-超级管理员，2-普通管理员
	Status        int8      `gorm:"column:status;type:tinyint;comment:管理员状态：1-启用，0-禁用;default:1;" json:"status"`                 // 管理员状态：1-启用，0-禁用
	LastLoginTime time.Time `gorm:"column:last_login_time;type:datetime;comment:最后登录时间;default:NULL;" json:"last_login_time"`    // 最后登录时间
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"create_time"` // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP;" json:"update_time"` // 更新时间
}

func (a *Admin) TableName() string {
	return "admin"
}
