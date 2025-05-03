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
	List(ctx context.Context, userId int64, offset, limit int) ([]domain.Article, error)
	GetByArtId(ctx context.Context, id int64) (domain.Article, error)
	GetPubByArtId(ctx context.Context, id int64) (domain.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (u *articleService) GetPubByArtId(ctx context.Context, id int64) (domain.Article, error) {
	return u.repo.GetPubByArtId(ctx, id)
}

func (u *articleService) GetByArtId(ctx context.Context, id int64) (domain.Article, error) {
	return u.repo.GetByArtId(ctx, id)
}

func (j *articleService) List(ctx context.Context, userId int64, offset, limit int) ([]domain.Article, error) {
	return j.repo.List(ctx, userId, offset, limit)
}

func (a *articleService) WithDraw(ctx context.Context, article domain.Article) error {
	article.Status = domain.ArtStatusPrivate
	return a.repo.SyncStatus(ctx, article)
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
