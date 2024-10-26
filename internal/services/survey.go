package services

import (
	"context"

	"github.com/benpsk/go-survey-api/internal/models"
)

func (s *Service) Store(ctx context.Context, survey models.SurveyInput) (*models.SurveyInput, error) {
	res, err := s.Repo.Store(ctx, survey)
	if err != nil {
		return nil, err
	}
	return s.formatSurvey(*res)
}

func (s *Service) Get(ctx context.Context, userId int) ([]models.SurveyInput, error) {
	surveys, err := s.Repo.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	var result []models.SurveyInput
	for _, res := range surveys {
		survey := models.SurveyInput{
			Id:        res.Id,
			UserId:    res.UserId,
			Name:      res.Name,
			PhoneNo:   res.PhoneNo,
			Gender:    res.Gender,
			Dob:       res.Dob.Format("2006-01-02"),
			CreatedAt: res.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: res.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, survey)
	}
	return result, nil
}

func (s *Service) GetById(ctx context.Context, id int) (*models.SurveyInput, error) {
	res, err := s.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.formatSurvey(*res)
}

func (s *Service) formatSurvey(res models.Survey) (*models.SurveyInput, error) {
	return &models.SurveyInput{
		Id:        res.Id,
		UserId:    res.UserId,
		Name:      res.Name,
		PhoneNo:   res.PhoneNo,
		Gender:    res.Gender,
		Dob:       res.Dob.Format("2006-01-02"),
		CreatedAt: res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: res.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *Service) Update(ctx context.Context, id int, survey models.SurveyInput) (*models.SurveyInput, error) {
	res, err := s.Repo.Update(ctx, id, survey)
	if err != nil {
		return nil, err
	}
	return s.formatSurvey(*res)
}
