package services

import "github.com/benpsk/go-survey-api/internal/repositories"

type Service struct {
	Repo *repositories.Repository
}

func New(repo *repositories.Repository) *Service {
	return &Service{Repo: repo}
}
