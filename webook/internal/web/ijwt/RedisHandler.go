package ijwt

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	AScretKey = []byte("ZD3oYULPnlBo2wqebduhFQjmrZdaFGaLzayCa8t8HWwxWKbRcGzaNLKkZ31ldeaM")
	RScretKey = []byte("od45zlxmSlBo2wqebduhFQjmrZdaFGaLzayCa8t8Hxfdsfwefwfdsvw3432s32sM")
)

type UserClaims struct {
	Uid int64 `json:"uid"`
	// UserAgent 是客户端的信息，用于校验
	UserAgent string
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Uid int64 `json:"uid"`
	jwt.RegisteredClaims
	//用于校验当前用户是否退出
	ssid string
}

func setRefreshJwt(c *gin.Context, uid int64) error {
	refreshClaims := RefreshClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)
	//签名
	tokenStr, err := token.SignedString([]byte(RScretKey))
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器问题")
		return err
	}
	c.Header("refresh-token", tokenStr)
	return nil

}

func setAccessJwt(c *gin.Context, uid int64) error {
	userClaims := UserClaims{
		Uid:       uid,
		UserAgent: c.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
	//签名
	tokenStr, err := token.SignedString([]byte(AScretKey))
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器问题")
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func SetJWT(c *gin.Context, uid int64) error {
	err := setAccessJwt(c, uid)
	err1 := setRefreshJwt(c, uid)
	return errors.Join(err, err1)
}
