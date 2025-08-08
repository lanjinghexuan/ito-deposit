// internal/pkg/job/scheduler.go
package job

import (
	"context"
	"github.com/robfig/cron/v3"
	"ito-deposit/internal/service"
	"log"
)

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler(jobSvc *service.AdminService) *Scheduler {
	c := cron.New(cron.WithSeconds()) // 支持到秒
	// 每 5 分钟执行一次
	_, err := c.AddFunc("*/10 * * * * *", func() {
		jobSvc.LockerStatus()
	})
	if err != nil {
		log.Fatalf("添加任务失败: %v", err)
	}
	return &Scheduler{cron: c}
}

func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("定时任务启动")
}

func (s *Scheduler) Stop(ctx context.Context) error {
	s.cron.Stop()
	log.Println("定时任务停止")
	return nil
}
