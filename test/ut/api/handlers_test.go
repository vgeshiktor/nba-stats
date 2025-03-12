// test/ut/api/handlers_test.go
package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/api"
	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/test/ut/mocks"
)



// --- Tests for API Endpoints ---

func TestLogPlayerStatsEndpoint(t *testing.T) {
	// Create a fake API handler using the fake services.
	handler := api.NewHandler(
		&mocks.FakePlayerStatsService{},
		&mocks.FakeAggregationService{},
		&mocks.FakePlayerService{},
		&mocks.FakeTeamService{},
		&mocks.FakeGameService{},
	)

	// Create a sample PlayerGameStats payload.
	stats := &domain.PlayerGameStats{
		ID:            "stats1",
		PlayerID:      "player1",
		GameID:        "game1",
		Points:        25,
		Rebounds:      5,
		Assists:       7,
		Steals:        2,
		Blocks:        1,
		Fouls:         3,
		Turnovers:     2,
		MinutesPlayed: 32.5,
	}

	payload, _ := json.Marshal(stats)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/player-stats", bytes.NewReader(payload))
	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	rr := httptest.NewRecorder()
	handler.LogPlayerStats(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, status)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if msg, ok := resp["message"]; !ok || !strings.Contains(msg, "logged successfully") {
		t.Errorf("unexpected response message: %v", resp)
	}
}

func TestGetPlayerAggregateEndpoint(t *testing.T) {
	handler := api.NewHandler(
		&mocks.FakePlayerStatsService{},
		&mocks.FakeAggregationService{},
		&mocks.FakePlayerService{},
		&mocks.FakeTeamService{},
		&mocks.FakeGameService{},
	)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/player-stats/player/player1", nil)
	req.Header.Set("Authorization", "dummy-token")
	rr := httptest.NewRecorder()

	handler.GetPlayerAggregate(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	var agg domain.AggregateStats
	if err := json.NewDecoder(rr.Body).Decode(&agg); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if agg.PlayerID != "player1" {
		t.Errorf("expected playerID 'player1', got %s", agg.PlayerID)
	}
}
