package service

import (
	"context"
	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/repository"
)

type ArticleService interface {
	Edit(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a articleService) Edit(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		err := a.repo.UpdateById(ctx, article)
		return article.Id, err
	}
	return a.repo.Create(ctx, article)
}
