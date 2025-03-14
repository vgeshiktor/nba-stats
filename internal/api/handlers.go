// internal/api/handlers.go
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vgeshiktor/nba-stats/internal/domain"
	"github.com/vgeshiktor/nba-stats/internal/service"

	"github.com/vgeshiktor/nba-stats/pkg/errors"
	"github.com/vgeshiktor/nba-stats/pkg/logger"
)

// Handler aggregates all service dependencies for handling API requests.
type Handler struct {
	// Services for stats logging and aggregation.
	PlayerStatsService service.PlayerStatsService
	AggregationService service.AggregationService

	// Services for managing players, teams, and games.
	PlayerService service.PlayerService
	TeamService   service.TeamService
	GameService   service.GameService
}

// NewHandler creates a new API handler instance.
func NewHandler(
	playerStatsService service.PlayerStatsService,
	aggregationService service.AggregationService,
	playerService service.PlayerService,
	teamService service.TeamService,
	gameService service.GameService,
) *Handler {
	return &Handler{
		PlayerStatsService: playerStatsService,
		AggregationService: aggregationService,
		PlayerService:      playerService,
		TeamService:        teamService,
		GameService:        gameService,
	}
}

// LogPlayerStats handles POST /api/v1/player-stats to log game statistics.
func (h *Handler) LogPlayerStats(w http.ResponseWriter, r *http.Request) {
	var stats domain.PlayerGameStats
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&stats); err != nil {
		errors.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.PlayerStatsService.LogPlayerStats(&stats); err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "Error logging player stats: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Player stats logged successfully"})
}

// GetPlayerAggregate handles GET /api/v1/player-stats/player/{playerId} to fetch player aggregates.
func (h *Handler) GetPlayerAggregate(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {

		errors.WriteError(w, http.StatusBadRequest, "Player ID not provided")
		return
	}
	playerID := parts[5]

	logger.Info("get player aggreggate for id:  %s", playerID)

	aggregate, err := h.AggregationService.GetPlayerAggregate(playerID)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "Error fetching player aggregate: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aggregate)
}

// GetTeamAggregate handles GET /api/v1/player-stats/team/{teamId} to fetch team aggregates.
func (h *Handler) GetTeamAggregate(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		errors.WriteError(w, http.StatusBadRequest, "Team ID not provided")
		return
	}
	teamID := parts[5]

	logger.Info("get team aggreggate for id:  %s", teamID)

	aggregate, err := h.AggregationService.GetTeamAggregate(teamID)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "Error fetching team aggregate: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aggregate)
}

// CreatePlayer handles POST /api/v1/players to create a new player.
func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var player domain.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		errors.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.PlayerService.CreatePlayer(&player); err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "Error creating player: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)
}

// GetPlayer handles GET /api/v1/players/{playerId} to retrieve a player's details.
func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		errors.WriteError(w, http.StatusBadRequest, "Player ID not provided")
		return
	}
	playerID := parts[4]

	player, err := h.PlayerService.GetPlayerByID(playerID)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "Error fetching player: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// CreateTeam handles POST /api/v1/teams to create a new team.
func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team domain.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		errors.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.TeamService.CreateTeam(&team); err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "Error creating team: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
}

// GetTeam handles GET /api/v1/teams/{teamId} to retrieve team details.
func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		errors.WriteError(w, http.StatusBadRequest, "Team ID not provided")
		return
	}
	teamID := parts[4]

	team, err := h.TeamService.GetTeamByID(teamID)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "Error fetching team: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

// CreateGame handles POST /api/v1/games to create a new game.
func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game domain.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		errors.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.GameService.CreateGame(&game); err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "Error creating game: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}

// GetGame handles GET /api/v1/games/{gameId} to retrieve game details.
func (h *Handler) GetGame(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		errors.WriteError(w, http.StatusBadRequest, "Game ID not provided")
		return
	}
	gameID := parts[4]

	game, err := h.GameService.GetGameByID(gameID)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "Error fetching game: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}
