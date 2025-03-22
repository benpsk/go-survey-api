package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/benpsk/go-survey-api/database"
	"github.com/jackc/pgx/v5"
)

const (
	migrationFile = "database/migration.sql"
	seedFile      = "database/seed.sql"
)

func main() {
	// Define CLI commands
	migrateCmd := flag.Bool("migrate", false, "Run database migrations")
	seedCmd := flag.Bool("seed", false, "Run database seeding")

	flag.Parse() // Parse command-line flags

	// Show help message if no flag is provided
	if !*migrateCmd && !*seedCmd {
		fmt.Println("Usage:")
		fmt.Println("  go run cmd/cli/main.go -migrate   # Run database migrations")
		fmt.Println("  go run cmd/cli/main.go -seed      # Run database seeding")
		os.Exit(1) // Exit with error code 1
	}

	db, err := database.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	// Absolute paths for the SQL files
	migrationFile, _ := filepath.Abs("database/migration.sql")
	seedFile, _ := filepath.Abs("database/seed.sql")

	// Execute based on the command
	if *migrateCmd {
		executeSQLFile(db, migrationFile, "Migration")
	}
	if *seedCmd {
		executeSQLFile(db, seedFile, "Seeding")
	}
}

// Reads and executes SQL from a file
func executeSQLFile(db *pgx.Conn, filePath string, operation string) {
	sqlData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading %s file: %v", operation, err)
	}

	_, err = db.Exec(context.Background(), string(sqlData))
	if err != nil {
		log.Fatalf("%s failed: %v", operation, err)
	}

	fmt.Printf("%s completed successfully!\n", operation)
}
