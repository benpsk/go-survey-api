package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/benpsk/go-survey-api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GetUserById(ctx context.Context, id int) (*models.UserResponse, error) {
	return s.Repo.GetUserById(ctx, id)
}

func (s *Service) Create(ctx context.Context, user models.User) (*models.UserResponse, error) {
	password, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = password
	return s.Repo.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, user models.User) (*models.Token, error) {
	u, err := s.Repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("Password does not match")
	}
	token, err := generateJwt(u.Id)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func generateJwt(userId int) (*models.Token, error) {
	exp := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	return &models.Token{
		Type:        "Bearer",
		AccessToken: signed,
		ExpiredAt:   exp,
	}, nil
}
