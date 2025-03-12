// internal/app/app.go
package app

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/vgeshiktor/nba-stats/internal/api"
	"github.com/vgeshiktor/nba-stats/internal/repository"
	"github.com/vgeshiktor/nba-stats/internal/service"
	"github.com/vgeshiktor/nba-stats/pkg/logger"
)

// AppConfig holds the configuration settings for the application.
type AppConfig struct {
	Port            string
	DBConnStr       string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewConfig reads environment variables and returns an AppConfig.
func NewConfig() AppConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		logger.Error("DATABASE_URL environment variable not set")
	}

	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 25)
	connMaxLifetime := getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	return AppConfig{
		Port:            port,
		DBConnStr:       dbConnStr,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}
}

// getEnvAsInt retrieves an environment variable as an integer.
func getEnvAsInt(name string, defaultValue int) int {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		logger.Error("Invalid value for %s, using default: %d", name, defaultValue)
		return defaultValue
	}
	return val
}

// getEnvAsDuration retrieves an environment variable as a time.Duration.
func getEnvAsDuration(name string, defaultValue time.Duration) time.Duration {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultValue
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		logger.Error("Invalid value for %s, using default: %v", name, defaultValue)
		return defaultValue
	}
	return val
}

// Initialize sets up the database connection, repositories, services, API handlers, and HTTP router.
func Initialize() *http.Server {
	// Load configuration
	config := NewConfig()

	// Establish database connection with the configured pool settings.
	db, err := repository.NewDB(config.DBConnStr, config.MaxOpenConns, config.MaxIdleConns, config.ConnMaxLifetime)
	if err != nil {
		logger.Error("Failed to connect to the database (conn str: %v): %v", config.DBConnStr, err)
	}

	// Run migrations to create necessary tables
	err = RunMigrations(db, "migrations/db_schema.sql")
	if err != nil {
		logger.Error("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	playerRepo := repository.NewPlayerRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	gameRepo := repository.NewGameRepository(db)
	statsRepo := repository.NewPlayerStatsRepository(db)

	// Initialize service layers
	statsService := service.NewPlayerStatsService(playerRepo, teamRepo, gameRepo, statsRepo)
	aggregationService := service.NewAggregationService(statsRepo)
	playerService := service.NewPlayerService(playerRepo)
	teamService := service.NewTeamService(teamRepo)
	gameService := service.NewGameService(gameRepo)

	// Initialize API handlers and register routes
	apiHandler := api.NewHandler(
		statsService,
		aggregationService,
		playerService,
		teamService,
		gameService,
	)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, apiHandler, db)

	// Create and return the configured HTTP server
	return &http.Server{
		Addr:         ":" + config.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
