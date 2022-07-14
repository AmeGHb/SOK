package mail_test

import (
	"log"
	"os"
	"testing"

	"mail-sender/config"
	"mail-sender/internal/mail"
	httpServer "mail-sender/internal/server"
)

func TestMail(t *testing.T) {

	TestConfig := config.Config{
		HttpPort:  4000,
		KafkaPort: 9092,
		MailPort:  4000,
		MailHost:  "localhost",
		MailFrom:  "NeveR2MorE@yandex.ru",
	}
	logger := log.New(os.Stderr, "", log.Lshortfile)

	go httpServer.New(&TestConfig, logger)

	testEmail := "NeveR5MorE@yandex.ru"
	testBody := "Your transaction has been successfully finished."

	err := mail.SendMail(testEmail, testBody, &TestConfig)
	if err != nil {
		t.Fatalf("/mail.SendMail()/ Mail send failed: \n%v", err)
	}
}
