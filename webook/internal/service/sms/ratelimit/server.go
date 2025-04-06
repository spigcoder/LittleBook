package ratelimit

import (
	"context"
	"fmt"

	"github.com/spigcoder/LittleBook/webook/internal/service/sms"
	"github.com/spigcoder/LittleBook/webook/pkg/ratelimiter"
)

var limitErr error = fmt.Errorf("触发限流")

type RateLimitServer struct {
	limiter ratelimiter.Limiter
	svc     sms.Service
}

func NewRateLimitServer(limiter ratelimiter.Limiter, svc sms.Service) sms.Service {
	return &RateLimitServer{
		limiter: limiter,
		svc:     svc,
	}
}

//这就是一个装饰器模式
func (r *RateLimitServer)Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limit, err :=  r.limiter.Limit(ctx, "sms_code")
	if err != nil  {
		return fmt.Errorf("限流服务异常:%w", err)	
	}
	if limit {
		return limitErr
	}
	//在前面添加一些功能
	err = r.svc.Send(ctx, tplId, args, numbers...)
	//在后面添加一些功能
	return err	
}
