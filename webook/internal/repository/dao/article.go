package dao

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"type:varchar(1024);not null"`
	Content  string `gorm:"type:BLOB;not null"`
	AuthorId int64  `gorm:"index:aid_ctime"`
	Status   uint8
	UTime    int64
	CTime    int64 `gorm:"index:aid_ctime"`
}

type PublishArticle struct {
	Article
}

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, article PublishArticle) (int64, error)
	SyncStatus(ctx context.Context, article Article) error
}

type GormArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{
		db: db,
	}
}

func (d *GormArticleDao) SyncStatus(ctx context.Context, article Article) error {
	now := time.Now().Unix()
	var err error
	d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? AND author_id = ?", article.Id, article.AuthorId).
			Updates(map[string]interface{}{
				"status": article.Status,
				"u_time": now,
			})
		if res.Error != nil {
			err = res.Error
			return res.Error
		}
		if res.RowsAffected == 0 {
			logrus.Errorf("有人在搞你或者id错误：art_id: %d, author_id: %d", article.Id, article.AuthorId)
			return fmt.Errorf("有人在搞你或者id错误：art_id: %d, author_id: %d", article.Id, article.AuthorId)
		}
		//线上库
		return tx.Model(&PublishArticle{}).Where("id = ?", article.Id).
			Updates(map[string]interface{}{
				"status": article.Status,
				"u_time": now,
			}).Error
	})
	return err
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
			"status":  article.Status,
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
			"status":  article.Status,
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
