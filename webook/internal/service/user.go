package service

import (
	"context"
	"errors"

	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
)

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
	Login(ctx context.Context, user domain.User) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Edit(ctx context.Context, user domain.User) error
}

type UserServiceX struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceX{repo: repo}
}

func (s *UserServiceX) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := s.repo.FindByPhone(ctx, phone)
	if err == repository.ErrUserNotFound {
		// 创建一个用户
		user := domain.User{
			Phone: phone,
		}
		err = s.repo.Create(ctx, user)
		if err != nil {
			return domain.User{}, err
		}
		return user, nil
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, err
}

func (s *UserServiceX) SignUp(ctx context.Context, user domain.User) error {
	HashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(HashPass)

	err = s.repo.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceX) Login(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := s.repo.FindByEmail(ctx, user.Email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (s *UserServiceX) Edit(ctx context.Context, user domain.User) error {
	//用不用校验用户不存在的问题，我认为不用，因为如果用户不存在，那么就不会调用这个方法
	return s.repo.Edit(ctx, user)
}
