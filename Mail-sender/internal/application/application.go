package application

import (
	"context"
	"log"
	"net/http"
	"os"

	"mail-sender/config"
	"mail-sender/internal/kafka"
	httpServer "mail-sender/internal/server"
)

var (
	s      *http.Server
	logger *log.Logger
)

func Start(ctx context.Context) {

	conf := config.New("config", ".\\config\\")
	logger := log.New(os.Stderr, "", log.Lshortfile)

	go httpServer.New(conf, logger)

	kafkaClient, err := kafka.New(
		&kafka.ToKafkaStartFunc{
			KafkaAdress: []string{"localhost:" + conf.GetKafkaPort()},
			Topic:       "transaction",
			GroupID:     "G1",
		})

	if err != nil {
		logger.Fatal(err)
	}

	defer kafkaClient.Reader.Close()

	for {
		err := kafkaClient.FetchProcessCommit(conf, logger)

		if err != nil {
			logger.Printf("Kafka runtime error. Error: %v", err)
		}
	}
}

func Stop() {
	_ = s.Shutdown(context.Background())
	logger.Printf("The application has been stopped")
}
