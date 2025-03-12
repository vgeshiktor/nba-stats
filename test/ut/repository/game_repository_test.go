// test/ut/repository/game_repository_test.go
package repository_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateGame_Success(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	// Initialize the game repository using the mocked database.
	repo := repository.NewGameRepository(db)

	// Prepare a sample game.
	game := &domain.Game{
		ID:       "game1",
		Date:     time.Now(),
		HomeTeam: "team1",
		AwayTeam: "team2",
	}

	// Expect an INSERT statement.
	mock.ExpectExec("INSERT INTO games").
		WithArgs(game.ID, game.Date, game.HomeTeam, game.AwayTeam).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call CreateGame.
	err = repo.CreateGame(game)
	if err != nil {
		t.Errorf("unexpected error on CreateGame: %s", err)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetGameByID_Success(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := repository.NewGameRepository(db)

	// Define expected game data.
	gameID := "game1"
	gameDate := time.Now()
	homeTeam := "team1"
	awayTeam := "team2"

	// Set up expected query and result rows.
	rows := sqlmock.NewRows([]string{"id", "date", "home_team", "away_team"}).
		AddRow(gameID, gameDate, homeTeam, awayTeam)
	mock.ExpectQuery("SELECT id, date, home_team, away_team FROM games WHERE id = \\$1").
		WithArgs(gameID).
		WillReturnRows(rows)

	// Call GetGameByID.
	game, err := repo.GetGameByID(gameID)
	if err != nil {
		t.Errorf("unexpected error on GetGameByID: %s", err)
	}

	// Verify returned game data.
	if game.ID != gameID {
		t.Errorf("expected game ID %s, got %s", gameID, game.ID)
	}
	if game.HomeTeam != homeTeam {
		t.Errorf("expected home team %s, got %s", homeTeam, game.HomeTeam)
	}
	if game.AwayTeam != awayTeam {
		t.Errorf("expected away team %s, got %s", awayTeam, game.AwayTeam)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetGameByID_NotFound(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := repository.NewGameRepository(db)

	gameID := "nonexistent"

	// Set up expected query returning no rows.
	mock.ExpectQuery("SELECT id, date, home_team, away_team FROM games WHERE id = \\$1").
		WithArgs(gameID).
		WillReturnError(sql.ErrNoRows)

	// Call GetGameByID.
	_, err = repo.GetGameByID(gameID)
	if err == nil {
		t.Error("expected error when game is not found, got nil")
	}
	// Optionally, you can check for a specific error:
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows error, got: %v", err)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
