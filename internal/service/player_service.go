// internal/service/player_service.go
package service

import (
	"errors"
	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/repository"

	"github.com/vgeshiktor/nba-stats/pkg/validator"
)

// PlayerService defines the methods related to player management.
type PlayerService interface {
	CreatePlayer(player *domain.Player) error
	GetPlayerByID(id string) (*domain.Player, error)
}

type playerService struct {
	playerRepo repository.PlayerRepository
}

// NewPlayerService creates a new instance of PlayerService.
func NewPlayerService(playerRepo repository.PlayerRepository) PlayerService {
	return &playerService{playerRepo: playerRepo}
}

// CreatePlayer validates and inserts a new player into the database.
func (s *playerService) CreatePlayer(player *domain.Player) error {
	if err := validator.ValidatePlayer(player); err != nil {
		return err
	}
	return s.playerRepo.CreatePlayer(player)
}

// GetPlayerByID fetches player details by ID.
func (s *playerService) GetPlayerByID(id string) (*domain.Player, error) {
	if id == "" {
		return nil, errors.New("player ID cannot be empty")
	}
	return s.playerRepo.GetPlayerByID(id)
}
