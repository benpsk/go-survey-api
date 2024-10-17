package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func Connect(dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	return conn, err
}
