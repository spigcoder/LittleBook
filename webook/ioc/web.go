package ioc

import (
	"github.com/spigcoder/LittleBook/webook/pkg/ginx/middleware/logs"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spigcoder/LittleBook/webook/internal/web"
	"github.com/spigcoder/LittleBook/webook/internal/web/middleware"
	"github.com/spigcoder/LittleBook/webook/pkg/ginx/middleware/ratelimit"
	"github.com/spigcoder/LittleBook/webook/pkg/ratelimiter"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	mds := []gin.HandlerFunc{
		cors.New(cors.Config{
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
		}),
		logs.NewMiddlewareBuilder().
			EnableRequest().
			EnableResponse().Build(),
		ratelimit.NewBuilder(ratelimiter.NewRedisSlideWindowLimiter(redisClient, time.Second, 100)).Build(),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/signup").IgnorePaths("/users/login").
			IgnorePaths("/users/code/send").IgnorePaths("/users/login_sms").
			Build(),
	}
	return mds
}
