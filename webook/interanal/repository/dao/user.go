package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound   = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	UserName string
	Birthday string
	Intro    string

	CreateTime int64
	UpdateTime int64
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
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

func (dao *UserDao) Edit(ctx context.Context, u User) error {
	u.UpdateTime = time.Now().UnixMilli()
	fmt.Println(u)
	return dao.db.WithContext(ctx).Model(&u).Updates(u).Error
}
