package ratelimit

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spigcoder/LittleBook/webook/pkg/ratelimiter"
)

type Builder struct {
	prefix  string
	limiter ratelimiter.Limiter
}

func NewBuilder(limiter ratelimiter.Limiter) *Builder {
	return &Builder{
		prefix:  "ip-limiter",
		limiter: limiter,
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			log.Println(err)
			// 这一步很有意思，就是如果这边出错了
			// 要怎么办？
			// 可以限流，也可以不限流
			// 限流：这种情况是你的下游比较的坑，你不能相信他，这是一种保守的策略
			// 不限流：这有两种可能，一种情况是你的下游涉及的比较好，你可以相信他
			//   				 另一种情况是你的服务的高可用性比较强，你要保证它的可用性，即使Redis挂了，也要发送数据
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.limiter.Limit(ctx, key)
}
