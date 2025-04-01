package ioc

import (
	"github.com/spigcoder/LittleBook/webook/internal/service/sms"
	"github.com/spigcoder/LittleBook/webook/internal/service/sms/mem"
)

func InitSms() sms.Service {
	return mem.NewMemService()
}
