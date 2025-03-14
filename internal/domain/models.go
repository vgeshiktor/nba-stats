// internal/domain/models.go
package domain

import (
	"time"
)

// Player represents an NBA player.
type Player struct {
	ID     string `json:"id"`      // Unique identifier for the player.
	Name   string `json:"name"`    // Player's full name.
	TeamID string `json:"team_id"` // Associated team identifier.
}

// Team represents an NBA team.
type Team struct {
	ID   string `json:"id"`   // Unique identifier for the team.
	Name string `json:"name"` // Name of the team.
}

// Game represents a single NBA game.
type Game struct {
	ID       string    `json:"id"`        // Unique identifier for the game.
	Date     time.Time `json:"date"`      // Date and time of the game.
	HomeTeam string    `json:"home_team"` // Home team identifier.
	AwayTeam string    `json:"away_team"` // Away team identifier.
}

// PlayerGameStats holds the statistics for a player in a specific game.
type 	PlayerGameStats struct {
	ID            string  `json:"id,omitempty"` // Unique identifier for the stats record (optional).
	PlayerID      string  `json:"player_id"`    // Identifier of the player.
	GameID        string  `json:"game_id"`      // Identifier of the game.
	Points        int     `json:"points"`       // Points scored.
	Rebounds      int     `json:"rebounds"`     // Rebounds recorded.
	Assists       int     `json:"assists"`      // Assists made.
	Steals        int     `json:"steals"`       // Steals recorded.
	Blocks        int     `json:"blocks"`       // Blocks recorded.
	Fouls         int     `json:"fouls"`        // Fouls committed (maximum allowed value: 6).
	Turnovers     int     `json:"turnovers"`    // Turnovers committed.
	MinutesPlayed float64 `json:"minutes_played"` // Minutes played in the game (range: 0 to 48.0).
}

// AggregateStats represents aggregated season statistics for a player or team.
type AggregateStats struct {
	// Either PlayerID or TeamID will be set.
	PlayerID      string  `json:"player_id,omitempty"`
	TeamID        string  `json:"team_id,omitempty"`
	GamesPlayed   int     `json:"games_played"`
	TotalPoints   int     `json:"total_points"`
	TotalRebounds int     `json:"total_rebounds"`
	TotalAssists  int     `json:"total_assists"`
	TotalSteals   int     `json:"total_steals"`
	TotalBlocks   int     `json:"total_blocks"`
	TotalFouls    int     `json:"total_fouls"`
	TotalTurnovers int    `json:"total_turnovers"`
	TotalMinutes  float64 `json:"total_minutes"`
	AvgPoints     float64 `json:"avg_points"`
	AvgRebounds   float64 `json:"avg_rebounds"`
	AvgAssists    float64 `json:"avg_assists"`
	AvgSteals     float64 `json:"avg_steals"`
	AvgBlocks     float64 `json:"avg_blocks"`
	AvgFouls      float64 `json:"avg_fouls"`
	AvgTurnovers  float64 `json:"avg_turnovers"`
	AvgMinutes    float64 `json:"avg_minutes"`
}
