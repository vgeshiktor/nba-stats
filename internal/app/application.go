package app

import (
	"log"
	"encoding/json"
	"net/http"
	"database/sql"

	config "github.com/vgeshiktor/nba-stats/config"
)


// Application struct holds all application components
type Application struct {
	config       config.Config
	logger       *log.Logger
	playerRepo   PlayerRepository
	teamRepo     TeamRepository
	statsRepo    GameStatsRepository
	statsService *StatsService
}

// Setup all application components
func setupApplication(db *sql.DB, logger *log.Logger) *Application {
	// Initialize repositories
	playerRepo := NewPlayerRepository(db)
	teamRepo := NewTeamRepository(db)
	statsRepo := NewGameStatsRepository(db)
	
	// Initialize services
	statsService := NewStatsService(playerRepo, teamRepo, statsRepo)
	
	return &Application{
		logger:       logger,
		playerRepo:   playerRepo,
		teamRepo:     teamRepo,
		statsRepo:    statsRepo,
		statsService: statsService,
	}
}

// Configure all HTTP routes
func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	
	// Register API endpoints
	mux.HandleFunc("/api/v1/player-stats", app.logPlayerStatsHandler)
	mux.HandleFunc("/api/v1/player-stats/player/", app.getPlayerSeasonAverageHandler)
	mux.HandleFunc("/api/v1/player-stats/team/", app.getTeamSeasonAverageHandler)
	
	return mux
}

// HTTP Handlers for API endpoints
func (app *Application) logPlayerStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Decode JSON request
	var input PlayerStatsInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.logger.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if err := validatePlayerStats(&input); err != nil {
		app.logger.Printf("Validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Process stats through service layer
	if err := app.statsService.LogPlayerStats(&input); err != nil {
		app.logger.Printf("Error logging stats: %v", err)
		http.Error(w, "Failed to log player statistics", http.StatusInternalServerError)
		return
	}
	
	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Player statistics logged successfully"})
}
