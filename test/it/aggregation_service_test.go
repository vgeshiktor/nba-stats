// test/integration/aggregation_service_test.go
package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/vgeshiktor/nba-stats/internal/app"
	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestAggregationServiceIntegration(t *testing.T) {
	// Use in-memory SQLite database for testing.
	os.Setenv("DATABASE_URL", ":memory:")

	// Initialize the app (runs migrations, sets up routes, etc.)
	server := app.Initialize()

	// --- Step 1: Create a Team ---
	newTeam := domain.Team{
		ID:   "team1",
		Name: "Test Team",
	}
	teamBody, err := json.Marshal(newTeam)
	assert.NoError(t, err)
	reqTeam, _ := http.NewRequest("POST", "/api/v1/teams", bytes.NewBuffer(teamBody))
	reqTeam.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	reqTeam.Header.Set("Authorization", "dummy-token")

	respTeam := httptest.NewRecorder()
	server.Handler.ServeHTTP(respTeam, reqTeam)
	assert.Equal(t, http.StatusCreated, respTeam.Code)

	// --- Step 2: Create a Player ---
	newPlayer := domain.Player{
		ID:     "player1",
		Name:   "John Doe",
		TeamID: "team1",
	}
	playerBody, err := json.Marshal(newPlayer)
	assert.NoError(t, err)
	reqPlayer, _ := http.NewRequest("POST", "/api/v1/players", bytes.NewBuffer(playerBody))
	reqPlayer.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	reqPlayer.Header.Set("Authorization", "dummy-token")

	respPlayer := httptest.NewRecorder()
	server.Handler.ServeHTTP(respPlayer, reqPlayer)
	assert.Equal(t, http.StatusCreated, respPlayer.Code)

	// --- Step 3: Create a Game ---
	newGame := domain.Game{
		ID:       "game1",
		Date:     time.Now(),
		HomeTeam: "team1",
		AwayTeam: "team2",
	}
	gameBody, err := json.Marshal(newGame)
	assert.NoError(t, err)
	reqGame, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(gameBody))
	reqGame.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	reqGame.Header.Set("Authorization", "dummy-token")

	respGame := httptest.NewRecorder()
	server.Handler.ServeHTTP(respGame, reqGame)
	assert.Equal(t, http.StatusCreated, respGame.Code)

	// --- Step 4: Log Multiple Player Stats ---
	// Log stats from two games for the same player
	statsEntries := []domain.PlayerGameStats{
		{
			ID:            "stats1",
			PlayerID:      "player1",
			GameID:        "game1",
			Points:        30,
			Rebounds:      8,
			Assists:       5,
			Steals:        2,
			Blocks:        1,
			Fouls:         3,
			Turnovers:     2,
			MinutesPlayed: 35.0,
		},
		{
			ID:            "stats2",
			PlayerID:      "player1",
			GameID:        "game1", // For simplicity, using same game ID
			Points:        20,
			Rebounds:      6,
			Assists:       4,
			Steals:        1,
			Blocks:        0,
			Fouls:         2,
			Turnovers:     3,
			MinutesPlayed: 30.0,
		},
	}

	for _, stats := range statsEntries {
		statsBody, err := json.Marshal(stats)
		assert.NoError(t, err)
		reqStats, _ := http.NewRequest("POST", "/api/v1/player-stats", bytes.NewBuffer(statsBody))
		reqStats.Header.Set("Content-Type", "application/json")

		// For endpoints with middleware, add a dummy Authorization header.
		reqStats.Header.Set("Authorization", "dummy-token")

		respStats := httptest.NewRecorder()
		server.Handler.ServeHTTP(respStats, reqStats)
		assert.Equal(t, http.StatusCreated, respStats.Code)
	}

	// --- Step 5: Retrieve and Verify Player Aggregate Stats ---
	reqAggPlayer, _ := http.NewRequest("GET", "/api/v1/player-stats/player/player1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqAggPlayer.Header.Set("Authorization", "dummy-token")
	
	respAggPlayer := httptest.NewRecorder()
	server.Handler.ServeHTTP(respAggPlayer, reqAggPlayer)
	assert.Equal(t, http.StatusOK, respAggPlayer.Code)

	var playerAgg domain.AggregateStats
	err = json.Unmarshal(respAggPlayer.Body.Bytes(), &playerAgg)
	assert.NoError(t, err)

	jsonPlayerAgg, _ := json.MarshalIndent(playerAgg, "", "")
	logger.Info("Player agg stats: %s", jsonPlayerAgg)

	// Since we logged two stat entries for one game, GamesPlayed should be 1 if the aggregation is grouped per game.
	assert.Equal(t, 1, playerAgg.GamesPlayed)
	// Totals should be the sum of both entries.
	expectedTotalPoints := 30 + 20
	assert.Equal(t, expectedTotalPoints, playerAgg.TotalPoints)
	// Average is computed as total divided by games played.
	expectedAvgPoints := float64(expectedTotalPoints) / float64(playerAgg.GamesPlayed)
	assert.Equal(t, expectedAvgPoints, playerAgg.AvgPoints)

	// --- Step 6: Retrieve and Verify Team Aggregate Stats ---
	reqAggTeam, _ := http.NewRequest("GET", "/api/v1/player-stats/team/team1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqAggTeam.Header.Set("Authorization", "dummy-token")

	respAggTeam := httptest.NewRecorder()
	server.Handler.ServeHTTP(respAggTeam, reqAggTeam)
	assert.Equal(t, http.StatusOK, respAggTeam.Code)

	var teamAgg domain.AggregateStats
	err = json.Unmarshal(respAggTeam.Body.Bytes(), &teamAgg)
	assert.NoError(t, err)
	assert.Equal(t, "team1", teamAgg.TeamID)
	// Totals and averages for the team should reflect the aggregated stats of its players.
	assert.Equal(t, expectedTotalPoints, teamAgg.TotalPoints)
}
