// internal/repository/player_stats_repository.go
package repository

import (
	"database/sql"

	"github.com/vgeshiktor/nba-stats/internal/domain"
)

// PlayerStatsRepository defines operations for player game statistics.
type PlayerStatsRepository interface {
	InsertPlayerStats(stats *domain.PlayerGameStats) error
	FetchPlayerAggregate(playerID string) (*domain.AggregateStats, error)
	FetchTeamAggregate(teamID string) (*domain.AggregateStats, error)
}

type playerStatsRepo struct {
	db *sql.DB
}

// NewPlayerStatsRepository returns a new instance of PlayerStatsRepository.
func NewPlayerStatsRepository(db *sql.DB) PlayerStatsRepository {
	return &playerStatsRepo{db: db}
}

// InsertPlayerStats stores a player's game statistics.
func (r *playerStatsRepo) InsertPlayerStats(stats *domain.PlayerGameStats) error {
	query := `
		INSERT INTO player_game_stats 
		(id, player_id, game_id, points, rebounds, assists, steals, blocks, fouls, turnovers, minutes_played)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(query, stats.ID, stats.PlayerID, stats.GameID, stats.Points, stats.Rebounds,
		stats.Assists, stats.Steals, stats.Blocks, stats.Fouls, stats.Turnovers, stats.MinutesPlayed)
	return err
}

// FetchPlayerAggregate calculates and returns aggregated statistics for a player.
func (r *playerStatsRepo) FetchPlayerAggregate(playerID string) (*domain.AggregateStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT game_id) as games_played,
			SUM(points) as total_points,
			SUM(rebounds) as total_rebounds,
			SUM(assists) as total_assists,
			SUM(steals) as total_steals,
			SUM(blocks) as total_blocks,
			SUM(fouls) as total_fouls,
			SUM(turnovers) as total_turnovers,
			SUM(minutes_played) as total_minutes
		FROM player_game_stats
		WHERE player_id = $1
	`
	row := r.db.QueryRow(query, playerID)

	var agg domain.AggregateStats
	agg.PlayerID = playerID
	err := row.Scan(&agg.GamesPlayed, &agg.TotalPoints, &agg.TotalRebounds, &agg.TotalAssists,
		&agg.TotalSteals, &agg.TotalBlocks, &agg.TotalFouls, &agg.TotalTurnovers, &agg.TotalMinutes)
	if err != nil {
		return nil, err
	}
	if agg.GamesPlayed > 0 {
		agg.AvgPoints = float64(agg.TotalPoints) / float64(agg.GamesPlayed)
		agg.AvgRebounds = float64(agg.TotalRebounds) / float64(agg.GamesPlayed)
		agg.AvgAssists = float64(agg.TotalAssists) / float64(agg.GamesPlayed)
		agg.AvgSteals = float64(agg.TotalSteals) / float64(agg.GamesPlayed)
		agg.AvgBlocks = float64(agg.TotalBlocks) / float64(agg.GamesPlayed)
		agg.AvgFouls = float64(agg.TotalFouls) / float64(agg.GamesPlayed)
		agg.AvgTurnovers = float64(agg.TotalTurnovers) / float64(agg.GamesPlayed)
		agg.AvgMinutes = agg.TotalMinutes / float64(agg.GamesPlayed)
	}
	return &agg, nil
}

// FetchTeamAggregate calculates and returns aggregated statistics for a team by joining player data.
func (r *playerStatsRepo) FetchTeamAggregate(teamID string) (*domain.AggregateStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT game_id) as games_played,
			SUM(ps.points) as total_points,
			SUM(ps.rebounds) as total_rebounds,
			SUM(ps.assists) as total_assists,
			SUM(ps.steals) as total_steals,
			SUM(ps.blocks) as total_blocks,
			SUM(ps.fouls) as total_fouls,
			SUM(ps.turnovers) as total_turnovers,
			SUM(ps.minutes_played) as total_minutes
		FROM player_game_stats ps
		INNER JOIN players p ON ps.player_id = p.id
		WHERE p.team_id = $1
	`
	row := r.db.QueryRow(query, teamID)

	var agg domain.AggregateStats
	agg.TeamID = teamID
	err := row.Scan(&agg.GamesPlayed, &agg.TotalPoints, &agg.TotalRebounds, &agg.TotalAssists,
		&agg.TotalSteals, &agg.TotalBlocks, &agg.TotalFouls, &agg.TotalTurnovers, &agg.TotalMinutes)
	if err != nil {
		return nil, err
	}
	if agg.GamesPlayed > 0 {
		agg.AvgPoints = float64(agg.TotalPoints) / float64(agg.GamesPlayed)
		agg.AvgRebounds = float64(agg.TotalRebounds) / float64(agg.GamesPlayed)
		agg.AvgAssists = float64(agg.TotalAssists) / float64(agg.GamesPlayed)
		agg.AvgSteals = float64(agg.TotalSteals) / float64(agg.GamesPlayed)
		agg.AvgBlocks = float64(agg.TotalBlocks) / float64(agg.GamesPlayed)
		agg.AvgFouls = float64(agg.TotalFouls) / float64(agg.GamesPlayed)
		agg.AvgTurnovers = float64(agg.TotalTurnovers) / float64(agg.GamesPlayed)
		agg.AvgMinutes = agg.TotalMinutes / float64(agg.GamesPlayed)
	}
	return &agg, nil
}
