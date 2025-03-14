package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetPlayerE2E(t *testing.T) {
	// Step 1: Create a player
	newPlayer := map[string]string{
		"id":      "player1",
		"name":    "John Doe",
		"team_id": "team1",
	}
	playerBody, _ := json.Marshal(newPlayer)

	req, _ := http.NewRequest("POST", "/api/v1/players", bytes.NewBuffer(playerBody))
	req.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	resp := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	// Step 2: Retrieve the player
	reqGet, _ := http.NewRequest("GET", "/api/v1/players/player1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqGet.Header.Set("Authorization", "dummy-token")

	respGet := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(respGet, reqGet)

	assert.Equal(t, http.StatusOK, respGet.Code)
	assert.Contains(t, respGet.Body.String(), `"id":"player1"`)
}

// Test fetching a non-existent player
func TestGetNonExistentPlayerE2E(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/players/nonexistent", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	resp := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), `Player not found`)
}
