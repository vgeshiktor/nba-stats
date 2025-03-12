// test/ut/repository/player_repository_test.go
package repository_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreatePlayer_Success(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	// Initialize the player repository using the mocked database.
	repo := repository.NewPlayerRepository(db)

	// Prepare a sample player.
	player := &domain.Player{
		ID:     "player1",
		Name:   "John Doe",
		TeamID: "team1",
	}

	// Expect an INSERT statement.
	mock.ExpectExec("INSERT INTO players").
		WithArgs(player.ID, player.Name, player.TeamID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call CreatePlayer.
	err = repo.CreatePlayer(player)
	if err != nil {
		t.Errorf("unexpected error on CreatePlayer: %s", err)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPlayerByID_Success(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := repository.NewPlayerRepository(db)

	// Define expected player data.
	playerID := "player1"
	expectedName := "John Doe"
	expectedTeamID := "team1"

	// Set up expected query and result rows.
	rows := sqlmock.NewRows([]string{"id", "name", "team_id"}).
		AddRow(playerID, expectedName, expectedTeamID)
	mock.ExpectQuery("SELECT id, name, team_id FROM players WHERE id = \\$1").
		WithArgs(playerID).
		WillReturnRows(rows)

	// Call GetPlayerByID.
	player, err := repo.GetPlayerByID(playerID)
	if err != nil {
		t.Errorf("unexpected error on GetPlayerByID: %s", err)
	}

	// Verify returned player data.
	if player.ID != playerID {
		t.Errorf("expected player id %s, got %s", playerID, player.ID)
	}
	if player.Name != expectedName {
		t.Errorf("expected player name %s, got %s", expectedName, player.Name)
	}
	if player.TeamID != expectedTeamID {
		t.Errorf("expected team id %s, got %s", expectedTeamID, player.TeamID)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPlayerByID_NotFound(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	repo := repository.NewPlayerRepository(db)

	playerID := "nonexistent"

	// Set up expected query returning no rows.
	mock.ExpectQuery("SELECT id, name, team_id FROM players WHERE id = \\$1").
		WithArgs(playerID).
		WillReturnError(sql.ErrNoRows)

	// Call GetPlayerByID.
	_, err = repo.GetPlayerByID(playerID)
	if err == nil {
		t.Error("expected error when player not found, got nil")
	}

	// Optionally, verify the error matches sql.ErrNoRows.
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected error sql.ErrNoRows, got: %v", err)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
