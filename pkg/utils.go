package pkg

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data"`
}

func Success(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{Data: data}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		Error(w, "Failed to encode error message", http.StatusInternalServerError)
	}
}

// Error sends a general error message to the client.
func Error(w http.ResponseWriter, msg string, status ...int) {
	response := map[string]map[string]string{"errors": {"message": msg}}
	sendErrorResponse(w, response, status...)
}

// ValidationError sends validation error messages to the client.
func ValidationError(w http.ResponseWriter, validationErrors map[string]string, status ...int) {
	response := map[string]map[string]string{"errors": validationErrors}
	sendErrorResponse(w, response, status...)
}

// sendErrorResponse writes a JSON error response to the client.
func sendErrorResponse(w http.ResponseWriter, res any, status ...int) {
	code := http.StatusInternalServerError // Default status 500
	if len(status) > 0 {
		code = status[0] // Use provided status if available
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode error message", http.StatusInternalServerError)
	}
}

func Auth(r *http.Request) (int, bool) {
	userId, ok := r.Context().Value("userId").(float64)
	return int(userId), ok
}
