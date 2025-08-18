package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库模型
type Lockers struct {
	Id            int64     `gorm:"column:id;type:bigint;primaryKey;not null;" json:"id"`
	LockerPointId int64     `gorm:"column:locker_point_id;type:bigint;comment:网点ID;not null;" json:"locker_point_id"`
	TypeId        int64     `gorm:"column:type_id;type:bigint;comment:柜型ID;not null;" json:"type_id"`
	Status        int8      `gorm:"column:status;type:tinyint;comment:状态：1-可用，2-使用中，3-故障，4-超时未取，5-维护中;default:1;" json:"status"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"update_time"`
}

type LockerPoint struct {
	Id      int64  `gorm:"column:id;type:bigint;primaryKey;not null;" json:"id"`
	Name    string `gorm:"column:name;type:varchar(255);comment:网点名称;not null;" json:"name"`
	Address string `gorm:"column:address;type:varchar(500);comment:网点地址;not null;" json:"address"`
}

// 异常快递柜信息
type AbnormalLocker struct {
	ID        int64  `json:"id"`
	PointName string `json:"point_name"`
	Address   string `json:"address"`
	Status    int    `json:"status"`
	Reason    string `json:"reason"`
	Time      string `json:"time"`
}

// 通知消息
type NotificationMessage struct {
	Type      string           `json:"type"`
	Message   string           `json:"message"`
	Lockers   []AbnormalLocker `json:"lockers"`
	Timestamp time.Time        `json:"timestamp"`
}

// 简单定时任务服务
type SimpleCronService struct {
	db    *gorm.DB
	redis *redis.Client
	cron  *cron.Cron
}

// 创建服务实例
func NewSimpleCronService() *SimpleCronService {
	// 连接数据库
	dsn := "root:password@tcp(127.0.0.1:3306)/ito_deposit?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 连接Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 测试Redis连接
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("连接Redis失败:", err)
	}

	return &SimpleCronService{
		db:    db,
		redis: rdb,
		cron:  cron.New(),
	}
}

// 检查异常快递柜
func (s *SimpleCronService) CheckAbnormalLockers() {
	log.Println("开始检查异常快递柜...")

	// 查询异常状态的快递柜
	var abnormalLockers []Lockers
	err := s.db.Where("status IN (3, 4, 5)").Find(&abnormalLockers).Error
	if err != nil {
		log.Printf("查询异常快递柜失败: %v", err)
		return
	}

	if len(abnormalLockers) == 0 {
		log.Println("没有发现异常快递柜")
		return
	}

	// 获取网点信息
	var abnormalList []AbnormalLocker
	for _, locker := range abnormalLockers {
		// 查询网点信息
		var point LockerPoint
		err := s.db.Where("id = ?", locker.LockerPointId).First(&point).Error
		if err != nil {
			log.Printf("查询网点信息失败，快递柜ID: %d", locker.Id)
			continue
		}

		// 确定异常原因
		reason := s.getAbnormalReason(int(locker.Status))

		abnormalList = append(abnormalList, AbnormalLocker{
			ID:        locker.Id,
			PointName: point.Name,
			Address:   point.Address,
			Status:    int(locker.Status),
			Reason:    reason,
			Time:      time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	// 发送通知
	if len(abnormalList) > 0 {
		s.SendNotification(abnormalList)
	}
}

// 发送通知
func (s *SimpleCronService) SendNotification(abnormalLockers []AbnormalLocker) {
	message := NotificationMessage{
		Type:      "abnormal_locker_alert",
		Message:   fmt.Sprintf("发现 %d 个异常快递柜，请及时处理", len(abnormalLockers)),
		Lockers:   abnormalLockers,
		Timestamp: time.Now(),
	}

	// 序列化消息
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化通知消息失败: %v", err)
		return
	}

	// 存储到Redis
	key := fmt.Sprintf("notifications:%s", time.Now().Format("2006-01-02"))
	ctx := context.Background()
	err = s.redis.LPush(ctx, key, string(messageBytes)).Err()
	if err != nil {
		log.Printf("存储通知到Redis失败: %v", err)
		return
	}

	// 设置过期时间（7天）
	s.redis.Expire(ctx, key, 7*24*time.Hour)

	log.Printf("异常快递柜通知已发送，共 %d 个异常快递柜", len(abnormalLockers))
}

// 获取异常原因
func (s *SimpleCronService) getAbnormalReason(status int) string {
	switch status {
	case 3:
		return "故障"
	case 4:
		return "超时未取"
	case 5:
		return "维护中"
	default:
		return "未知异常"
	}
}

// 启动定时任务
func (s *SimpleCronService) Start() {
	// 每5分钟检查一次异常快递柜
	_, err := s.cron.AddFunc("*/5 * * * *", s.CheckAbnormalLockers)
	if err != nil {
		log.Fatal("添加定时任务失败:", err)
	}

	// 立即执行一次检查
	s.CheckAbnormalLockers()

	// 启动定时任务
	s.cron.Start()
	log.Println("定时任务已启动，每5分钟检查一次异常快递柜")

	// 保持程序运行
	select {}
}

// 停止定时任务
func (s *SimpleCronService) Stop() {
	s.cron.Stop()
	log.Println("定时任务已停止")
}

// 主函数
func main() {
	service := NewSimpleCronService()

	// 启动定时任务
	service.Start()
}
