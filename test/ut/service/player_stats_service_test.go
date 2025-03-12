// test/ut/service/player_stats_service_test.go
package service_test

import (
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/service"
	"github.com/vgeshiktor/nba-stats/test/ut/mocks"
)

// --- Test functions ---

func TestLogPlayerStats_Success(t *testing.T) {
	// Create fake repositories.
	playerRepo := &mocks.FakePlayerRepo{}
	teamRepo := &mocks.FakeTeamRepo{}
	gameRepo := &mocks.FakeGameRepo{}
	statsRepo := &mocks.FakePlayerStatsRepo{}

	statsService := service.NewPlayerStatsService(playerRepo, teamRepo, gameRepo, statsRepo)

	// Valid player game statistics.
	stats := &domain.PlayerGameStats{
		ID:            "stats1",
		PlayerID:      "valid",
		GameID:        "game1",
		Points:        30,
		Rebounds:      5,
		Assists:       7,
		Steals:        2,
		Blocks:        1,
		Fouls:         3,
		Turnovers:     2,
		MinutesPlayed: 35.0,
	}

	err := statsService.LogPlayerStats(stats)
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	if !statsRepo.Inserted {
		t.Errorf("Expected stats to be inserted")
	}
}

func TestLogPlayerStats_InvalidFouls(t *testing.T) {
	playerRepo := &mocks.FakePlayerRepo{}
	teamRepo := &mocks.FakeTeamRepo{}
	gameRepo := &mocks.FakeGameRepo{}
	statsRepo := &mocks.FakePlayerStatsRepo{}

	statsService := service.NewPlayerStatsService(playerRepo, teamRepo, gameRepo, statsRepo)

	// Create stats with invalid fouls (> 6)
	stats := &domain.PlayerGameStats{
		ID:            "stats1",
		PlayerID:      "valid",
		GameID:        "game1",
		Points:        30,
		Rebounds:      5,
		Assists:       7,
		Steals:        2,
		Blocks:        1,
		Fouls:         7, // Invalid value.
		Turnovers:     2,
		MinutesPlayed: 35.0,
	}

	err := statsService.LogPlayerStats(stats)
	if err == nil {
		t.Errorf("Expected error due to invalid fouls, got success")
	}
}
