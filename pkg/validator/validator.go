// pkg/validator/validator.go
package validator

import (
	"errors"
	"github.com/vgeshiktor/nba-stats/internal/domain"
)

// ValidatePlayer ensures a player's data is valid.
func ValidatePlayer(player *domain.Player) error {
	if player.ID == "" || player.Name == "" || player.TeamID == "" {
		return errors.New("player ID, name, and team ID cannot be empty")
	}
	return nil
}

// ValidateTeam ensures a team's data is valid.
func ValidateTeam(team *domain.Team) error {
	if team.ID == "" || team.Name == "" {
		return errors.New("team ID and name cannot be empty")
	}
	return nil
}

// ValidateGame ensures a game's data is valid.
func ValidateGame(game *domain.Game) error {
	if game.ID == "" || game.HomeTeam == "" || game.AwayTeam == "" {
		return errors.New("game ID, home team, and away team cannot be empty")
	}
	return nil
}

// ValidatePlayerStats ensures player statistics are valid.
func ValidatePlayerStats(stats *domain.PlayerGameStats) error {
	if stats.PlayerID == "" || stats.GameID == "" {
		return errors.New("player ID and game ID cannot be empty")
	}
	if stats.Fouls > 6 {
		return errors.New("fouls cannot exceed 6")
	}
	if stats.MinutesPlayed < 0 || stats.MinutesPlayed > 48 {
		return errors.New("minutes played must be between 0 and 48")
	}
	return nil
}
