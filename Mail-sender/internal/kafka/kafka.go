package kafka

import (
	"context"
	"errors"
	"strings"

	"mail-sender/config"
	"mail-sender/internal/mail"

	kafka "github.com/segmentio/kafka-go"
)

// Client - клиент очереди Kafka.
type ClientR struct {
	Reader *kafka.Reader
}

func New(data *ToKafkaStartFunc) (*ClientR, error) {

	brokers := data.KafkaAdress
	topic := data.Topic
	groupId := data.GroupID

	if len(brokers) == 0 || brokers[0] == "" || topic == "" || groupId == "" {
		return nil, errors.New("not given all parameters for kafka client")
	}

	client := ClientR{}

	client.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupId,
		MinBytes: 10e0,
		MaxBytes: 10e6,
	})

	return &client, nil
}

func (c *ClientR) getMessage() (kafka.Message, error) {
	msg, err := c.Reader.ReadMessage(context.Background())
	return msg, err
}

// FetchProcessCommit сначала выбирает сообщение из очереди,
// потом обрабатывает, после чего подтверждает.
func (c *ClientR) FetchProcessCommit(conf *config.Config) error {
	// Выборка очередного сообщения из Kafka.
	msg, err := c.Reader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	// Обработка сообщения
	err = messagesHandler(msg.Value, conf)
	if err != nil {
		return err
	}

	// Подтверждение сообщения как обработанного.
	err = c.Reader.CommitMessages(context.Background(), msg)
	return err
}

func messagesHandler(value []byte, conf *config.Config) error {

	messageV := strings.Split(string(value), " ")
	email, status := messageV[0], messageV[1]
	var body string

	if status == "ok" {
		body = "Your transaction has been successfully finished."
	} else {
		body = "Your transaction has not been finished."
	}

	return mail.SendMail(email, body, conf)
}
