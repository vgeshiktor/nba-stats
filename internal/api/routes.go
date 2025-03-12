// internal/api/routes.go
package api

import (
	"database/sql"
	"net/http"

	"github.com/vgeshiktor/nba-stats/pkg/errors"
)

// RegisterRoutes maps URL endpoints to the corresponding handler functions.
// Each endpoint is wrapped with Logging, Authentication, and RequestTracing middleware.
func RegisterRoutes(mux *http.ServeMux, handler *Handler, db *sql.DB) {
	// Define a middleware chain.
	chain := func(h http.Handler) http.Handler {
		return ChainMiddleware(h, RequestTracingMiddleware, AuthenticationMiddleware, LoggingMiddleware)
	}

	// Player stats endpoints.
	mux.Handle("/api/v1/player-stats", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.LogPlayerStats(w, r)
			return
		}
		errors.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	})))

	mux.Handle("/api/v1/player-stats/player/", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetPlayerAggregate(w, r)
			return
		}
		errors.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	})))

	mux.Handle("/api/v1/player-stats/team/", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetTeamAggregate(w, r)
			return
		}
		errors.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	})))

	// Player management endpoints.
	mux.Handle("/api/v1/players", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreatePlayer(w, r)
			return
		}
		errors.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	})))
	mux.Handle("/api/v1/players/", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetPlayer(w, r)
			return
		}
		errors.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	})))

	// Team management endpoints.
	mux.Handle("/api/v1/teams", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateTeam(w, r)
			return
		}
		errors.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	})))
	mux.Handle("/api/v1/teams/", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetTeam(w, r)
			return
		}
		errors.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	})))

	// Game management endpoints.
	mux.Handle("/api/v1/games", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateGame(w, r)
			return
		}
		errors.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	})))
	mux.Handle("/api/v1/games/", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetGame(w, r)
			return
		}
		errors.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	})))
	
	mux.HandleFunc("/health/live", LivenessProbeHandler)
	mux.HandleFunc("/health/ready", ReadinessProbeHandler(db))
}
