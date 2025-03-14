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

	"github.com/stretchr/testify/assert"
)

func TestLogPlayerStatsAndGetAggregate(t *testing.T) {
	os.Setenv("DATABASE_URL", ":memory:")

	server := app.Initialize()

	// Create a player (required for logging stats).
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

	// Create a game (required for logging stats).
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

	// Log player stats.
	newStats := domain.PlayerGameStats{
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
	}
	statsBody, err := json.Marshal(newStats)
	assert.NoError(t, err)
	reqStats, _ := http.NewRequest("POST", "/api/v1/player-stats", bytes.NewBuffer(statsBody))
	reqStats.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	reqStats.Header.Set("Authorization", "dummy-token")

	respStats := httptest.NewRecorder()
	server.Handler.ServeHTTP(respStats, reqStats)
	assert.Equal(t, http.StatusCreated, respStats.Code)

	// Retrieve player aggregate stats.
	reqAgg, _ := http.NewRequest("GET", "/api/v1/player-stats/player/player1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqAgg.Header.Set("Authorization", "dummy-token")

	respAgg := httptest.NewRecorder()
	server.Handler.ServeHTTP(respAgg, reqAgg)
	assert.Equal(t, http.StatusOK, respAgg.Code)

	var agg domain.AggregateStats
	err = json.Unmarshal(respAgg.Body.Bytes(), &agg)
	assert.NoError(t, err)

	// Since there is only one game, totals should equal the logged stats.
	assert.Equal(t, 1, agg.GamesPlayed)
	assert.Equal(t, newStats.Points, agg.TotalPoints)
	// Average equals total divided by games played.
	expectedAvgPoints := float64(newStats.Points) / float64(agg.GamesPlayed)
	assert.Equal(t, expectedAvgPoints, agg.AvgPoints)
}
