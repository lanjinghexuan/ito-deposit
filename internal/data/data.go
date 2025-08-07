package data

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"time"

	"ito-deposit/internal/basic/pkg"
	"ito-deposit/internal/conf"

	"github.com/redis/go-redis/v9"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewCityRepo, NewNearbyRepo, pkg.NewSendSms, NewAdminRepo, NewCabinetCellRepo)

// Data .
type Data struct {
	DB    *gorm.DB      // 保持具体类型以确保项目正常运行
	Redis *redis.Client // 保持具体类型以确保项目正常运行
	// 新增用于测试的接口字段
	DBI    DBInterface
	RedisI RedisInterface
	Mq     rocketmq.Producer
}

// 添加接口适配层，不影响原有代码
// DBInterface 定义数据库操作接口
type DBInterface interface {
	Table(name string) DBInterface
	WithContext(ctx context.Context) DBInterface
	Where(query interface{}, args ...interface{}) DBInterface
	Select(fields string) DBInterface
	Limit(limit int) DBInterface
	Find(dest interface{}) error
	Transaction(fc func(tx DBInterface) error) error // 修改参数类型为接口
	Update(column string, value interface{}) *gorm.DB
	Create(value interface{}) error // 添加Create方法
}

// RedisInterface 定义Redis操作接口
type RedisInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}

// 添加接口适配实现
func (d *Data) GetDBInterface() DBInterface {
	return &dbAdapter{db: d.DB}
}

func (d *Data) GetRedisInterface() RedisInterface {
	return &redisAdapter{client: d.Redis}
}

// NewDBAdapter 用于外部包创建 DBInterface 适配器
func NewDBAdapter(db *gorm.DB) DBInterface {
	return &dbAdapter{db: db}
}

// 数据库适配层
type dbAdapter struct {
	db *gorm.DB
}

// 实现DBInterface的所有方法...
func (a *dbAdapter) Table(name string) DBInterface {
	return &dbAdapter{db: a.db.Table(name)}
}

func (a *dbAdapter) WithContext(ctx context.Context) DBInterface {
	return &dbAdapter{db: a.db.WithContext(ctx)}
}

func (a *dbAdapter) Where(query interface{}, args ...interface{}) DBInterface {
	return &dbAdapter{db: a.db.Where(query, args...)}
}

func (a *dbAdapter) Select(fields string) DBInterface {
	return &dbAdapter{db: a.db.Select(fields)}
}

func (a *dbAdapter) Limit(limit int) DBInterface {
	return &dbAdapter{db: a.db.Limit(limit)}
}

func (a *dbAdapter) Find(dest interface{}) error {
	return a.db.Find(dest).Error
}

func (d *dbAdapter) Transaction(fn func(DBInterface) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		// ✅ 强制事务中所有读操作也走主库
		txAdapter := &dbAdapter{
			db: tx.Clauses(dbresolver.Write),
		}
		return fn(txAdapter)
	})
}

func (a *dbAdapter) Update(column string, value interface{}) *gorm.DB {
	return a.db.Update(column, value)
}

// 添加Create方法实现
func (a *dbAdapter) Create(value interface{}) error {
	return a.db.Create(value).Error
}

// Redis适配层
type redisAdapter struct {
	client *redis.Client
}

// 实现RedisInterface的所有方法...
func (a *redisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return a.client.Set(ctx, key, value, expiration)
}

func (a *redisAdapter) Get(ctx context.Context, key string) *redis.StringCmd {
	return a.client.Get(ctx, key)
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {

	helper := log.NewHelper(logger)
	masterDsn := c.Database.Master
	slaveDsn := c.Database.Slave
	if masterDsn == "" {
		panic("config error: database.master is required")
	}
	if slaveDsn == "" {
		panic("config error: database.slave is required")
	}
	//主库
	db, err := gorm.Open(mysql.Open(masterDsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect to master database", err)
		panic("failed to connect database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	//从库
	err = db.Use(
		dbresolver.Register(
			dbresolver.Config{
				Replicas: []gorm.Dialector{
					mysql.Open(slaveDsn), // 从库（读）
				},
				Policy: dbresolver.RandomPolicy{}, // 读负载均衡策略
			},
		).
			SetMaxIdleConns(10).
			SetMaxOpenConns(50).
			SetConnMaxLifetime(time.Hour),
	)
	if err != nil {
		helper.Errorf("failed to register dbresolver: %v", err)
		panic(err)
	}
	if err := db.AutoMigrate(&City{}, &LockerPoint{}); err != nil {
		helper.Errorf("自动迁移数据库表结构失败: %v", err)
	} else {
		helper.Info("自动迁移数据库表结构成功")
	}
	//redis
	sqlDB, err = db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	redisDB := RedisInit(c)

	mq, err := rocketmq.NewProducer(
		producer.WithGroupName("deposit_group"),
		producer.WithNameServer([]string{"14.103.243.149:9876"}),
	)
	if err != nil {
		panic(err)
	}
	if err := mq.Start(); err != nil {
		panic(err)
	}

	cleanup := func() {
		helper.Info("closing the data resources")
		if err := mq.Shutdown(); err != nil {
			helper.Errorf("failed to shutdown mq: %v", err)
		}
	}

	consumerIns, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("deposit_group"),
		consumer.WithNameServer([]string{"14.103.243.149:9876"}),
	)
	if err != nil {
		panic(err)
	}

	err = StartDepositConsumer(consumerIns, db)
	if err != nil {
		panic(err)
	}

	if err := consumerIns.Start(); err != nil {
		panic(err)
	}

	return &Data{
		DB:    db,
		Redis: redisDB,
		Mq:    mq,
	}, cleanup, nil
}

func RedisInit(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})
	// 可选：设置连接超时 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查是否能连接到 Redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("无法连接到 Redis: %v", err))
		// 或者使用日志：log.Fatalf("无法连接到 Redis: %v", err)
	}

	return rdb
}

// 导出的接口工厂方法
func GetDBInterface(db DBInterface) DBInterface {
	return db
}

func GetRedisInterface(redis RedisInterface) RedisInterface {
	return redis
}
