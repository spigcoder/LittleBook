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
	SyncStatus(ctx context.Context, article domain.Article) error
	List(ctx context.Context, id int64, offset int, limit int) ([]domain.Article, error)
	GetByArtId(ctx context.Context, id int64) (domain.Article, error)
	GetPubByArtId(ctx context.Context, id int64) (domain.Article, error)
}

type CacheArticleRepository struct {
	dao      dao.ArticleDao
	UserRepo UserRepository
}

func NewArtilceRepository(dao dao.ArticleDao, repository UserRepository) ArticleRepository {
	return &CacheArticleRepository{
		dao:      dao,
		UserRepo: repository,
	}
}

func (r *CacheArticleRepository) GetPubByArtId(ctx context.Context, id int64) (domain.Article, error) {
	art, err := r.dao.GetByArtId(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	aut, err := r.UserRepo.FindById(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	res := convertDaoToDomain(art)
	res.Author.Name = aut.UserName
	return res, nil
}
func (r *CacheArticleRepository) GetByArtId(ctx context.Context, id int64) (domain.Article, error) {
	res, err := r.dao.GetByArtId(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return convertDaoToDomain(res), nil
}

func (r *CacheArticleRepository) List(ctx context.Context, id int64, offset int, limit int) ([]domain.Article, error) {
	articles, err := r.dao.List(ctx, id, offset, limit)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Article, len(articles))
	for i, article := range articles {
		result[i] = convertDaoToDomain(article)
	}
	go func() {

	}()
	return result, nil
}

func (r *CacheArticleRepository) SyncStatus(ctx context.Context, article domain.Article) error {
	return r.dao.SyncStatus(ctx, convertDomainToDao(article))
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

func convertDomainToDao(article domain.Article) dao.Article {
	return dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Id:       article.Id,
		Status:   article.Status.ToUInt8(),
	}
}

func convertDaoToDomain(article dao.Article) domain.Article {
	return domain.Article{
		Title:   article.Title,
		Content: article.Content,
		Author: domain.Author{
			Id: article.AuthorId,
		},
		Id:     article.Id,
		Status: domain.ArtStatus(article.Status),
		CTime:  article.CTime,
		UTime:  article.UTime,
	}
}
