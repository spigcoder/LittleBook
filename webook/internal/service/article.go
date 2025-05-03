package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/events"
	"github.com/spigcoder/LittleBook/webook/internal/events/article"
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
	pro  article.Producer
}

func NewArticleService(repo repository.ArticleRepository, pro article.Producer) ArticleService {
	return &articleService{
		repo: repo,
		pro:  pro,
	}
}

func (u *articleService) GetPubByArtId(ctx context.Context, id int64) (domain.Article, error) {
	art, err := u.repo.GetPubByArtId(ctx, id)
	if err == nil {
		//	发送阅读时间，后端得到后可以异步的进行阅读数的增加
		go func() {
			err := u.pro.Send(ctx, events.ReadEvent{
				//这里如果消费者要使用art，usr的数据，让他自己去数据库查询
				Uid: art.Author.Id,
				Aid: art.Id,
			})
			if err != nil {
				logrus.Error("发送消息失败：", err)
			}
		}()
	}
	return art, err
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
