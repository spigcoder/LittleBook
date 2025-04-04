package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound   = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)	
	FindByPhone(ctx context.Context, phone string) (User, error)
	Edit(ctx context.Context, u User) error
}

type GormUserDao struct {
	db *gorm.DB
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password string
	UserName string
	Birthday string
	Intro    string

	CreateTime int64
	UpdateTime int64
}

func NewUserDao(db *gorm.DB) *GormUserDao {
	return &GormUserDao{
		db: db,
	}
}

func (dao *GormUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

func (dao *GormUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *GormUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GormUserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	//跟底层强耦合，因为这里假设底层使用了MySQLj
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 { // Duplicate entry
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *GormUserDao) Edit(ctx context.Context, u User) error {
	u.UpdateTime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&u).Updates(u).Error
}
