// test/ut/service/game_service_test.go
package service_test

import (
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/service"
	"github.com/vgeshiktor/nba-stats/test/ut/mocks"
)

func TestCreateGame_Success(t *testing.T) {
	repo := &mocks.FakeGameRepo{}
	gameService := service.NewGameService(repo)

	game := &domain.Game{
		ID:       "game1",
		HomeTeam: "team1",
		AwayTeam: "team2",
	}

	err := gameService.CreateGame(game)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestCreateGame_InvalidInput(t *testing.T) {
	repo := &mocks.FakeGameRepo{}
	gameService := service.NewGameService(repo)

	// Test with missing game ID.
	game := &domain.Game{
		ID:       "",
		HomeTeam: "team1",
		AwayTeam: "team2",
	}

	err := gameService.CreateGame(game)
	if err == nil {
		t.Errorf("expected error for missing game ID, got success")
	}
}

func TestGetGameByID_Success(t *testing.T) {
	repo := &mocks.FakeGameRepo{}
	gameService := service.NewGameService(repo)

	game, err := gameService.GetGameByID("game1")
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
	if game.ID != "game1" {
		t.Errorf("expected game ID 'game1', got %s", game.ID)
	}
}

func TestGetGameByID_InvalidID(t *testing.T) {
	repo := &mocks.FakeGameRepo{}
	gameService := service.NewGameService(repo)

	_, err := gameService.GetGameByID("")
	if err == nil {
		t.Errorf("expected error for empty game ID, got success")
	}
}
