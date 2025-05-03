package article

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/spigcoder/LittleBook/webook/internal/events"
)

type Producer interface {
	Send(ctx context.Context, evt events.ReadEvent) error
}

type KafkaProducer struct {
	pro sarama.SyncProducer
}

func (p KafkaProducer) Send(ctx context.Context, evt events.ReadEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	//这里可以有一个简单的重试逻辑
	for i := 0; i < 3; i++ {
		_, _, err = p.pro.SendMessage(&sarama.ProducerMessage{
			Topic: events.ArticleRead,
			Value: sarama.ByteEncoder(data),
		})
		if err == nil {
			break
		}
	}
	return err
}

func NewKafkaProducer(pro sarama.SyncProducer) Producer {
	return &KafkaProducer{
		pro: pro,
	}
}
