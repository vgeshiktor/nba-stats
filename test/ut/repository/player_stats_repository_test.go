// test/ut/repository/player_stats_repository_test.go
package repository_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertPlayerStats_Success(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := repository.NewPlayerStatsRepository(db)

	// Prepare a sample PlayerGameStats record.
	stats := &domain.PlayerGameStats{
		ID:            "stats1",
		PlayerID:      "player1",
		GameID:        "game1",
		Points:        25,
		Rebounds:      8,
		Assists:       5,
		Steals:        2,
		Blocks:        1,
		Fouls:         3,
		Turnovers:     2,
		MinutesPlayed: 35.5,
	}

	// Expect an INSERT statement.
	mock.ExpectExec("INSERT INTO player_game_stats").
		WithArgs(stats.ID, stats.PlayerID, stats.GameID, stats.Points, stats.Rebounds,
			stats.Assists, stats.Steals, stats.Blocks, stats.Fouls, stats.Turnovers, stats.MinutesPlayed).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call InsertPlayerStats.
	err = repo.InsertPlayerStats(stats)
	if err != nil {
		t.Errorf("unexpected error on InsertPlayerStats: %v", err)
	}

	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestFetchPlayerAggregate_Success(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := repository.NewPlayerStatsRepository(db)

	// Define expected aggregate values.
	playerID := "player1"
	gamesPlayed := 3
	totalPoints := 90
	totalRebounds := 24
	totalAssists := 15
	totalSteals := 6
	totalBlocks := 3
	totalFouls := 9
	totalTurnovers := 6
	totalMinutes := 105.0

	// Set up expected query and result rows.
	rows := sqlmock.NewRows([]string{
		"games_played",
		"total_points",
		"total_rebounds",
		"total_assists",
		"total_steals",
		"total_blocks",
		"total_fouls",
		"total_turnovers",
		"total_minutes",
	}).AddRow(
		gamesPlayed,
		totalPoints,
		totalRebounds,
		totalAssists,
		totalSteals,
		totalBlocks,
		totalFouls,
		totalTurnovers,
		totalMinutes,
	)

	mock.ExpectQuery("SELECT (.+) FROM player_game_stats WHERE player_id = \\$1").
		WithArgs(playerID).
		WillReturnRows(rows)

	// Call FetchPlayerAggregate.
	agg, err := repo.FetchPlayerAggregate(playerID)
	if err != nil {
		t.Errorf("unexpected error on FetchPlayerAggregate: %v", err)
	}

	// Verify returned aggregate values.
	if agg.GamesPlayed != gamesPlayed {
		t.Errorf("expected games_played %d, got %d", gamesPlayed, agg.GamesPlayed)
	}
	if agg.TotalPoints != totalPoints {
		t.Errorf("expected total_points %d, got %d", totalPoints, agg.TotalPoints)
	}
	if agg.TotalRebounds != totalRebounds {
		t.Errorf("expected total_rebounds %d, got %d", totalRebounds, agg.TotalRebounds)
	}
	if agg.TotalAssists != totalAssists {
		t.Errorf("expected total_assists %d, got %d", totalAssists, agg.TotalAssists)
	}
	if agg.TotalSteals != totalSteals {
		t.Errorf("expected total_steals %d, got %d", totalSteals, agg.TotalSteals)
	}
	if agg.TotalBlocks != totalBlocks {
		t.Errorf("expected total_blocks %d, got %d", totalBlocks, agg.TotalBlocks)
	}
	if agg.TotalFouls != totalFouls {
		t.Errorf("expected total_fouls %d, got %d", totalFouls, agg.TotalFouls)
	}
	if agg.TotalTurnovers != totalTurnovers {
		t.Errorf("expected total_turnovers %d, got %d", totalTurnovers, agg.TotalTurnovers)
	}
	if agg.TotalMinutes != totalMinutes {
		t.Errorf("expected total_minutes %f, got %f", totalMinutes, agg.TotalMinutes)
	}

	// Averages are calculated based on gamesPlayed.
	expectedAvgPoints := float64(totalPoints) / float64(gamesPlayed)
	if agg.AvgPoints != expectedAvgPoints {
		t.Errorf("expected avg_points %f, got %f", expectedAvgPoints, agg.AvgPoints)
	}

	// Ensure all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestFetchPlayerAggregate_Error(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := repository.NewPlayerStatsRepository(db)

	playerID := "nonexistent"

	// Set up expected query returning an error.
	mock.ExpectQuery("SELECT (.+) FROM player_game_stats WHERE player_id = \\$1").
		WithArgs(playerID).
		WillReturnError(sql.ErrNoRows)

	// Call FetchPlayerAggregate.
	_, err = repo.FetchPlayerAggregate(playerID)
	if err == nil {
		t.Error("expected error when aggregate not found, got nil")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows error, got: %v", err)
	}

	// Ensure all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestFetchTeamAggregate_Success(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := repository.NewPlayerStatsRepository(db)

	// Define expected aggregate values for team.
	teamID := "team1"
	gamesPlayed := 4
	totalPoints := 120
	totalRebounds := 32
	totalAssists := 20
	totalSteals := 8
	totalBlocks := 4
	totalFouls := 10
	totalTurnovers := 8
	totalMinutes := 140.0

	// Set up expected query and result rows.
	rows := sqlmock.NewRows([]string{
		"games_played",
		"total_points",
		"total_rebounds",
		"total_assists",
		"total_steals",
		"total_blocks",
		"total_fouls",
		"total_turnovers",
		"total_minutes",
	}).AddRow(
		gamesPlayed,
		totalPoints,
		totalRebounds,
		totalAssists,
		totalSteals,
		totalBlocks,
		totalFouls,
		totalTurnovers,
		totalMinutes,
	)

	mock.ExpectQuery("SELECT (.+) FROM player_game_stats ps INNER JOIN players p ON ps.player_id = p.id WHERE p.team_id = \\$1").
		WithArgs(teamID).
		WillReturnRows(rows)

	// Call FetchTeamAggregate.
	agg, err := repo.FetchTeamAggregate(teamID)
	if err != nil {
		t.Errorf("unexpected error on FetchTeamAggregate: %v", err)
	}

	// Verify returned aggregate values.
	if agg.GamesPlayed != gamesPlayed {
		t.Errorf("expected games_played %d, got %d", gamesPlayed, agg.GamesPlayed)
	}
	if agg.TotalPoints != totalPoints {
		t.Errorf("expected total_points %d, got %d", totalPoints, agg.TotalPoints)
	}
	if agg.TotalRebounds != totalRebounds {
		t.Errorf("expected total_rebounds %d, got %d", totalRebounds, agg.TotalRebounds)
	}
	if agg.TotalAssists != totalAssists {
		t.Errorf("expected total_assists %d, got %d", totalAssists, agg.TotalAssists)
	}
	if agg.TotalSteals != totalSteals {
		t.Errorf("expected total_steals %d, got %d", totalSteals, agg.TotalSteals)
	}
	if agg.TotalBlocks != totalBlocks {
		t.Errorf("expected total_blocks %d, got %d", totalBlocks, agg.TotalBlocks)
	}
	if agg.TotalFouls != totalFouls {
		t.Errorf("expected total_fouls %d, got %d", totalFouls, agg.TotalFouls)
	}
	if agg.TotalTurnovers != totalTurnovers {
		t.Errorf("expected total_turnovers %d, got %d", totalTurnovers, agg.TotalTurnovers)
	}
	if agg.TotalMinutes != totalMinutes {
		t.Errorf("expected total_minutes %f, got %f", totalMinutes, agg.TotalMinutes)
	}

	// Ensure averages are calculated properly.
	expectedAvgPoints := float64(totalPoints) / float64(gamesPlayed)
	if agg.AvgPoints != expectedAvgPoints {
		t.Errorf("expected avg_points %f, got %f", expectedAvgPoints, agg.AvgPoints)
	}

	// Ensure all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestFetchTeamAggregate_Error(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %v", err)
	}
	defer db.Close()

	repo := repository.NewPlayerStatsRepository(db)

	teamID := "nonexistent"

	// Set up expected query returning an error.
	mock.ExpectQuery("SELECT (.+) FROM player_game_stats ps INNER JOIN players p ON ps.player_id = p.id WHERE p.team_id = \\$1").
		WithArgs(teamID).
		WillReturnError(sql.ErrNoRows)

	// Call FetchTeamAggregate.
	_, err = repo.FetchTeamAggregate(teamID)
	if err == nil {
		t.Error("expected error when team aggregate not found, got nil")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows error, got: %v", err)
	}

	// Ensure all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
