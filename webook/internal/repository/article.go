package repository

import (
	"context"
	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Publish(ctx context.Context, article domain.Article) (int64, error)
}

type CacheArticleRepository struct {
	dao dao.ArticleDao
}

func NewArtilceRepository(dao dao.ArticleDao) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}

func convertDomainToDao(article domain.Article) dao.Article {
	return dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Id:       article.Id,
	}
}

func (r *CacheArticleRepository) Publish(ctx context.Context, article domain.Article) (int64, error) {
	return r.dao.Sync(ctx, convertDomainToDao(article))
}

func (r *CacheArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return r.dao.UpdateById(ctx, convertDomainToDao(article))
}
func (r *CacheArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return r.dao.Insert(ctx, convertDomainToDao(article))
}
