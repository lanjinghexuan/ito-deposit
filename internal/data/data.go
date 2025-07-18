package data

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ito/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	DB *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	db, _ := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	return &Data{
		DB: db,
	}, cleanup, nil
}
