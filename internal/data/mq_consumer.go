package data

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"
	"ito-deposit/internal/biz"
	"log"
	"time"
)

type DepositOrderMessage struct {
	OrderNo           string  `json:"order_no"`
	UserID            int64   `json:"user_id"`
	ScheduledDuration int32   `json:"scheduled_duration"`
	Price             float64 `json:"price"`
	LockerID          int64   `json:"locker_id"`
}

func StartDepositConsumer(mq rocketmq.PushConsumer, db *gorm.DB) error {
	return mq.Subscribe("deposit_create", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			var orderMsg DepositOrderMessage
			if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
				log.Println("MQ消息解析失败:", err)
				continue
			}

			log.Printf("收到订单消息: %+v\n", orderMsg)

			// ✅ 先查 local_message 做幂等判断
			var localMessage biz.LocalMessage
			if err := db.Table("local_message").Where("business_key = ?", orderMsg.OrderNo).First(&localMessage).Error; err != nil {
				log.Println("查询 local_message 失败:", err)
				return consumer.ConsumeRetryLater, nil
			}

			if localMessage.Status == 2 {
				log.Println("消息已被消费，跳过:", orderMsg.OrderNo)
				return consumer.ConsumeSuccess, nil
			}

			// ✅ 开始插入历史记录
			record := biz.DepositRecord{
				OrderNo:   orderMsg.OrderNo,
				UserId:    int32(orderMsg.UserID),
				LockerId:  int32(orderMsg.LockerID),
				Price:     orderMsg.Price,
				CreatedAt: time.Now(),
			}
			if err := db.Table("deposit_record").Create(&record).Error; err != nil {
				log.Println("写入 deposit_record 失败:", err)
				return consumer.ConsumeRetryLater, nil
			}

			// ✅ 更新 local_message 状态为已消费（2）
			if err := db.Table("local_message").Where("business_key = ?", orderMsg.OrderNo).
				Update("status", 2).Error; err != nil {
				log.Println("更新 local_message 状态失败:", err)
			}

			log.Println("消费完成，订单号:", orderMsg.OrderNo)
		}
		return consumer.ConsumeSuccess, nil
	})
}
