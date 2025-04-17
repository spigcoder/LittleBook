package ioc

import (
	"github.com/spf13/viper"
)

func InitViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/Users/x/go/src/LittleBook/webook/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
