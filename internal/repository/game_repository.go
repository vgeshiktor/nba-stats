// internal/repository/game_repository.go
package repository

import (
	"database/sql"
	"github.com/vgeshiktor/nba-stats/internal/domain"
)

// GameRepository defines operations on Game data.
type GameRepository interface {
	CreateGame(game *domain.Game) error
	GetGameByID(id string) (*domain.Game, error)
}

type gameRepo struct {
	db *sql.DB
}

// NewGameRepository returns a new instance of GameRepository.
func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepo{db: db}
}

// CreateGame inserts a new game record into the database.
func (r *gameRepo) CreateGame(game *domain.Game) error {
	query := `INSERT INTO games (id, date, home_team, away_team) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, game.ID, game.Date, game.HomeTeam, game.AwayTeam)
	return err
}

// GetGameByID retrieves a game by its ID.
func (r *gameRepo) GetGameByID(id string) (*domain.Game, error) {
	query := `SELECT id, date, home_team, away_team FROM games WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var game domain.Game
	if err := row.Scan(&game.ID, &game.Date, &game.HomeTeam, &game.AwayTeam); err != nil {
		return nil, err
	}
	return &game, nil
}
