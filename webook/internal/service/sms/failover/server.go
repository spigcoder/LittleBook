package failover

import (
	"context"
	"sync/atomic"

	"github.com/spigcoder/LittleBook/webook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	//当前server出错了多少次
	cnt uint32
	//当前使用的server下标
	idx uint32
	//阈值
	threshold uint32
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, threshold uint32) sms.Service {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx, cnt := atomic.LoadUint32(&t.idx), atomic.LoadUint32(&t.cnt)
	if cnt >= t.threshold {
		newIdx := (idx + 1) % uint32(len(t.svcs))
		if atomic.CompareAndSwapUint32(&t.idx, idx, newIdx) {
			//重置计数器
			atomic.StoreUint32(&t.cnt, 0)
		}
		idx = atomic.LoadUint32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case nil:
		atomic.StoreUint32(&t.cnt, 0)
		return nil	
	case context.DeadlineExceeded:
		atomic.AddUint32(&t.cnt, 1)
		return err
	default:
		//这里也可以做服务商的切换
		return err
	}
	return err
}
