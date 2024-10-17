package main

import (
	"context"
	"log"
	"net/http"

	"github.com/benpsk/go-survey-api/config"
	"github.com/benpsk/go-survey-api/database"
	"github.com/benpsk/go-survey-api/internal"
)

func main() {
	conn, err := database.Connect(config.DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	container := internal.NewContainer(conn)
	mux := internal.Router(container)

	log.Printf("Server running on: %v", config.PORT)
	log.Fatal(http.ListenAndServe(":"+config.PORT, mux))
}
