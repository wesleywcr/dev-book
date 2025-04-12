package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error sends an error response with the specified status code and error message.
func Error(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

// JSON sends a JSON response with the specified status code and data.

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
