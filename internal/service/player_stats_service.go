// internal/service/player_stats_service.go
package service

import (
	"errors"
	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/repository"

	"github.com/vgeshiktor/nba-stats/pkg/validator"
)

// PlayerStatsService defines operations for logging player game statistics.
type PlayerStatsService interface {
	LogPlayerStats(stats *domain.PlayerGameStats) error
}

type playerStatsService struct {
	playerRepo repository.PlayerRepository
	teamRepo   repository.TeamRepository
	gameRepo   repository.GameRepository
	statsRepo  repository.PlayerStatsRepository
}

// NewPlayerStatsService creates a new instance of PlayerStatsService.
func NewPlayerStatsService(playerRepo repository.PlayerRepository, teamRepo repository.TeamRepository, gameRepo repository.GameRepository, statsRepo repository.PlayerStatsRepository) PlayerStatsService {
	return &playerStatsService{
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
		gameRepo:   gameRepo,
		statsRepo:  statsRepo,
	}
}

// LogPlayerStats validates and stores player game statistics.
func (s *playerStatsService) LogPlayerStats(stats *domain.PlayerGameStats) error {
	if err := validator.ValidatePlayerStats(stats); err != nil {
		return err
	}

	// Ensure player exists
	_, err := s.playerRepo.GetPlayerByID(stats.PlayerID)
	if err != nil {
		return errors.New("player not found")
	}

	// Ensure game exists
	_, err = s.gameRepo.GetGameByID(stats.GameID)
	if err != nil {
		return errors.New("game not found")
	}

	// Store the stats
	return s.statsRepo.InsertPlayerStats(stats)
}
