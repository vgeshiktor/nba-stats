package integration_test

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://nba_user:nba_password@localhost:5432/nba_stats?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	assert.NoError(t, err, "Database should be accessible")
}
