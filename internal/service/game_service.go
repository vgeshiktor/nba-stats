// internal/service/game_service.go
package service

import (
	"errors"
	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/repository"

	"github.com/vgeshiktor/nba-stats/pkg/logger"
	"github.com/vgeshiktor/nba-stats/pkg/validator"
)

// GameService defines the methods related to game management.
type GameService interface {
	CreateGame(game *domain.Game) error
	GetGameByID(id string) (*domain.Game, error)
}

type gameService struct {
	gameRepo repository.GameRepository
}

// NewGameService creates a new instance of GameService.
func NewGameService(gameRepo repository.GameRepository) GameService {
	return &gameService{gameRepo: gameRepo}
}

// CreateGame validates and inserts a new game into the database.
func (s *gameService) CreateGame(game *domain.Game) error {
	if err := validator.ValidateGame(game); err != nil {
		return err
	}
	logger.Info("creating game: %v", game)
	return s.gameRepo.CreateGame(game)
}

// GetGameByID fetches game details by ID.
func (s *gameService) GetGameByID(id string) (*domain.Game, error) {
	if id == "" {
		return nil, errors.New("game ID cannot be empty")
	}

	logger.Info("Get game by id: %v", id)
	return s.gameRepo.GetGameByID(id)
}
