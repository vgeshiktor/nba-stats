package repository_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
)

// Function to create tables in the test database
func setupTestDatabase(t *testing.T) *sql.DB {
	// Open a connection to the in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:") // For in-memory DB
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Read the SQL migration script
	schema, err := os.ReadFile("./db_schema.sql")
	if err != nil {
		t.Fatalf("Failed to read schema file: %v", err)
	}

	// Execute the schema SQL to create tables
	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatalf("Failed to execute schema: %v", err)
	}

	return db
}

// Example test for checking the setup
func TestDatabaseConnectionAndTables(t *testing.T) {
	// Set up the test database and create tables
	db := setupTestDatabase(t)
	defer db.Close()

	cwd, _ := os.Getwd()
	t.Logf("working directory: %s",  cwd)

	// Verify that the tables were created (check players table exists)
	var exists bool
	err := db.QueryRow("SELECT count(*) > 0 FROM sqlite_master WHERE type='table' AND name='players'").Scan(&exists)
	if err != nil || !exists {
		t.Fatalf("players table not created: %v", err)
	}

	// Optionally, check other tables as well (e.g., teams, games)
	err = db.QueryRow("SELECT count(*) > 0 FROM sqlite_master WHERE type='table' AND name='teams'").Scan(&exists)
	if err != nil || !exists {
		t.Fatalf("teams table not created: %v", err)
	}

	// If no errors, the database setup is successful
	assert.True(t, exists, "Database tables are set up correctly")
}
