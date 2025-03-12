package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/app"

	"github.com/stretchr/testify/assert"
)

func TestCreatePlayer(t *testing.T) {
	server := app.Initialize()

	playerData := map[string]string{
		"id":      "player1",
		"name":    "John Doe",
		"team_id": "team1",
	}

	body, _ := json.Marshal(playerData)
	req, _ := http.NewRequest("POST", "/api/v1/players", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	resp := httptest.NewRecorder()
	server.Handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), `"id":"player1"`)
}

func TestGetPlayerByID(t *testing.T) {
	server := app.Initialize()

		playerData := map[string]string{
		"id":      "player2",
		"name":    "Mike Tyson",
		"team_id": "team1",
	}

	body, _ := json.Marshal(playerData)
	req, _ := http.NewRequest("POST", "/api/v1/players", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	resp := httptest.NewRecorder()
	server.Handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), `"id":"player2"`)

	req, _ = http.NewRequest("GET", "/api/v1/players/player2", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	resp = httptest.NewRecorder()

	server.Handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"id":"player2"`)
}
