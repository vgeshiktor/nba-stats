// internal/pkg/errors/error_handling.go
package errors

import (
	"encoding/json"
	"net/http"
)

// APIError represents a structured error response.
type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// WriteError writes a structured error response to the client.
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	apiError := APIError{
		Message: message,
		Code:    statusCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(apiError)
}
