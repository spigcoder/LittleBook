package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"time"
)

type BatchHandler[T any] struct {
	fn        func(msg []*sarama.ConsumerMessage, event []T) error
	batchSize int
	timeout   time.Duration
}

func NewBatchHandler[T any](fn func(msg []*sarama.ConsumerMessage, event []T) error, batchSize int, timeout time.Duration) BatchHandler[T] {
	return BatchHandler[T]{
		fn:        fn,
		batchSize: batchSize,
		timeout:   timeout,
	}
}

func (b BatchHandler[T]) Setup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (b BatchHandler[T]) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (b BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgCh := claim.Messages()
	batchSize := b.batchSize
	for {
		ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
	Loop:
		for i := 0; i < batchSize; i++ {
			select {
			case <-ctx.Done():
				break Loop
			case msg := <-msgCh:
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					logrus.Errorf("json序列化失败: topic: %s, message: %s, partition:%d, offset:%d, error: %t", msg.Topic, msg.Value, msg.Partition, msg.Offset, err)
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		err := b.fn(msgs, ts)
		if err != nil {
			logrus.Error(err)
		}
		for _, msg := range msgs {
			session.MarkMessage(msg, "")
		}
	}
}
