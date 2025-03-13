package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/api"

	"github.com/stretchr/testify/assert"
)

// Test Logging Middleware
func TestLoggingMiddleware(t *testing.T) {
	// Create a simple handler that returns OK
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	})

	// Wrap the handler with middleware
	handler := api.LoggingMiddleware(nextHandler)

	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"message": "success"`)
}

// Test Auth Middleware - Success Case
func TestAuthMiddleware_Success(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "authorized"}`))
	})

	handler := api.AuthenticationMiddleware(nextHandler)

	req, _ := http.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"message": "authorized"`)
}

// Test Auth Middleware - No Token
func TestAuthMiddleware_NoToken(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := api.AuthenticationMiddleware(nextHandler)

	req, _ := http.NewRequest("GET", "/secure", nil) // No Auth Header
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), `"Unauthorized: missing token"`)
}

// TODO: add token validation and uncomment this test
// Test Auth Middleware - Invalid Token
// func TestAuthMiddleware_InvalidToken(t *testing.T) {
// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	})

// 	handler := api.AuthenticationMiddleware(nextHandler)

// 	req, _ := http.NewRequest("GET", "/secure", nil)
// 	req.Header.Set("Authorization", "Bearer invalid-token") // Invalid token
// 	resp := httptest.NewRecorder()

// 	handler.ServeHTTP(resp, req)

// 	assert.Equal(t, http.StatusUnauthorized, resp.Code)
// 	assert.Contains(t, resp.Body.String(), `"error": "Invalid token"`)
// }
