package service

import (
	"context"
	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	WithDraw(ctx context.Context, article domain.Article) error
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a *articleService) WithDraw(ctx context.Context, article domain.Article) error {
	article.Status = domain.ArtStatusPrivate
	return a.repo.SyncStatus(article)
}
func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArtStatusPublished
	return a.repo.Publish(ctx, article)
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArtStatusUnPublished
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	}
	return a.repo.Create(ctx, article)
}
