//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spigcoder/LittleBook/webook/internal/repository"
	"github.com/spigcoder/LittleBook/webook/internal/repository/cache"
	"github.com/spigcoder/LittleBook/webook/internal/repository/dao"
	"github.com/spigcoder/LittleBook/webook/internal/service"
	"github.com/spigcoder/LittleBook/webook/internal/web"
	"github.com/spigcoder/LittleBook/webook/ioc"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//最原始的依赖
		ioc.InitDB, ioc.InitRedis,
		// 数据库
		dao.NewUserDao,
		dao.NewArticleDao,
		// 缓存
		cache.NewUserCache,
		cache.NewCodeCache,
		// repository
		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewArtilceRepository,
		// 第三方服务
		ioc.InitSms,
		// service
		service.NewUserService,
		service.NewArticleService,
		service.NewCodeService,
		// handler
		web.NewArticleHandler,
		web.NewUserHandler,
		// 中间件
		ioc.InitMiddlewares,
		// web server
		ioc.InitGin,
	)
	return new(gin.Engine)
}

func InitTestDB() *gorm.DB {
	wire.Build(
		ioc.InitDB,
	)
	return new(gorm.DB)
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		ProviderSet,
		dao.NewArticleDao,
		service.NewArticleService,
		web.NewArticleHandler,
		repository.NewArtilceRepository,
	)
	return new(web.ArticleHandler)
}
