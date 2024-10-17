package user

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	Conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{Conn: conn}
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := "SELECT id, name, email, password FROM users WHERE email=$1"

	err := repo.Conn.QueryRow(ctx, query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetUserById(ctx context.Context, id int) (*UserResponse, error) {
	var user UserResponse
	query := "SELECT id, name, email FROM users WHERE id=$1"

	err := repo.Conn.QueryRow(ctx, query, id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) Create(ctx context.Context, user User) (*UserResponse, error) {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Conn.QueryRow(ctx, query, user.Name, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return repo.GetUserById(ctx, user.Id)
}
