package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/app"
	"github.com/vgeshiktor/nba-stats/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetPlayer(t *testing.T) {
	// Set up environment for in-memory SQLite testing.
	os.Setenv("DATABASE_URL", ":memory:")

	// Initialize the application (runs migrations, sets up router, etc.)
	server := app.Initialize()

	// Create a new player.
	newPlayer := domain.Player{
		ID:     "player1",
		Name:   "John Doe",
		TeamID: "team1",
	}
	playerBody, err := json.Marshal(newPlayer)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/players", bytes.NewBuffer(playerBody))
	req.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	resp := httptest.NewRecorder()
	server.Handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Retrieve the created player.
	reqGet, _ := http.NewRequest("GET", "/api/v1/players/player1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqGet.Header.Set("Authorization", "dummy-token")

	respGet := httptest.NewRecorder()
	server.Handler.ServeHTTP(respGet, reqGet)
	assert.Equal(t, http.StatusOK, respGet.Code)

	var retrievedPlayer domain.Player
	err = json.Unmarshal(respGet.Body.Bytes(), &retrievedPlayer)
	assert.NoError(t, err)
	assert.Equal(t, newPlayer.ID, retrievedPlayer.ID)
	assert.Equal(t, newPlayer.Name, retrievedPlayer.Name)
	assert.Equal(t, newPlayer.TeamID, retrievedPlayer.TeamID)
}
