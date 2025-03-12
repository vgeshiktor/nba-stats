package app

import (
	"database/sql"
	"os"
	"github.com/vgeshiktor/nba-stats/pkg/logger"
)

// RunMigrations reads the SQL migration file and executes it against the database.
func RunMigrations(db *sql.DB, schemaPath string) error {
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	// Execute the migration script. Using Exec ensures that all statements in the file are run.
	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}

	logger.Info("Database migrations applied successfully.")
	return nil
}
