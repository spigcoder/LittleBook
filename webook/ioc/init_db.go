package ioc

import (
	"github.com/spf13/viper"
	"github.com/spigcoder/LittleBook/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	db, err := gorm.Open(mysql.Open(conf.Dsn))
	if err != nil {
		panic(err)
	}
	dao.InitTables(db)
	return db
}
