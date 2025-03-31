package repository

import (
	"context"
	"github.com/spigcoder/LittleBook/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository struct {
	codeCache *cache.CodeCache
}

func NewCodeRepository(codeCache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		codeCache: codeCache,
	}
}

func (c *CodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return c.codeCache.Set(ctx, biz, phone, code)
}

func (c *CodeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return c.codeCache.Verify(ctx, biz, phone, inputCode)
}