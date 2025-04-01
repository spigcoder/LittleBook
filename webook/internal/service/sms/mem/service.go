package mem

import (
	"context"
	"fmt"
)

type MemService struct {
}

func (c *MemService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
