// test/ut/service/repo_mocks.go
package mocks

import (
	"errors"

	"github.com/vgeshiktor/nba-stats/internal/domain"
)

// -------------------------
// Fake Player Repository
// -------------------------

// FakePlayerRepo implements the repository.PlayerRepository interface.
type FakePlayerRepo struct{}

func (r *FakePlayerRepo) CreatePlayer(player *domain.Player) error {
	if player.ID == "" {
		return errors.New("player ID cannot be empty")
	}
	return nil
}

func (r *FakePlayerRepo) GetPlayerByID(id string) (*domain.Player, error) {
	if id == "valid" {
		return &domain.Player{
			ID:     "valid",
			Name:   "Test Player",
			TeamID: "team1",
		}, nil
	}
	return nil, errors.New("player not found")
}

// -------------------------
// Fake Team Repository
// -------------------------

// FakeTeamRepo implements the repository.TeamRepository interface.
type FakeTeamRepo struct{}

func (r *FakeTeamRepo) CreateTeam(team *domain.Team) error {
	if team.ID == "" {
		return errors.New("team ID cannot be empty")
	}
	return nil
}

func (r *FakeTeamRepo) GetTeamByID(id string) (*domain.Team, error) {
	if id == "team1" {
		return &domain.Team{
			ID:   "team1",
			Name: "Test Team",
		}, nil
	}
	return nil, errors.New("team not found")
}

// -------------------------
// Fake Game Repository
// -------------------------

// FakeGameRepo implements the repository.GameRepository interface.
type FakeGameRepo struct{}

func (r *FakeGameRepo) CreateGame(game *domain.Game) error {
	if game.ID == "" {
		return errors.New("game ID cannot be empty")
	}
	return nil
}

func (r *FakeGameRepo) GetGameByID(id string) (*domain.Game, error) {
	if id == "game1" {
		return &domain.Game{
			ID:       "game1",
			HomeTeam: "team1",
			AwayTeam: "team2",
		}, nil
	}
	return nil, errors.New("game not found")
}

// -------------------------
// Fake Player Stats Repository
// -------------------------

// FakePlayerStatsRepo implements the repository.PlayerStatsRepository interface.
type FakePlayerStatsRepo struct {
	Inserted bool
}

func (r *FakePlayerStatsRepo) InsertPlayerStats(stats *domain.PlayerGameStats) error {
	if stats.PlayerID == "" || stats.GameID == "" {
		return errors.New("invalid stats: missing playerID or gameID")
	}
	r.Inserted = true
	return nil
}

func (r *FakePlayerStatsRepo) FetchPlayerAggregate(playerID string) (*domain.AggregateStats, error) {
	if playerID == "valid" {
		return &domain.AggregateStats{
			PlayerID:    "valid",
			GamesPlayed: 1,
			TotalPoints: 30,
			AvgPoints:   30,
		}, nil
	}
	return nil, errors.New("aggregate not found")
}

func (r *FakePlayerStatsRepo) FetchTeamAggregate(teamID string) (*domain.AggregateStats, error) {
	if teamID == "team1" {
		return &domain.AggregateStats{
			TeamID:      "team1",
			GamesPlayed: 1,
			TotalPoints: 100,
			AvgPoints:   100,
		}, nil
	}
	return nil, errors.New("aggregate not found")
}
