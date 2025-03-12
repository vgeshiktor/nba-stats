// internal/repository/team_repository.go
package repository

import (
	"database/sql"
	"github.com/vgeshiktor/nba-stats/internal/domain"
)

// TeamRepository defines operations on Team data.
type TeamRepository interface {
	CreateTeam(team *domain.Team) error
	GetTeamByID(id string) (*domain.Team, error)
}

type teamRepo struct {
	db *sql.DB
}

// NewTeamRepository returns a new instance of TeamRepository.
func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepo{db: db}
}

// CreateTeam inserts a new team record into the database.
func (r *teamRepo) CreateTeam(team *domain.Team) error {
	query := `INSERT INTO teams (id, name) VALUES ($1, $2)`
	_, err := r.db.Exec(query, team.ID, team.Name)
	return err
}

// GetTeamByID retrieves a team by its ID.
func (r *teamRepo) GetTeamByID(id string) (*domain.Team, error) {
	query := `SELECT id, name FROM teams WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var team domain.Team
	if err := row.Scan(&team.ID, &team.Name); err != nil {
		return nil, err
	}
	return &team, nil
}
