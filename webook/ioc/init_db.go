package ioc

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/spigcoder/LittleBook/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitDB() *gorm.DB {
	type config struct {
		Dsn string
	}
	var conf config
	err := viper.UnmarshalKey("db.mysql", &conf)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(conf.Dsn), &gorm.Config{
		Logger: logger.New(logrus.New(), logger.Config{
			SlowThreshold: time.Millisecond * 50,
			LogLevel:      logger.Info,
			Colorful:      true,
		}),
	})
	if err != nil {
		panic(err)
	}
	dao.InitTables(db)
	return db
}
