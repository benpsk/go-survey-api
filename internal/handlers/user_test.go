package handlers_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/benpsk/go-survey-api/internal/handlers"
	"github.com/benpsk/go-survey-api/internal/repositories"
	"github.com/benpsk/go-survey-api/internal/services"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	db        *pgx.Conn
	container *handlers.Handler
)

// Setup runs before all tests to initialize the DB connection.
func TestMain(m *testing.M) {
	var err error
	dsn := "postgres://postgres:postgres@localhost:5432/go_survey_api_test?sslmode=disable"
	db, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	repo := repositories.New(db)
	service := services.New(repo)
	container = handlers.New(service)

	os.Exit(m.Run())
}

// Clean up the database between tests.
func cleanDB() {
	_, _ = db.Exec(context.Background(), "DROP TABLE IF EXISTS users")
	_, _ = db.Exec(context.Background(), `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    )
  `)
}

// Test the User Handler by interacting with the real database.
func TestUserHandler_User_Success(t *testing.T) {
	cleanDB()
	_, _ = db.Exec(context.Background(), `
		INSERT INTO users (name, email, password) 
		VALUES ('John Doe', 'john@mail.com', 'password')
	`)
	req := httptest.NewRequest("GET", "/user", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userId", float64(1)))

	rr := httptest.NewRecorder()
	container.User(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status expected: %v, got: %v", http.StatusOK, rr.Code)
	}
	// Decode the expected JSON
	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"id":    float64(1),
			"name":  "John Doe",
			"email": "john@mail.com",
		},
	}

	var actual map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&actual)
	// Compare the JSON objects
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Body expected: %v, got: %v", expected, actual)
	}
}

func TestUserHandler_User_Unauthorized(t *testing.T) {
	cleanDB()
	_, _ = db.Exec(context.Background(), `
		INSERT INTO users (name, email, password) 
		VALUES ('John Doe', 'john@mail.com', 'password')
	`)
	req := httptest.NewRequest("GET", "/user", nil)
	rr := httptest.NewRecorder()
	container.User(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Status expected: %v, got: %v", http.StatusOK, rr.Code)
	}
	// Decode the expected JSON
	expected := map[string]interface{}{
		"errors": map[string]interface{}{
			"message": "Unauthorized",
		},
	}
	var actual map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&actual)
	// Compare the JSON objects
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Body expected: %v, got: %v", expected, actual)
	}
}

func TestUserHandler_User_Not_Found(t *testing.T) {
	cleanDB()
	_, _ = db.Exec(context.Background(), `
		INSERT INTO users (name, email, password) 
		VALUES ('John Doe', 'john@mail.com', 'password')
	`)
	req := httptest.NewRequest("GET", "/user", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userId", float64(100)))
	rr := httptest.NewRecorder()
	container.User(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Status expected: %v, got: %v", http.StatusOK, rr.Code)
	}
	// Decode the expected JSON
	expected := map[string]interface{}{
		"errors": map[string]interface{}{
			"message": "User not found",
		},
	}
	var actual map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&actual)
	// Compare the JSON objects
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Body expected: %v, got: %v", expected, actual)
	}
}

func TestUserHandler_Register_Validation_Success_Fail(t *testing.T) {
	cleanDB()
	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]interface{}
		status   int
	}{
		{
			name: "Name is required",
			input: map[string]string{
				"name":     "",
				"email":    "john@mail.com",
				"password": "password",
			},
			expected: map[string]interface{}{
				"errors": map[string]interface{}{
					"name": "Name is required",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Email is required",
			input: map[string]string{
				"name":     "John Doe",
				"email":    "",
				"password": "password",
			},
			expected: map[string]interface{}{
				"errors": map[string]interface{}{
					"email": "Email is required",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Email should be a valid email",
			input: map[string]string{
				"name":     "John Doe",
				"email":    "abc",
				"password": "password",
			},
			expected: map[string]interface{}{
				"errors": map[string]interface{}{
					"email": "Invalid email format",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Password is required",
			input: map[string]string{
				"name":     "John Doe",
				"email":    "john@mail.com",
				"password": "",
			},
			expected: map[string]interface{}{
				"errors": map[string]interface{}{
					"password": "Password is required",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "All fields are required",
			input: map[string]string{
				"name":     "",
				"email":    "",
				"password": "",
			},
			expected: map[string]interface{}{
				"errors": map[string]interface{}{
					"name":     "Name is required",
					"email":    "Email is required",
					"password": "Password is required",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Max character should not exceed",
			input: map[string]string{
				"name":     "The error indicates that the value passed to bcrypt.CompareHashAndPassword as the hashis not a valid bcrypt hash. This usually happens when you try to compare a plain text password as if it were a bcrypt hash.",
				"email":    "The error indicates that the value passed to bcrypt.CompareHashAndPassword as the hashis not a valid bcrypt hash. This usually happens when you try to compare a plain text password as if it were a bcrypt hash@mail.com",
				"password": "The error indicates that the value passed to bcrypt.CompareHashAndPassword as the hashis not a valid bcrypt hash. This usually happens when you try to compare a plain text password as if it were a bcrypt hash.",
			},
			expected: map[string]interface{}{
				"errors": map[string]interface{}{
					"name":     "Name must not exceed 100 characters",
					"email":    "Email must not exceed 100 characters",
					"password": "Password must not exceed 100 characters",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Valid input",
			input: map[string]string{
				"name":     "John Doe",
				"email":    "john@mail.com",
				"password": "password",
			},
			expected: map[string]interface{}{
				"data": map[string]interface{}{
					"id":    float64(1),
					"name":  "John Doe",
					"email": "john@mail.com",
				},
			},
			status: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest("POST", "/register", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			container.Register(rr, req)

			var response map[string]interface{}
			json.NewDecoder(rr.Body).Decode(&response)

			if rr.Code != tc.status {
				t.Errorf("Status expected: %v, got: %v", tc.status, rr.Code)
			}
			if !reflect.DeepEqual(tc.expected, response) {
				t.Errorf("Body expected: %v, got: %v", tc.expected, response)
			}
		})
	}
}

func TestUserHandler_Register_Password_Hashed(t *testing.T) {
	cleanDB()
	input := map[string]string{
		"name":     "John Doe",
		"email":    "john@mail.com",
		"password": "password",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/register", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	container.Register(rr, req)

	var actualPassword string
	err := db.QueryRow(context.Background(), "SELECT password FROM users WHERE id=$1", 1).Scan(&actualPassword)
	if err != nil {
		t.Errorf("Query error: %v", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(actualPassword), []byte(input["password"]))
	if err != nil {
		t.Errorf("Password hashed not match err: %v", err)
	}
	if rr.Code != http.StatusCreated {
		t.Errorf("Status expected: %v, got: %v", http.StatusCreated, rr.Code)
	}
}

func TestUserHandler_Login_Validation(t *testing.T) {
	cleanDB()
	_, _ = db.Exec(context.Background(), `
		INSERT INTO users (name, email, password) 
		VALUES ('John Doe', 'john@mail.com', 'password')
	`)
	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]interface{}
		status   int
	}{
		{
			name: "Input are required",
			input: map[string]string{
				"email":    "",
				"password": "",
			},
			expected: map[string]interface{}{
				"errors": map[string]interface{}{
					"email":    "Email is required",
					"password": "Password is required",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "User not found",
			input: map[string]string{
				"email":    "hello@mail.com",
				"password": "password",
			},
			expected: map[string]interface{}{
				"errors": map[string]interface{}{
					"message": "User not found",
				},
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Password does not match",
			input: map[string]string{
				"email":    "john@mail.com",
				"password": "a",
			},
			expected: map[string]interface{}{
				"errors": map[string]interface{}{
					"message": "Password does not match",
				},
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest("POST", "/login", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			container.Login(rr, req)

			var response map[string]interface{}
			json.NewDecoder(rr.Body).Decode(&response)

			if rr.Code != tc.status {
				t.Errorf("Status expected: %v, got: %v", tc.status, rr.Code)
			}
			if !reflect.DeepEqual(tc.expected, response) {
				t.Errorf("Body expected: %v, got: %v", tc.expected, response)
			}
		})
	}
}

func TestUserHandler_Login_Success(t *testing.T) {
	cleanDB()
	_, _ = db.Exec(context.Background(), `
		INSERT INTO users (name, email, password) 
		VALUES ('John Doe', 'john@mail.com', '$2a$10$TbKnsnI8pB/KovbdnbvbbOybf1SESd0o8nB7y/iCwkYtoLa2vhjiu')
	`)

	input := map[string]string{
		"email":    "john@mail.com",
		"password": "password",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	container.Login(rr, req)

	type tokenResponse struct {
		Data struct {
			Type        string  `json:"type"`
			AccessToken string  `json:"access_token"`
			ExpiredAt   float64 `json:"expired_at"`
		} `json:"data"`
	}
	var response tokenResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if rr.Code != http.StatusOK {
		t.Errorf("Status expected: %v, got: %v", http.StatusOK, rr.Code)
	}
	if response.Data.Type != "Bearer" {
		t.Errorf("Token type expected: %v, got: %v", "Bearer", response.Data.Type)
	}
	if response.Data.AccessToken == "" {
		t.Errorf("Access token got: %v", response.Data.Type)
	}
	if response.Data.ExpiredAt == 0 {
		t.Errorf("Access token got: %v", response.Data.ExpiredAt)
	}
}
