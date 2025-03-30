package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	path []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.path = append(l.path, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range l.path {
			if path == c.Request.URL.Path {
				return
			}
		}
		see := sessions.Default(c)
		userId := see.Get("userId")
		if userId == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//更新Session过期时间
		updateTime := see.Get("update_time")
		now := time.Now().UnixMilli()
		see.Set("userId", userId)
		see.Options(sessions.Options{
			MaxAge: 60*30,
		})
		if updateTime == nil {
			see.Set("update_time", now)
			see.Save()
			return
		}
		if now-updateTime.(int64) > 60*1000*15 {
			see.Set("update_time", now)
			see.Save()
			return
		}
	}
}
