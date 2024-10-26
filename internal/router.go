package internal

import (
	"net/http"

	"github.com/benpsk/go-survey-api/internal/handlers"
	"github.com/benpsk/go-survey-api/internal/middlewares"
)

func Router(handler *handlers.Handler) *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", handler.Register)
	mux.HandleFunc("POST /login", handler.Login)

	mux.HandleFunc("/user", middlewares.Auth(handler.User))
	mux.HandleFunc("POST /survey", middlewares.Auth(handler.Store))
	mux.HandleFunc("/surveys", middlewares.Auth(handler.Get))
	mux.HandleFunc("/surveys/{id}", middlewares.Auth(handler.GetById))
	mux.HandleFunc("PUT /surveys/{id}", middlewares.Auth(handler.Update))
	return mux
}
