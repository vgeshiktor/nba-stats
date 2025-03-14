package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test unauthorized access
func TestUnauthorizedAccessE2E(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/player-stats/player/player1", nil) // Secure endpoint
	resp := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), `Unauthorized`)
}

// Test access with valid token
func TestAuthorizedAccessE2E(t *testing.T) {
		// Step 1: Create a player
	newPlayer := map[string]string{
		"id":      "player1",
		"name":    "John Doe",
		"team_id": "team1",
	}
	playerBody, _ := json.Marshal(newPlayer)

	req, _ := http.NewRequest("POST", "/api/v1/players", bytes.NewBuffer(playerBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token") // Assuming valid token logic exists

	resp := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}
