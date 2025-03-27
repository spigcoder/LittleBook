package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/spigcoder/LittleBook/webook/interanal/web"
)

type LoginJWTMiddlewareBuilder struct {
	path []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.path = append(l.path, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range l.path {
			if path == c.Request.URL.Path {
				return
			}
		}
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// seg := strings.Split(tokenHeader, " ")
		// if len(seg) != 2 {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }
		// tokenStr := seg[1]
		tokenStr := tokenHeader
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return web.ScretKey, nil
		})
		if err != nil || token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
