package article

import (
	"context"
	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/repository/dao"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}

type CacheArticleReaderRepository struct {
	dao dao.ArticleDao
}

func NewCacheArticleReaderRepository(dao dao.ArticleDao) ArticleReaderRepository {
	return &CacheArticleReaderRepository{
		dao: dao,
	}
}

func (r *CacheArticleReaderRepository) Save(ctx context.Context, article domain.Article) (int64, error) {
	return 0, nil
}
