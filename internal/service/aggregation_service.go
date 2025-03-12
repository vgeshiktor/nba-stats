// internal/service/aggregation_service.go
package service

import (
	"errors"
	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/repository"
)

// AggregationService defines operations for retrieving aggregate statistics.
type AggregationService interface {
	GetPlayerAggregate(playerID string) (*domain.AggregateStats, error)
	GetTeamAggregate(teamID string) (*domain.AggregateStats, error)
}

type aggregationService struct {
	statsRepo repository.PlayerStatsRepository
}

// NewAggregationService creates a new instance of AggregationService.
func NewAggregationService(statsRepo repository.PlayerStatsRepository) AggregationService {
	return &aggregationService{statsRepo: statsRepo}
}

// GetPlayerAggregate retrieves the season averages for a specific player.
func (s *aggregationService) GetPlayerAggregate(playerID string) (*domain.AggregateStats, error) {
	if playerID == "" {
		return nil, errors.New("player ID cannot be empty")
	}
	return s.statsRepo.FetchPlayerAggregate(playerID)
}

// GetTeamAggregate retrieves the season averages for a specific team.
func (s *aggregationService) GetTeamAggregate(teamID string) (*domain.AggregateStats, error) {
	if teamID == "" {
		return nil, errors.New("team ID cannot be empty")
	}
	return s.statsRepo.FetchTeamAggregate(teamID)
}
