package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Handler[T any] func(m *sarama.ConsumerMessage, t T) error

func NewHandler[T any](fn func(m *sarama.ConsumerMessage, t T) error) Handler[T] {
	return fn
}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msg := claim.Messages()
	var err error
	for msg := range msg {
		var t T
		err = json.Unmarshal(msg.Value, &t)
		if err != nil {
			logrus.Errorf("json序列化失败: topic: %s, message: %s, partition:%d, offset:%d, error: %t", msg.Topic, msg.Value, msg.Partition, msg.Offset, err)
			continue
		}
		for i := 0; i < 3; i++ {
			err = h(msg, t)
			if err == nil {
				break
			}
			logrus.Errorf("重试失败：topic: %s, message: %s, partition:%d, offset:%d, error: %t", msg.Topic, msg.Value, msg.Partition, msg.Offset, err)
		}
		if err != nil {
			logrus.Errorf("重试超过次数限制：topic: %s, message: %s, partition:%d, offset:%d, error: %t", msg.Topic, msg.Value, msg.Partition, msg.Offset, err)
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return err
}
