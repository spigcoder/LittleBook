package service

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/spigcoder/LittleBook/webook/internal/repository"
	"github.com/spigcoder/LittleBook/webook/internal/service/sms"
)

var tmpId string = ""

type CodeService struct {
	repo *repository.CodeRepository
	smsSvc  sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo: repo,
		smsSvc:  smsSvc,
	}
}

func (svc *CodeService) randomCode() string {
	key := rand.IntN(1000000)
	return fmt.Sprintf("%6d", key)
}

func (csv *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := csv.randomCode()
	err := csv.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	//这里证明设置Redis正常
	err = csv.smsSvc.Send(ctx, tmpId, []string{code}, phone)
	//这里如果出问题怎么办：能不能把Reids中的内容删除 -- 不可以，因为你不知道这里出的是什么问题
	return err
}

func (csv *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return csv.repo.Verify(ctx, biz, phone, inputCode)
}
