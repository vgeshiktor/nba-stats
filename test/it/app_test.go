package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vgeshiktor/nba-stats/internal/app"
)

func TestAppInitialization(t *testing.T) {
	// Initialize the app (this starts everything: DB, middleware, router)
	server := app.Initialize()

	// Create a request to test the health check endpoint
	req, _ := http.NewRequest("GET", "/health/live", nil)
	resp := httptest.NewRecorder()
	server.Handler.ServeHTTP(resp, req)

	// Ensure the app responds correctly
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"status":"ok"`)
}
