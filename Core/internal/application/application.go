package application

import (
	"context"
	"log"
	"os"

	"transaction/config"
	"transaction/internal/adapters/db/postgresql"
	"transaction/internal/adapters/http"
	"transaction/internal/adapters/kafka"
	"transaction/internal/users/db"

	"golang.org/x/sync/errgroup"
)

var (
	s *http.Server
)

func Start(ctx context.Context) {

	conf := config.New()
	logger := log.New(os.Stderr, "", log.Lshortfile)

	postgreSQLClient, err := postgresql.NewClient(ctx, *conf.DatabaseConfig)
	if err != nil {
		log.Fatalf("%v", err)
	}

	defer func() {
		postgreSQLClient.Close()
	}()

	userRepository := db.NewRepository(postgreSQLClient, logger)

	kafkaClient, err := kafka.New([]string{"localhost:9092"}, "transaction", "G1")
	if err != nil {
		log.Fatalf("Kafka client error: %v", err)
	}

	s, err := http.New(ctx, conf, userRepository, kafkaClient)
	if err != nil {
		log.Fatalln("http server creating failed")
	}

	var g errgroup.Group
	g.Go(func() error {
		return s.Start()
	})

	err = g.Wait()
	if err != nil {
		log.Printf("http server start failed")
	}
}

func Stop() {
	_ = s.Stop(context.Background())
	log.Printf("The application has been stopped")
}
