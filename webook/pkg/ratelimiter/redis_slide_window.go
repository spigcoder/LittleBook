package ratelimiter

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed slide_window.lua
var slideWinLimtiLua string

type RedisSlideWindowLimiter struct {
	cmd      redis.Cmdable
	//窗口大小
	interval time.Duration
	//阈值
	rate     int
}

func NewRedisSlideWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSlideWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (l *RedisSlideWindowLimiter) Limit(ctx context.Context, key string) (limited bool, err error) {
	return l.cmd.Eval(ctx, slideWinLimtiLua, []string{key},
		l.interval.Milliseconds(), l.rate, time.Now().UnixMilli()).Bool()
}
