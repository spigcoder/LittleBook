package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spigcoder/LittleBook/webook/internal/service/sms"
)

type Service struct {
	svc sms.Service
	key string
}

func NewService(svc sms.Service, key string) sms.Service {
	return &Service{
		svc: svc,
		key: key,
	}
}

// 这里是为了安全性，我们要保证调用我们这个服务的调用方是我们允许的，有完整的申请和审批流程
func (s *Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	var claims Claims
	token, err := jwt.ParseWithClaims(biz, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token 不合法")
	}
	//下面就调用我们的短信服务
	return s.svc.Send(ctx, claims.tplId, args, numbers...)

}

type Claims struct {
	jwt.RegisteredClaims
	tplId string
}
