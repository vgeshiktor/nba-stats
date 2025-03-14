// test/e2e/player_stats_e2e_test.go
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

func TestPlayerStatsE2E(t *testing.T) {
	// --- Step 1: Create a Team ---
	newTeam := domain.Team{
		ID:   "team1",
		Name: "Team One",
	}
	teamBody, err := json.Marshal(newTeam)
	assert.NoError(t, err)
	reqTeam, _ := http.NewRequest("POST", "/api/v1/teams", bytes.NewBuffer(teamBody))
	reqTeam.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	reqTeam.Header.Set("Authorization", "dummy-token")

	respTeam := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(respTeam, reqTeam)
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
	testServer.Handler.ServeHTTP(respPlayer, reqPlayer)
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
	testServer.Handler.ServeHTTP(respGame, reqGame)
	assert.Equal(t, http.StatusCreated, respGame.Code)

	// --- Step 4: Log Two Player Stats Entries ---
	statsEntries := []domain.PlayerGameStats{
		{
			ID:            "stats1",
			PlayerID:      "player1",
			GameID:        "game1",
			Points:        25,
			Rebounds:      5,
			Assists:       4,
			Steals:        2,
			Blocks:        1,
			Fouls:         2,
			Turnovers:     3,
			MinutesPlayed: 30.0,
		},
		{
			ID:            "stats2",
			PlayerID:      "player1",
			GameID:        "game1", // Same game to test aggregation grouping
			Points:        35,
			Rebounds:      7,
			Assists:       6,
			Steals:        1,
			Blocks:        0,
			Fouls:         3,
			Turnovers:     2,
			MinutesPlayed: 40.0,
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
		testServer.Handler.ServeHTTP(respStats, reqStats)
		assert.Equal(t, http.StatusCreated, respStats.Code)
	}

	// --- Step 5: Retrieve and Verify Player Aggregate Stats ---
	reqAgg, _ := http.NewRequest("GET", "/api/v1/player-stats/player/player1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqAgg.Header.Set("Authorization", "dummy-token")

	respAgg := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(respAgg, reqAgg)
	assert.Equal(t, http.StatusOK, respAgg.Code)

	var agg domain.AggregateStats
	err = json.Unmarshal(respAgg.Body.Bytes(), &agg)
	assert.NoError(t, err)

	// Since both stat entries belong to the same game, GamesPlayed should be 1.
	assert.Equal(t, 1, agg.GamesPlayed)

	// Totals are the sum of each field from both entries.
	expectedTotalPoints := 25 + 35
	expectedTotalRebounds := 5 + 7
	expectedTotalAssists := 4 + 6
	expectedTotalSteals := 2 + 1
	expectedTotalBlocks := 1 + 0
	expectedTotalFouls := 2 + 3
	expectedTotalTurnovers := 3 + 2
	expectedTotalMinutes := 30.0 + 40.0

	assert.Equal(t, expectedTotalPoints, agg.TotalPoints)
	assert.Equal(t, expectedTotalRebounds, agg.TotalRebounds)
	assert.Equal(t, expectedTotalAssists, agg.TotalAssists)
	assert.Equal(t, expectedTotalSteals, agg.TotalSteals)
	assert.Equal(t, expectedTotalBlocks, agg.TotalBlocks)
	assert.Equal(t, expectedTotalFouls, agg.TotalFouls)
	assert.Equal(t, expectedTotalTurnovers, agg.TotalTurnovers)
	assert.Equal(t, expectedTotalMinutes, agg.TotalMinutes)

	// Averages should be totals divided by games played (i.e., totals themselves since GamesPlayed is 1).
	assert.Equal(t, float64(expectedTotalPoints)/1.0, agg.AvgPoints)
	assert.Equal(t, float64(expectedTotalRebounds)/1.0, agg.AvgRebounds)
	assert.Equal(t, float64(expectedTotalAssists)/1.0, agg.AvgAssists)
	assert.Equal(t, float64(expectedTotalSteals)/1.0, agg.AvgSteals)
	assert.Equal(t, float64(expectedTotalBlocks)/1.0, agg.AvgBlocks)
	assert.Equal(t, float64(expectedTotalFouls)/1.0, agg.AvgFouls)
	assert.Equal(t, float64(expectedTotalTurnovers)/1.0, agg.AvgTurnovers)
	assert.Equal(t, expectedTotalMinutes/1.0, agg.AvgMinutes)
}
