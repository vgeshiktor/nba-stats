// test/ut/repository/team_repository_test.go
package repository_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTeam_Success(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	// Initialize the team repository using the mocked database.
	repo := repository.NewTeamRepository(db)

	// Prepare a sample team.
	team := &domain.Team{
		ID:   "team1",
		Name: "Test Team",
	}

	// Expect an INSERT statement.
	mock.ExpectExec("INSERT INTO teams").
		WithArgs(team.ID, team.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call CreateTeam.
	err = repo.CreateTeam(team)
	if err != nil {
		t.Errorf("unexpected error on CreateTeam: %v", err)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetTeamByID_Success(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := repository.NewTeamRepository(db)

	// Define expected team data.
	teamID := "team1"
	expectedName := "Test Team"

	// Set up expected query and result rows.
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(teamID, expectedName)
	mock.ExpectQuery("SELECT id, name FROM teams WHERE id = \\$1").
		WithArgs(teamID).
		WillReturnRows(rows)

	// Call GetTeamByID.
	team, err := repo.GetTeamByID(teamID)
	if err != nil {
		t.Errorf("unexpected error on GetTeamByID: %v", err)
	}

	// Verify returned team data.
	if team.ID != teamID {
		t.Errorf("expected team id %s, got %s", teamID, team.ID)
	}
	if team.Name != expectedName {
		t.Errorf("expected team name %s, got %s", expectedName, team.Name)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetTeamByID_NotFound(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := repository.NewTeamRepository(db)

	teamID := "nonexistent"

	// Set up expected query returning an error (simulate no rows found).
	mock.ExpectQuery("SELECT id, name FROM teams WHERE id = \\$1").
		WithArgs(teamID).
		WillReturnError(sql.ErrNoRows)

	// Call GetTeamByID.
	_, err = repo.GetTeamByID(teamID)
	if err == nil {
		t.Error("expected error when team not found, got nil")
	}

	// Optionally, verify that the error is sql.ErrNoRows.
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected error sql.ErrNoRows, got: %v", err)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
