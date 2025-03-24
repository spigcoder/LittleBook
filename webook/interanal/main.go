package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spigcoder/LittleBook/webook/interanal/web"
)

func main() {
	server := gin.Default()
	user := web.NewUserHandler()
	user.RegisterRoutes(server)
	server.Run(":8080")
}
