package data

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"ito-deposit/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		fmt.Println("err:", err)
		panic("failed to connect database")
	}
	redisDB := RedisInit(c)
	return &Data{
		DB:    db,
		Redis: redisDB,
	}, cleanup, nil
}

func RedisInit(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:            c.Redis.Addr,
		DisableIdentity: true,
	})
	return rdb
}
