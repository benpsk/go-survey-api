package survey

import (
	"context"
)

type SurveyService struct {
	Repo *SurveyRepository
}

func NewSurveyService(repo *SurveyRepository) *SurveyService {
	return &SurveyService{Repo: repo}
}

func (s *SurveyService) Store(ctx context.Context, survey SurveyInput) (*SurveyInput, error) {
	res, err := s.Repo.Store(ctx, survey)
	if err != nil {
		return nil, err
	}
	return s.formatSurvey(*res)
}

func (s *SurveyService) Get(ctx context.Context, userId int) ([]SurveyInput, error) {
	surveys, err := s.Repo.getByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	var result []SurveyInput
	for _, res := range surveys {
		survey := SurveyInput{
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

func (s *SurveyService) GetById(ctx context.Context, id int) (*SurveyInput, error) {
	res, err := s.Repo.getById(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.formatSurvey(*res)
}

func (s *SurveyService) formatSurvey(res Survey) (*SurveyInput, error) {
	return &SurveyInput{
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

func (s *SurveyService) Update(ctx context.Context, id int, survey SurveyInput) (*SurveyInput, error) {
	res, err := s.Repo.Update(ctx, id, survey)
	if err != nil {
		return nil, err
	}
	return s.formatSurvey(*res)
}
