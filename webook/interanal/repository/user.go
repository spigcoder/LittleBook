package repository

import (
	"context"

	"github.com/spigcoder/LittleBook/webook/interanal/domain"
	"github.com/spigcoder/LittleBook/webook/interanal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) FindByEmail(c context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(c, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		UserName: u.UserName,
		Password: u.Password,
	}, nil
}

func (repo *UserRepository) Create(c context.Context, u domain.User) error {
	err := repo.dao.Insert(c, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

	return err
}

func (repo *UserRepository) Edit(c context.Context, u domain.User) error {
	err := repo.dao.Edit(c, dao.User{
		Birthday: u.Birthday,
		Intro:    u.Intro,
		UserName: u.UserName,
		Id:       u.Id,
	})
	return err
}

func (u *UserRepository) GetUserById(id int) (user domain.User, err error) {
	// TODO: implement
	return
}
