package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func Connect(dsn string) (*pgx.Conn, error) {
  if db != nil {
    return db, nil
  }
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
  db = conn
	return db, err
}
