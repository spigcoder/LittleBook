package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
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

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

type GormArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{
		db: db,
	}
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
