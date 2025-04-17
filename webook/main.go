package main

import (
	"github.com/spigcoder/LittleBook/webook/ioc"
	"github.com/spigcoder/LittleBook/webook/startup"
)

func main() {
	ioc.InitViper()
	ioc.InitLogrus()
	server := startup.InitWebServer()
	server.Run(":8080")
}
