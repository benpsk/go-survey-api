package user

import (
	"context"
	"errors"
	"time"

	"github.com/benpsk/go-survey-api/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetUserById(ctx context.Context, id int) (*UserResponse, error) {
	return s.Repo.GetUserById(ctx, id)
}

func (s *UserService) Create(ctx context.Context, user User) (*UserResponse, error) {
	password, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = password
	return s.Repo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, user User) (*Token, error) {
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

func generateJwt(userId int) (*Token, error) {
	exp := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return nil, err
	}
	return &Token{
		Type:        "Bearer",
		AccessToken: signed,
		ExpiredAt:   exp,
	}, nil
}
