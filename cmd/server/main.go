package main

import (
	"context"
	"log"
	"net/http"

	"github.com/benpsk/go-survey-api/config"
	"github.com/benpsk/go-survey-api/database"
	"github.com/benpsk/go-survey-api/internal"
	"github.com/benpsk/go-survey-api/internal/handlers"
	"github.com/benpsk/go-survey-api/internal/repositories"
	"github.com/benpsk/go-survey-api/internal/services"
)

func main() {
	conn, err := database.Connect(config.DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	repo := repositories.New(conn)
	service := services.New(repo)
	handler := handlers.New(service)

	mux := internal.Router(handler)

	log.Printf("Server running on: %v", config.PORT)
	log.Fatal(http.ListenAndServe(":"+config.PORT, mux))
}
