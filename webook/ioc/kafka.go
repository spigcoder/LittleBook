package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/spigcoder/LittleBook/webook/internal/events"
	"github.com/spigcoder/LittleBook/webook/internal/events/article"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addr []string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	SaramaCfg := sarama.NewConfig()
	SaramaCfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, SaramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	pro, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return pro
}

func NewConsumer(conumer *article.KafkaConsumer) []events.Consumer {
	return []events.Consumer{conumer}
}
