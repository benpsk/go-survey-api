package repositories

import (
	"context"
	"errors"

	"github.com/benpsk/go-survey-api/internal/models"
)

func (repo *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email, password FROM users WHERE email=$1"

	err := repo.Conn.QueryRow(ctx, query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return &user, nil
}

func (repo *Repository) GetUserById(ctx context.Context, id int) (*models.UserResponse, error) {
	var user models.UserResponse
	query := "SELECT id, name, email FROM users WHERE id=$1"

	err := repo.Conn.QueryRow(ctx, query, id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *Repository) Create(ctx context.Context, user models.User) (*models.UserResponse, error) {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Conn.QueryRow(ctx, query, user.Name, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return repo.GetUserById(ctx, user.Id)
}
