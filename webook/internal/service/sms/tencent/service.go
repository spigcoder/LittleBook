package tencent

import (
	"context"
	"fmt"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(appId, signName string, client *sms.Client) *Service {
	return &Service{
		appId:    &appId,
		signName: &signName,
		client:   client,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = &tplId
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	resp, err := s.client.SendSmsWithContext(ctx, req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败, code: %s, msg: %s",
				*status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(args []string) []*string {
	res := make([]*string, 0, len(args))
	for _, arg := range args {
		res = append(res, &arg)
	}
	return res
}
