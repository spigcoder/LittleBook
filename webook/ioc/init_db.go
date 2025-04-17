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

// 自定义彩色 Writer
type colorWriter struct {
	logrusLogger *logrus.Logger
}

func InitDB() *gorm.DB {
	type config struct {
		Dsn string
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true, // 强制颜色输出
		FullTimestamp: true,
	})
	var conf config
	err := viper.UnmarshalKey("db.mysql", &conf)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(conf.Dsn), &gorm.Config{
		Logger: logger.New(
			logrus.New(),
			logger.Config{
				SlowThreshold: time.Millisecond * 50,
				LogLevel:      logger.Info, // 只显示警告和错误
				Colorful:      true,
			},
		),
	})
	if err != nil {
		panic(err)
	}
	dao.InitTables(db)
	return db
}
