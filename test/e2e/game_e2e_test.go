// test/e2e/game_e2e_test.go
package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vgeshiktor/nba-stats/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetGameE2E(t *testing.T) {
	// --- Step 1: Create a new game ---
	newGame := domain.Game{
		ID:       "game1",
		Date:     time.Now(),
		HomeTeam: "team1",
		AwayTeam: "team2",
	}
	gameBody, err := json.Marshal(newGame)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(gameBody))
	req.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	resp := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// --- Step 2: Retrieve the created game ---
	reqGet, _ := http.NewRequest("GET", "/api/v1/games/game1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqGet.Header.Set("Authorization", "dummy-token")

	respGet := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(respGet, reqGet)
	assert.Equal(t, http.StatusOK, respGet.Code)

	var retrievedGame domain.Game
	err = json.Unmarshal(respGet.Body.Bytes(), &retrievedGame)
	assert.NoError(t, err)

	// --- Assert that the retrieved game matches the created game ---
	assert.Equal(t, newGame.ID, retrievedGame.ID)
	assert.Equal(t, newGame.HomeTeam, retrievedGame.HomeTeam)
	assert.Equal(t, newGame.AwayTeam, retrievedGame.AwayTeam)
	// Optionally, you can check that the Date field is non-zero.
	assert.False(t, retrievedGame.Date.IsZero(), "Game date should be set")
}
