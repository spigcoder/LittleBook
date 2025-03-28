package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spigcoder/LittleBook/webook/interanal/repository"
	"github.com/spigcoder/LittleBook/webook/interanal/repository/dao"
	"github.com/spigcoder/LittleBook/webook/interanal/service"
	"github.com/spigcoder/LittleBook/webook/interanal/web"
	"github.com/spigcoder/LittleBook/webook/interanal/web/middleware"
	"github.com/spigcoder/LittleBook/webook/pkg/ginx/middleware/ratelimit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	user := initUser(db)

	server := initWebServer()
	user.RegisterRoutes(server)
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		//这里用来配置允许的域名
		// AllowOrigins:     []string{"https://foo.com"},
		//如果没有这个，那就是默认所有的都可以
		// AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//只有加上这个，前端才能拿到这个header
		ExposeHeaders: []string{"x-jwt-token"},
		//允许带cookie之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	//启用session
	// store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("oez06bfpafdk77ocfcpc0tyrc5izmr9r"), []byte("tisjdqf9omlwdztf6codcmeslh352bpv"))
	// if err != nil {
	// 	panic(err)
	// }
	// server.Use(sessions.Sessions("webook", store))
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	userDao := dao.NewUserDao(db)
	userRepo := repository.NewUserRepository(userDao)
	userSvc := service.NewUserService(userRepo)
	user := web.NewUserHandler(userSvc)
	return user
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	dao.InitTables(db)
	return db
}
