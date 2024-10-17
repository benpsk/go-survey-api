package internal

import (
	"net/http"

	"github.com/benpsk/go-survey-api/internal/middleware"
)

func Router(c *Container) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", c.UserHandler.Register)
	mux.HandleFunc("POST /login", c.UserHandler.Login)

	mux.HandleFunc("/user", middleware.Auth(c.UserHandler.User))
	mux.HandleFunc("POST /survey", middleware.Auth(c.SurveyHandler.Store))
	mux.HandleFunc("/surveys", middleware.Auth(c.SurveyHandler.Get))
	mux.HandleFunc("/surveys/{id}", middleware.Auth(c.SurveyHandler.GetById))
	mux.HandleFunc("PUT /surveys/{id}", middleware.Auth(c.SurveyHandler.Update))
	return mux
}
