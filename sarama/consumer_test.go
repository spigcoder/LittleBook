package sarama

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	cg, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, "test_group", cfg)
	assert.NoError(t, err)
	err = cg.Consume(context.Background(), []string{"test_topic"}, TestHandler{})
	t.Log(err)
}

type TestHandler struct {
}

func (t TestHandler) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("SetUp")
	return nil
}

func (t TestHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Println("Cleanup")
	return nil
}

func (t TestHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msg := claim.Messages()
	for msg := range msg {
		fmt.Printf("Message claimed: %+v\n", string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
