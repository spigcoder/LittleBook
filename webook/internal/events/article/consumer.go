package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spigcoder/LittleBook/webook/internal/events"
	"github.com/spigcoder/LittleBook/webook/internal/repository"
	"github.com/spigcoder/LittleBook/webook/pkg/saramax"
	"time"
)

type KafkaConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
}

func NewKafkaConsumer(client sarama.Client, repo repository.InteractiveRepository) *KafkaConsumer {
	return &KafkaConsumer{
		client: client,
		repo:   repo,
	}
}

func (c *KafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, "Interactive", sarama.NewConfig())
	if err != nil {
		return err
	}
	go func() {
		for {
			err := cg.Consume(context.Background(), []string{events.ArticleRead}, saramax.NewHandler[events.ReadEvent](c.Consume))
			if err != nil {
				logrus.Error("消费消息失败：", err)
			}
		}
	}()
	return nil
}

func (c *KafkaConsumer) Consume(msg *sarama.ConsumerMessage, event events.ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.repo.IncrReadCnt(ctx, "article", event.Aid)
}
