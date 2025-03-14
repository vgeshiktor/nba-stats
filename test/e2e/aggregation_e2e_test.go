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

func TestAggregationFlowE2E(t *testing.T) {
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

	// --- Step 2: Create Two Players ---
	players := []domain.Player{
		{ID: "player1", Name: "John Doe", TeamID: "team1"},
		{ID: "player2", Name: "Jane Doe", TeamID: "team1"},
	}

	for _, player := range players {
		playerBody, err := json.Marshal(player)
		assert.NoError(t, err)
		reqPlayer, _ := http.NewRequest("POST", "/api/v1/players", bytes.NewBuffer(playerBody))
		reqPlayer.Header.Set("Content-Type", "application/json")

		// For endpoints with middleware, add a dummy Authorization header.
		reqPlayer.Header.Set("Authorization", "dummy-token")

		respPlayer := httptest.NewRecorder()
		testServer.Handler.ServeHTTP(respPlayer, reqPlayer)
		assert.Equal(t, http.StatusCreated, respPlayer.Code)
	}

	// --- Step 3: Create Two Games ---
	games := []domain.Game{
		{ID: "game1", Date: time.Now(), HomeTeam: "team1", AwayTeam: "team2"},
		{ID: "game2", Date: time.Now().Add(24 * time.Hour), HomeTeam: "team1", AwayTeam: "team3"},
	}

	for _, game := range games {
		gameBody, err := json.Marshal(game)
		assert.NoError(t, err)
		reqGame, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBuffer(gameBody))
		reqGame.Header.Set("Content-Type", "application/json")

		// For endpoints with middleware, add a dummy Authorization header.
		reqGame.Header.Set("Authorization", "dummy-token")

		respGame := httptest.NewRecorder()
		testServer.Handler.ServeHTTP(respGame, reqGame)
		assert.Equal(t, http.StatusCreated, respGame.Code)
	}

	// --- Step 4: Log Player Stats for Two Games ---
	statsEntries := []domain.PlayerGameStats{
		{ID: "stats1", PlayerID: "player1", GameID: "game1", Points: 30, Rebounds: 8, Assists: 5, Steals: 2, Blocks: 1, Fouls: 3, Turnovers: 2, MinutesPlayed: 35.0},
		{ID: "stats2", PlayerID: "player1", GameID: "game2", Points: 20, Rebounds: 6, Assists: 4, Steals: 1, Blocks: 0, Fouls: 2, Turnovers: 3, MinutesPlayed: 30.0},
		{ID: "stats3", PlayerID: "player2", GameID: "game1", Points: 15, Rebounds: 5, Assists: 3, Steals: 1, Blocks: 1, Fouls: 1, Turnovers: 1, MinutesPlayed: 25.0},
		{ID: "stats4", PlayerID: "player2", GameID: "game2", Points: 10, Rebounds: 4, Assists: 2, Steals: 0, Blocks: 0, Fouls: 2, Turnovers: 2, MinutesPlayed: 20.0},
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

	// --- Step 5: Retrieve and Verify Player Aggregated Stats ---
	reqAggPlayer, _ := http.NewRequest("GET", "/api/v1/player-stats/player/player1", nil)
		
	// For endpoints with middleware, add a dummy Authorization header.
	reqAggPlayer.Header.Set("Authorization", "dummy-token")

	respAggPlayer := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(respAggPlayer, reqAggPlayer)
	assert.Equal(t, http.StatusOK, respAggPlayer.Code)

	var playerAgg domain.AggregateStats
	err = json.Unmarshal(respAggPlayer.Body.Bytes(), &playerAgg)
	assert.NoError(t, err)

	// Expected stats for player1
	assert.Equal(t, 2, playerAgg.GamesPlayed)
	assert.Equal(t, 30+20, playerAgg.TotalPoints)
	assert.Equal(t, (30+20)/2.0, playerAgg.AvgPoints)
	assert.Equal(t, 8+6, playerAgg.TotalRebounds)
	assert.Equal(t, (8+6)/2.0, playerAgg.AvgRebounds)

	// --- Step 6: Retrieve and Verify Team Aggregated Stats ---
	reqAggTeam, _ := http.NewRequest("GET", "/api/v1/player-stats/team/team1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqAggTeam.Header.Set("Authorization", "dummy-token")

	respAggTeam := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(respAggTeam, reqAggTeam)
	assert.Equal(t, http.StatusOK, respAggTeam.Code)

	var teamAgg domain.AggregateStats
	err = json.Unmarshal(respAggTeam.Body.Bytes(), &teamAgg)
	assert.NoError(t, err)

	// Expected stats for team1 (Sum of player1 and player2 stats)
	expectedTeamPoints := (30 + 20) + (15 + 10)
	expectedTeamRebounds := (8 + 6) + (5 + 4)

	assert.Equal(t, "team1", teamAgg.TeamID)
	assert.Equal(t, 2, teamAgg.GamesPlayed)
	assert.Equal(t, expectedTeamPoints, teamAgg.TotalPoints)
	assert.Equal(t, float64(expectedTeamPoints)/2.0, teamAgg.AvgPoints)
	assert.Equal(t, expectedTeamRebounds, teamAgg.TotalRebounds)
	assert.Equal(t, float64(expectedTeamRebounds)/2.0, teamAgg.AvgRebounds)
}
