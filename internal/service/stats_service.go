package service

import (
	"fmt"

	models "github.com/vgeshiktor/nba-stats/pkg/models"
)

// StatsService implements business logic for player statistics
type StatsService struct {
	playerRepo PlayerRepository
	teamRepo   TeamRepository
	statsRepo  GameStatsRepository
}

func NewStatsService(
	playerRepo PlayerRepository,
	teamRepo TeamRepository,
	statsRepo GameStatsRepository,
) *StatsService {
	return &StatsService{
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
		statsRepo:  statsRepo,
	}
}

func (s *StatsService) LogPlayerStats(stats *PlayerStatsInput) error {
	// Verify player exists
	player, err := s.playerRepo.GetByID(stats.PlayerID)
	if err != nil {
		return fmt.Errorf("invalid player ID: %w", err)
	}
	
	// Additional validation or business logic would go here
	
	// Save stats
	return s.statsRepo.CreatePlayerStats(stats)
}

func (s *StatsService) GetPlayerSeasonAverage(playerID int) (*SeasonAverage, error) {
	// Verify player exists
	_, err := s.playerRepo.GetByID(playerID)
	if err != nil {
		return nil, fmt.Errorf("invalid player ID: %w", err)
	}
	
	return s.statsRepo.GetPlayerSeasonStats(playerID)
}
