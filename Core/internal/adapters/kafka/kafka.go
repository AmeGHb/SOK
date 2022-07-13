package kafka

import (
	"context"
	"errors"

	kafka "github.com/segmentio/kafka-go"
)

// Client - клиент очереди Kafka.
type ClientW struct {
	Writer *kafka.Writer
}

func New(brokers []string, topic string, groupId string) (*ClientW, error) {

	if len(brokers) == 0 || brokers[0] == "" || topic == "" || groupId == "" {
		return nil, errors.New("not given all parameters for kafka client")
	}

	client := ClientW{}

	client.Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers[0]),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &client, nil
}

// SendMessages отправляет сообщения в Kafka.
func (c *ClientW) SendMessages(messages []kafka.Message) error {
	err := c.Writer.WriteMessages(context.Background(), messages...)
	return err
}
