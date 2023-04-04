package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectPostgres postgre
func ConnectPostgres(source string) (*pgxpool.Pool, error) {
	log.Println(source)
	var err error
	client, err := pgxpool.New(context.Background(), source)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background()); err != nil {
		return nil, err
	}

	log.Println("Successfully connected with postgres db!")

	return client, nil
}
