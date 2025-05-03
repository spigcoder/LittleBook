package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9094"}, cfg)
	assert.NoError(t, err)
	p, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("hello kafka"),
		//生产者和消费者之间进行数据传递
		Headers: []sarama.RecordHeader{
			{Key: []byte("key"), Value: []byte("value")},
		},
		Metadata: map[string]any{"MetaData": "MetaDataValue"},
	})
	assert.NoError(t, err)
	t.Log(p, offset)
}

func TestPRoducerAsync(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	producer, _ := sarama.NewAsyncProducer([]string{"127.0.0.1:9094"}, cfg)
	msgChan := producer.Input()
	msgChan <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Key:   sarama.StringEncoder("key"),
		Value: sarama.StringEncoder("hello 111kafka"),
		//生产者和消费者之间进行数据传递
		Headers: []sarama.RecordHeader{
			{Key: []byte("key"), Value: []byte("value")},
		},
		Metadata: map[string]any{"MetaData": "MetaDataValue"},
	}
	select {
	case msg := <-producer.Successes():
		val, _ := msg.Value.Encode()
		t.Log(string(val))
	case err := <-producer.Errors():
		val, _ := err.Msg.Value.Encode()
		t.Log(string(val), err.Err)
	}
}
