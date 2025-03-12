// test/ut/service/player_service_test.go
package service_test

import (
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/service"
	"github.com/vgeshiktor/nba-stats/test/ut/mocks"
)

func TestCreatePlayer_Success(t *testing.T) {
	// Use FakePlayerRepo from our mocks package.
	fakeRepo := &mocks.FakePlayerRepo{}
	playerService := service.NewPlayerService(fakeRepo)

	player := &domain.Player{
		ID:     "valid",
		Name:   "John Doe",
		TeamID: "team1",
	}

	err := playerService.CreatePlayer(player)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestCreatePlayer_Invalid(t *testing.T) {
	fakeRepo := &mocks.FakePlayerRepo{}
	playerService := service.NewPlayerService(fakeRepo)

	// Test with missing Name.
	player := &domain.Player{
		ID:     "valid",
		Name:   "",
		TeamID: "team1",
	}

	err := playerService.CreatePlayer(player)
	if err == nil {
		t.Errorf("expected error for missing name, got nil")
	}
}

func TestGetPlayerByID_Success(t *testing.T) {
	fakeRepo := &mocks.FakePlayerRepo{}
	playerService := service.NewPlayerService(fakeRepo)

	player, err := playerService.GetPlayerByID("valid")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if player.ID != "valid" {
		t.Errorf("expected player ID 'valid', got: %s", player.ID)
	}
}

func TestGetPlayerByID_Invalid(t *testing.T) {
	fakeRepo := &mocks.FakePlayerRepo{}
	playerService := service.NewPlayerService(fakeRepo)

	_, err := playerService.GetPlayerByID("invalid")
	if err == nil {
		t.Errorf("expected error for invalid player ID, got nil")
	}
}
