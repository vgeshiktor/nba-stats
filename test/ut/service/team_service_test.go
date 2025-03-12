// test/ut/service/team_service_test.go
package service_test

import (
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/service"
	"github.com/vgeshiktor/nba-stats/test/ut/mocks"
)

func TestCreateTeam_Success(t *testing.T) {
	// Use FakeTeamRepo from our mocks package.
	fakeRepo := &mocks.FakeTeamRepo{}
	teamService := service.NewTeamService(fakeRepo)

	team := &domain.Team{
		ID:   "team1",
		Name: "Test Team",
	}

	err := teamService.CreateTeam(team)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestCreateTeam_Invalid(t *testing.T) {
	fakeRepo := &mocks.FakeTeamRepo{}
	teamService := service.NewTeamService(fakeRepo)

	// Test with missing Name.
	team := &domain.Team{
		ID:   "team1",
		Name: "",
	}

	err := teamService.CreateTeam(team)
	if err == nil {
		t.Errorf("expected error for missing team name, got nil")
	}

	// Test with missing ID.
	team = &domain.Team{
		ID:   "",
		Name: "Test Team",
	}
	err = teamService.CreateTeam(team)
	if err == nil {
		t.Errorf("expected error for missing team ID, got nil")
	}
}

func TestGetTeamByID_Success(t *testing.T) {
	fakeRepo := &mocks.FakeTeamRepo{}
	teamService := service.NewTeamService(fakeRepo)

	team, err := teamService.GetTeamByID("team1")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if team.ID != "team1" {
		t.Errorf("expected team ID 'team1', got: %s", team.ID)
	}
}

func TestGetTeamByID_Invalid(t *testing.T) {
	fakeRepo := &mocks.FakeTeamRepo{}
	teamService := service.NewTeamService(fakeRepo)

	_, err := teamService.GetTeamByID("invalid")
	if err == nil {
		t.Errorf("expected error for invalid team ID, got nil")
	}
}
