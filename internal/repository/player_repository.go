// internal/repository/player_repository.go
package repository

import (
	"database/sql"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/pkg/logger"
)

// PlayerRepository defines operations on Player data.
type PlayerRepository interface {
	CreatePlayer(player *domain.Player) error
	GetPlayerByID(id string) (*domain.Player, error)
}

type playerRepo struct {
	db *sql.DB
}

// NewPlayerRepository returns a new instance of PlayerRepository.
func NewPlayerRepository(db *sql.DB) PlayerRepository {
	return &playerRepo{db: db}
}

// CreatePlayer inserts a new player record into the database.
func (r *playerRepo) CreatePlayer(player *domain.Player) error {
	query := `INSERT INTO players (id, name, team_id) VALUES ($1, $2, $3)`
	logger.Info("Running query: %s", query)
	_, err := r.db.Exec(query, player.ID, player.Name, player.TeamID)
	return err
}

// GetPlayerByID retrieves a player by its ID.
func (r *playerRepo) GetPlayerByID(id string) (*domain.Player, error) {
	query := `SELECT id, name, team_id FROM players WHERE id = $1`
	logger.Info("Running query: %s", query)
	row := r.db.QueryRow(query, id)
	var player domain.Player
	if err := row.Scan(&player.ID, &player.Name, &player.TeamID); err != nil {
		return nil, err
	}
	return &player, nil
}
