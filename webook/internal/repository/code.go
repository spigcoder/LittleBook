package repository

import (
	"context"

	"github.com/spigcoder/LittleBook/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type CacheCodeRepository struct {
	codeCache *cache.RedisCodeCache
}

func NewCodeRepository(codeCache *cache.RedisCodeCache) *CacheCodeRepository {
	return &CacheCodeRepository{
		codeCache: codeCache,
	}
}

func (c *CacheCodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return c.codeCache.Set(ctx, biz, phone, code)
}

func (c *CacheCodeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return c.codeCache.Verify(ctx, biz, phone, inputCode)
}
