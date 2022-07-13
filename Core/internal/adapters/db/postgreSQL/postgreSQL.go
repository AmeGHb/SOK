package postgresql

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"transaction/config"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, pConf config.DatabaseConfig) (pool *pgxpool.Pool, err error) {

	connection := pConf.Url

	err = doWithAttempts(func() error {

		ctx, cancel := context.WithTimeout(
			ctx, time.Duration(pConf.Seconds)*time.Second,
		)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, connection)
		if err != nil {
			return err
		}
		return nil
	}, pConf)

	if err != nil {
		log.Fatalf("Failed to connect to the database. Error: %v", err)
	}

	return
}

func doWithAttempts(function func() error, pConf config.DatabaseConfig) (err error) {

	var attemptNumber int = 1

	for pConf.Attempts > 0 {

		if err = function(); err != nil {
			log.Printf("Failed to connnect to the database. Attempt = %d", attemptNumber)
			time.Sleep(pConf.Delay)

			pConf.Attempts--
			attemptNumber++

			continue
		}
		return nil
	}
	return
}
