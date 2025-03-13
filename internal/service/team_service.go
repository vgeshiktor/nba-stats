// internal/service/team_service.go
package service

import (
	"errors"
	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/repository"

	"github.com/vgeshiktor/nba-stats/pkg/validator"
	"github.com/vgeshiktor/nba-stats/pkg/logger"
)

// TeamService defines the methods related to team management.
type TeamService interface {
	CreateTeam(team *domain.Team) error
	GetTeamByID(id string) (*domain.Team, error)
}

type teamService struct {
	teamRepo repository.TeamRepository
}

// NewTeamService creates a new instance of TeamService.
func NewTeamService(teamRepo repository.TeamRepository) TeamService {
	return &teamService{teamRepo: teamRepo}
}

// CreateTeam validates and inserts a new team into the database.
func (s *teamService) CreateTeam(team *domain.Team) error {
	if err := validator.ValidateTeam(team); err != nil {
		return err
	}

	return s.teamRepo.CreateTeam(team)
}

// GetTeamByID fetches team details by ID.
func (s *teamService) GetTeamByID(id string) (*domain.Team, error) {
	if id == "" {
		return nil, errors.New("team ID cannot be empty")
	}
	logger.Info("Getting team by id: %s", id)

	return s.teamRepo.GetTeamByID(id)
}
