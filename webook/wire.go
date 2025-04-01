//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spigcoder/LittleBook/webook/internal/repository"
	"github.com/spigcoder/LittleBook/webook/internal/repository/cache"
	"github.com/spigcoder/LittleBook/webook/internal/repository/dao"
	"github.com/spigcoder/LittleBook/webook/internal/service"
	"github.com/spigcoder/LittleBook/webook/internal/web"
	"github.com/spigcoder/LittleBook/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		ioc.InitSms,
		service.NewUserService,
		service.NewCodeService,
		web.NewUserHandler,
		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
