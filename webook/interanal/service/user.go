package service

import (
	"context"

	"github.com/spigcoder/LittleBook/webook/interanal/domain"
)

type UserService struct {
}

func (s *UserService) SignUp(ctx context.Context, user domain.User) error {
	// TODO 加密并且存储起来
	return nil
}
