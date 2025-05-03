package repository

import "context"

type InteractiveRepository interface {
	Like(ctx context.Context, uid, aid int64) error
	DisLike(ctx context.Context, uid, aid int64) error
	IsLike(ctx context.Context, uid, aid int64) (bool, error)
	IncrReadCnt(ctx context.Context, biz string, aid int64) error
}
