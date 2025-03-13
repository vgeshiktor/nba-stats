// internal/repository/db.go
package repository

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/vgeshiktor/nba-stats/pkg/logger"
)

// NewDB establishes a database connection using the provided connection string
// and applies connection pool settings.
func NewDB(connStr string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (*sql.DB, error) {
	var driverName string
	if connStr == ":memory:" {
		driverName = "sqlite3"
	} else {
		driverName = "postgres"
	}

	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connectivity.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool parameters.
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	logger.Info("Database connection pool established (MaxOpenConns: %d, MaxIdleConns: %d, ConnMaxLifetime: %s)",
		maxOpenConns, maxIdleConns, connMaxLifetime)
	return db, nil
}