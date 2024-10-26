package repositories

import "github.com/jackc/pgx/v5"

type Repository struct {
	Conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repository {
	return &Repository{
		Conn: conn,
	}
}
