package data

import "time"

// TODO: 城市表
type City struct {
	Id          int32     `gorm:"column:id;type:int;comment:城市ID，主键自增;primaryKey;not null;" json:"id"`                                     //TODO: 城市ID，主键自增
	Name        string    `gorm:"column:name;type:varchar(50);comment:城市名称;not null;" json:"name"`                                         // 城市名称（如"郑州市"）
	Province    string    `gorm:"column:province;type:varchar(50);comment:所属省份;not null;" json:"province"`                                 // 所属省份（如"河南省"）
	CityCode    string    `gorm:"column:city_code;type:varchar(20);comment:城市行政区划代码;not null;" json:"city_code"`                           // 城市行政区划代码
	Status      int8      `gorm:"column:status;type:tinyint(1);comment:状态：1-服务开通 0-服务未开通;not null;default:1;" json:"status"`               // 状态：1-服务开通 0-服务未开通
	IsSupported int8      `gorm:"column:is_supported;type:tinyint(1);comment:是否支持寄存服务：1-支持 0-不支持;not null;default:0;" json:"is_supported"` // 是否支持寄存服务：1-支持 0-不支持
	Longitude   float64   `gorm:"column:longitude;type:decimal(10, 6);comment:城市中心点经度;default:NULL;" json:"longitude"`                     // 城市中心点经度
	Latitude    float64   `gorm:"column:latitude;type:decimal(10, 6);comment:城市中心点纬度;default:NULL;" json:"latitude"`                       // 城市中心点纬度
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`    // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;not null;default:CURRENT_TIMESTAMP;" json:"update_time"`    // 更新时间
}

// TODO: 用户搜索记录表
type SearchHistory struct {
	Id          int32     `gorm:"column:id;type:int;comment:记录ID，主键自增;primaryKey;not null;" json:"id"`                                  //TODO: 记录ID，主键自增
	UserId      int32     `gorm:"column:user_id;type:int;comment:用户ID（未登录为NULL）;default:NULL;" json:"user_id"`                          // 用户ID（未登录为NULL）
	SessionId   string    `gorm:"column:session_id;type:varchar(100);comment:会话ID;not null;" json:"session_id"`                         // 会话ID
	CityId      int32     `gorm:"column:city_id;type:int;comment:搜索所在城市ID;not null;" json:"city_id"`                                    // 搜索所在城市ID
	Keyword     string    `gorm:"column:keyword;type:varchar(100);comment:搜索关键词;not null;" json:"keyword"`                              // 搜索关键词
	SearchType  int8      `gorm:"column:search_type;type:tinyint;comment:搜索类型：1-火车站 2-地铁站 3-景点 4-商圈;default:NULL;" json:"search_type"`  // 搜索类型：1-火车站 2-地铁站 3-景点 4-商圈
	ResultCount int32     `gorm:"column:result_count;type:int;comment:搜索结果数量;default:NULL;" json:"result_count"`                        // 搜索结果数量
	DeviceInfo  string    `gorm:"column:device_info;type:varchar(255);comment:设备信息;default:NULL;" json:"device_info"`                   // 设备信息
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"` // 创建时间
}

// TODO: 寄存点位信息表
type StoragePoint struct {
	Id              int32     `gorm:"column:id;type:int;comment:寄存点ID，主键自增;primaryKey;not null;" json:"id"`                                 //TODO: 寄存点ID，主键自增
	CityId          int32     `gorm:"column:city_id;type:int;comment:所在城市ID;not null;" json:"city_id"`                                      // 所在城市ID
	Name            string    `gorm:"column:name;type:varchar(100);comment:寄存点名称（如"中交锦兰荟南门处丰巢柜"）;not null;" json:"name"`                    //TODO: 寄存点名称（如"中交锦兰荟南门处丰巢柜"）
	Address         string    `gorm:"column:address;type:varchar(255);comment:详细地址;not null;" json:"address"`                               // 详细地址
	Longitude       float64   `gorm:"column:longitude;type:decimal(10, 6);comment:经度坐标;not null;" json:"longitude"`                         // 经度坐标
	Latitude        float64   `gorm:"column:latitude;type:decimal(10, 6);comment:纬度坐标;not null;" json:"latitude"`                           // 纬度坐标
	Type            int8      `gorm:"column:type;type:tinyint;comment:点位类型：1-火车站 2-地铁站 3-景点 4-商圈 5-社区;not null;" json:"type"`               // 点位类型：1-火车站 2-地铁站 3-景点 4-商圈 5-社区
	TotalLarge      int32     `gorm:"column:total_large;type:int;comment:大格口总数;not null;default:0;" json:"total_large"`                     // 大格口总数
	TotalMedium     int32     `gorm:"column:total_medium;type:int;comment:中格口总数;not null;default:0;" json:"total_medium"`                   // 中格口总数
	TotalSmall      int32     `gorm:"column:total_small;type:int;comment:小格口总数;not null;default:0;" json:"total_small"`                     // 小格口总数
	AvailableLarge  int32     `gorm:"column:available_large;type:int;comment:可用大格口数;not null;default:0;" json:"available_large"`            // 可用大格口数
	AvailableMedium int32     `gorm:"column:available_medium;type:int;comment:可用中格口数;not null;default:0;" json:"available_medium"`          // 可用中格口数
	AvailableSmall  int32     `gorm:"column:available_small;type:int;comment:可用小格口数;not null;default:0;" json:"available_small"`            // 可用小格口数
	Status          int8      `gorm:"column:status;type:tinyint;comment:状态：1-正常 0-离线 2-维护中;not null;default:1;" json:"status"`              // 状态：1-正常 0-离线 2-维护中
	Distance        int32     `gorm:"column:distance;type:int;comment:距离市中心距离(米);default:NULL;" json:"distance"`                            // 距离市中心距离(米)
	Operator        string    `gorm:"column:operator;type:varchar(50);comment:运营商;default:NULL;" json:"operator"`                           // 运营商
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"` // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;not null;default:CURRENT_TIMESTAMP;" json:"update_time"` // 更新时间
}

// TODO: 用户位置记录表
type UserLocation struct {
	Id           int32     `gorm:"column:id;type:int;comment:记录ID，主键自增;primaryKey;not null;" json:"id"`                                  //TODO: 记录ID，主键自增
	UserId       int32     `gorm:"column:user_id;type:int;comment:用户ID（未登录为NULL）;default:NULL;" json:"user_id"`                          // 用户ID（未登录为NULL）
	SessionId    string    `gorm:"column:session_id;type:varchar(100);comment:会话ID;not null;" json:"session_id"`                         // 会话ID
	CityId       int32     `gorm:"column:city_id;type:int;comment:当前所在城市ID;not null;" json:"city_id"`                                    // 当前所在城市ID
	Address      string    `gorm:"column:address;type:varchar(255);comment:格式化地址;default:NULL;" json:"address"`                          // 格式化地址
	Longitude    float64   `gorm:"column:longitude;type:decimal(10, 6);comment:经度坐标;not null;" json:"longitude"`                         // 经度坐标
	Latitude     float64   `gorm:"column:latitude;type:decimal(10, 6);comment:纬度坐标;not null;" json:"latitude"`                           // 纬度坐标
	LocationTime time.Time `gorm:"column:location_time;type:datetime;comment:定位时间;not null;" json:"location_time"`                       // 定位时间
	IpAddress    string    `gorm:"column:ip_address;type:varchar(50);comment:IP地址;default:NULL;" json:"ip_address"`                      // IP地址
	DeviceInfo   string    `gorm:"column:device_info;type:varchar(255);comment:设备信息;default:NULL;" json:"device_info"`                   // 设备信息
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"` // 创建时间
}
