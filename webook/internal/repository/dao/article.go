package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"type:varchar(1024);not null"`
	Content  string `gorm:"type:BLOB;not null"`
	AuthorId int64  `gorm:"index:aid_ctime"`
	CTime    int64  `gorm:"index:aid_ctime"`
	UTime    int64
}

type PublishArticle struct {
	Article
}

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, article PublishArticle) (int64, error)
}

type GormArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{
		db: db,
	}
}

func (d *GormArticleDao) Upsert(ctx context.Context, article PublishArticle) (int64, error) {
	now := time.Now().UnixMilli()
	article.CTime = now
	article.UTime = now
	//设置如果存在更新哪些字段
	err := d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   article.Title,
			"content": article.Content,
			"u_time":  time.Now().UnixMilli(),
		})}).Create(&article).Error
	return article.Id, err
}

func (d *GormArticleDao) Sync(ctx context.Context, article Article) (int64, error) {
	id := article.Id
	err := d.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txDao := NewArticleDao(tx)
		if article.Id == 0 {
			id, err = txDao.Insert(ctx, article)
		} else {
			err = txDao.UpdateById(ctx, article)
		}
		if err != nil {
			return err
		}
		//操作线上库
		id, err = txDao.Upsert(ctx, PublishArticle{Article: article})
		return err
	})
	return id, err
}

func (d *GormArticleDao) UpdateById(ctx context.Context, article Article) error {
	arc := d.db.WithContext(ctx).Model(&Article{}).Where("id = ? AND author_id = ?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"u_time":  time.Now().UnixMilli(),
		})
	if arc.Error != nil {
		return arc.Error
	}
	if arc.RowsAffected == 0 {
		return fmt.Errorf("有人搞你：arcile_id:%d, author_id:%d", article.Id, article.AuthorId)
	}
	return nil
}

func (d *GormArticleDao) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CTime = now
	article.UTime = now
	err := d.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}
