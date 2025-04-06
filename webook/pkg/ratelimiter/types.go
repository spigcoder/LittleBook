package ratelimiter

import (
	"context"
)

type Limiter interface {
	// Limit 是否达到限流阈值
	Limit(ctx context.Context, key string) (limited bool, err error)
}
