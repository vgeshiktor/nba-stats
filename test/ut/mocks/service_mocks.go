// test/ut/service/service_mocks.go
package mocks

import (
	"github.com/vgeshiktor/nba-stats/internal/domain"
)

// --- Fake Service Implementations for API Testing ---

type FakePlayerStatsService struct{}

func (s *FakePlayerStatsService) LogPlayerStats(stats *domain.PlayerGameStats) error {
	return nil
}

type FakeAggregationService struct{}

func (s *FakeAggregationService) GetPlayerAggregate(playerID string) (*domain.AggregateStats, error) {
	return &domain.AggregateStats{
		PlayerID:    playerID,
		GamesPlayed: 1,
		TotalPoints: 30,
		AvgPoints:   30,
	}, nil
}

func (s *FakeAggregationService) GetTeamAggregate(teamID string) (*domain.AggregateStats, error) {
	return &domain.AggregateStats{
		TeamID:      teamID,
		GamesPlayed: 1,
		TotalPoints: 100,
		AvgPoints:   100,
	}, nil
}

type FakePlayerService struct{}

func (s *FakePlayerService) CreatePlayer(player *domain.Player) error { return nil }
func (s *FakePlayerService) GetPlayerByID(id string) (*domain.Player, error) {
	return &domain.Player{ID: id, Name: "Test Player", TeamID: "team1"}, nil
}

type FakeTeamService struct{}

func (s *FakeTeamService) CreateTeam(team *domain.Team) error { return nil }
func (s *FakeTeamService) GetTeamByID(id string) (*domain.Team, error) {
	return &domain.Team{ID: id, Name: "Test Team"}, nil
}

type FakeGameService struct{}

func (s *FakeGameService) CreateGame(game *domain.Game) error { return nil }
func (s *FakeGameService) GetGameByID(id string) (*domain.Game, error) {
	return &domain.Game{ID: id, HomeTeam: "team1", AwayTeam: "team2"}, nil
}