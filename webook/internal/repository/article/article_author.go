package article

import (
	"context"
	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/repository"
	"github.com/spigcoder/LittleBook/webook/internal/repository/dao"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}

type CacheArticleAuthorRepository struct {
	dao dao.ArticleDao
}

func NewArtilceAuthorRepository(dao dao.ArticleDao) ArticleAuthorRepository {
	return &repository.CacheArticleRepository{
		dao: dao,
	}
}

func (r *CacheArticleAuthorRepository) Update(ctx context.Context, article domain.Article) error {
	return r.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}
func (r *CacheArticleAuthorRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return r.dao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}
