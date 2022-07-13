package kafka

import (
	"context"
	"errors"
	"strings"

	"mail-sender/internal/mail"

	kafka "github.com/segmentio/kafka-go"
)

// Client - клиент очереди Kafka.
type ClientR struct {
	Reader *kafka.Reader
}

func New(brokers []string, topic string, groupId string) (*ClientR, error) {

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
func (c *ClientR) FetchProcessCommit(SMTPPort int) error {
	// Выборка очередного сообщения из Kafka.
	msg, err := c.Reader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	// Обработка сообщения
	err = messagesHandler(msg.Key, msg.Value, SMTPPort)
	if err != nil {
		return err
	}

	// Подтверждение сообщения как обработанного.
	err = c.Reader.CommitMessages(context.Background(), msg)
	return err
}

func messagesHandler(key, value []byte, SMTPPort int) error {

	messageV := strings.Split(string(value), " ")
	email := messageV[0]
	status := messageV[1]
	var body string

	if status == "ok" {
		body = "Your transaction has been successfully finished."
	} else {
		body = "Your transaction has not been finished."
	}

	return mail.SendMail(email, body)
}
