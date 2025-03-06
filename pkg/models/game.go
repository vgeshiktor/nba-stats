package models

import (
	"errors"
	"time"
)

// GameStatus represents the current state of a game
type GameStatus string

const (
	GameStatusScheduled GameStatus = "SCHEDULED"
	GameStatusInProgress GameStatus = "IN_PROGRESS"
	GameStatusCompleted  GameStatus = "COMPLETED"
	GameStatusCancelled  GameStatus = "CANCELLED"
)

// Game represents an NBA basketball game between two teams
type Game struct {
	ID          int        `json:"id" db:"id"`
	HomeTeamID  int        `json:"home_team_id" db:"home_team_id"`
	AwayTeamID  int        `json:"away_team_id" db:"away_team_id"`
	Date        time.Time  `json:"date" db:"game_date"`
	Season      string     `json:"season" db:"season"` // Format: "2024-2025"
	Venue       string     `json:"venue" db:"venue"`
	Status      GameStatus `json:"status" db:"status"`
	HomeScore   *int       `json:"home_score,omitempty" db:"home_score"`
	AwayScore   *int       `json:"away_score,omitempty" db:"away_score"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// NewGame creates a new Game instance with basic validation
func NewGame(homeTeamID, awayTeamID int, date time.Time, season, venue string) (*Game, error) {
	if homeTeamID <= 0 {
		return nil, errors.New("home team ID must be positive")
	}
	if awayTeamID <= 0 {
		return nil, errors.New("away team ID must be positive")
	}
	if homeTeamID == awayTeamID {
		return nil, errors.New("home team and away team cannot be the same")
	}
	if date.IsZero() {
		return nil, errors.New("game date cannot be empty")
	}
	if season == "" {
		return nil, errors.New("season cannot be empty")
	}

	now := time.Now()
	return &Game{
		HomeTeamID: homeTeamID,
		AwayTeamID: awayTeamID,
		Date:       date,
		Season:     season,
		Venue:      venue,
		Status:     GameStatusScheduled,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// UpdateScore updates the game score and status
func (g *Game) UpdateScore(homeScore, awayScore int) error {
	if homeScore < 0 {
		return errors.New("home score cannot be negative")
	}
	if awayScore < 0 {
		return errors.New("away score cannot be negative")
	}

	g.HomeScore = &homeScore
	g.AwayScore = &awayScore
	g.UpdatedAt = time.Now()
	return nil
}

// MarkAsInProgress marks the game as in progress
func (g *Game) MarkAsInProgress() error {
	if g.Status == GameStatusCancelled {
		return errors.New("cannot start a cancelled game")
	}
	if g.Status == GameStatusCompleted {
		return errors.New("cannot restart a completed game")
	}

	g.Status = GameStatusInProgress
	g.UpdatedAt = time.Now()
	return nil
}

// MarkAsCompleted marks the game as completed
func (g *Game) MarkAsCompleted() error {
	if g.Status == GameStatusCancelled {
		return errors.New("cannot complete a cancelled game")
	}
	if g.Status != GameStatusInProgress {
		return errors.New("only in-progress games can be completed")
	}
	if g.HomeScore == nil || g.AwayScore == nil {
		return errors.New("game score must be set before marking as completed")
	}

	g.Status = GameStatusCompleted
	g.UpdatedAt = time.Now()
	return nil
}

// MarkAsCancelled marks the game as cancelled
func (g *Game) MarkAsCancelled() error {
	if g.Status == GameStatusCompleted {
		return errors.New("cannot cancel a completed game")
	}

	g.Status = GameStatusCancelled
	g.UpdatedAt = time.Now()
	return nil
}

// GetWinningTeamID returns the ID of the winning team, or 0 if game isn't completed or is a tie
func (g *Game) GetWinningTeamID() int {
	if g.Status != GameStatusCompleted || g.HomeScore == nil || g.AwayScore == nil {
		return 0
	}

	if *g.HomeScore > *g.AwayScore {
		return g.HomeTeamID
	} else if *g.AwayScore > *g.HomeScore {
		return g.AwayTeamID
	}
	
	return 0 // Tie game
}
