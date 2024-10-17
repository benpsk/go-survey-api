package internal

import (
	"github.com/benpsk/go-survey-api/internal/survey"
	"github.com/benpsk/go-survey-api/internal/user"
	"github.com/jackc/pgx/v5"
)

type Container struct {
	Db            *pgx.Conn
	UserHandler   *user.UserHandler
	SurveyHandler *survey.SurveyHandler
}

func NewContainer(db *pgx.Conn) *Container {
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	surveyRepo := survey.NewSurveyRepository(db)
	surveyService := survey.NewSurveyService(surveyRepo)
	surveyHandler := survey.NewSurveyHandler(surveyService)

	return &Container{
		Db:            db,
		UserHandler:   userHandler,
		SurveyHandler: surveyHandler,
	}
}
