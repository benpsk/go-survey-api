package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/benpsk/go-survey-api/database"
	"github.com/benpsk/go-survey-api/internal"
	"github.com/benpsk/go-survey-api/internal/handlers"
	"github.com/benpsk/go-survey-api/internal/repositories"
	"github.com/benpsk/go-survey-api/internal/services"
)

func main() {
	if os.Getenv("ENV") == "PRODUCTION" {
		logFile := generateLog()
		defer logFile.Close()
		os.Stderr = logFile
		log.SetOutput(logFile)
	}

	conn, err := database.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	repo := repositories.New(conn)
	service := services.New(repo)
	handler := handlers.New(service)

	mux := internal.Router(handler)

	log.Printf("Server running on: %v", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), mux))
}

func generateLog() *os.File {
	date := time.Now().Format("2006-01-02")
	fileName := filepath.Join("logs", "errors-"+date+".log")
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	return logFile
}
