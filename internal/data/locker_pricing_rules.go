package data

import "time"

type LockerPricingRules struct {
	Id               int64     `gorm:"column:id;type:bigint;primaryKey;not null;" json:"id"`
	NetworkId        int64     `gorm:"column:network_id;type:bigint;comment:网点ID;not null;" json:"network_id"`                              // 网点ID
	RuleName         string    `gorm:"column:rule_name;type:varchar(50);comment:规则名称;default:默认规则;" json:"rule_name"`                       // 规则名称
	FeeType          int8      `gorm:"column:fee_type;type:tinyint;comment:1-计时收费 2-按日收费;not null;" json:"fee_type"`                        // 1-计时收费 2-按日收费
	LockerType       int8      `gorm:"column:locker_type;type:tinyint;comment:1-小柜子 2-大柜子;not null;" json:"locker_type"`                    // 1-小柜子 2-大柜子
	FreeDuration     float64   `gorm:"column:free_duration;type:decimal(5, 1);comment:免费时长(小时);not null;default:0.0;" json:"free_duration"` // 免费时长(小时)
	IsDepositEnabled int8      `gorm:"column:is_deposit_enabled;type:tinyint;comment:是否启用押金;default:0;" json:"is_deposit_enabled"`          // 是否启用押金
	IsAdvancePay     int8      `gorm:"column:is_advance_pay;type:tinyint;comment:是否启用预付;default:0;" json:"is_advance_pay"`                  // 是否启用预付
	HourlyRate       float64   `gorm:"column:hourly_rate;type:decimal(10, 2);comment:每小时费用;default:0.00;" json:"hourly_rate"`               // 每小时费用
	DailyCap         float64   `gorm:"column:daily_cap;type:decimal(10, 2);comment:24小时封顶价;default:NULL;" json:"daily_cap"`                 // 24小时封顶价
	DailyRate        float64   `gorm:"column:daily_rate;type:decimal(10, 2);comment:每日费用;default:0.00;" json:"daily_rate"`                  // 每日费用
	AdvanceAmount    float64   `gorm:"column:advance_amount;type:decimal(10, 2);comment:预付金额;default:0.00;" json:"advance_amount"`          // 预付金额
	DepositAmount    float64   `gorm:"column:deposit_amount;type:decimal(10, 2);comment:押金金额;default:0.00;" json:"deposit_amount"`          // 押金金额
	Status           int8      `gorm:"column:status;type:tinyint;comment:1-生效 0-失效;default:1;" json:"status"`                               // 1-生效 0-失效
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;" json:"updated_at"`
}
