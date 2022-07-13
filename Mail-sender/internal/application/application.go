package application

import (
	"context"
	"log"
	"net/http"
	"os"

	"mail-sender/config"
	kafka "mail-sender/internal/kafka"
)

var (
	s *http.Server
)

func Start(ctx context.Context) {

	conf := config.New()
	logger := log.New(os.Stderr, "", log.Lshortfile)

	kafkaClient, err := kafka.New([]string{"localhost:9092"}, "transaction", "G1")
	if err != nil {
		logger.Fatal(err)
	}

	defer kafkaClient.Reader.Close()

	for {
		smtpPort := conf.SMTPPort
		kafkaClient.FetchProcessCommit(smtpPort)
	}
}

func Stop() {
	_ = s.Shutdown(context.Background())
	log.Printf("The application has been stopped")
}
