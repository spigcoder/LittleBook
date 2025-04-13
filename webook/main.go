package main

import (
	"github.com/spf13/viper"
)

func main() {
	initViper()
	server := InitWebServer()
	server.Run(":8080")
}

func initLog() {

}

func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
