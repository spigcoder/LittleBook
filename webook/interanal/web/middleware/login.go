package middleware

import (
	"net/http"

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
	}
}
