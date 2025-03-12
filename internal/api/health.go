package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// jsonResponse is a helper function for writing JSON responses.
func jsonResponse(w http.ResponseWriter, status int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}

// LivenessProbeHandler responds to the liveness probe.
func LivenessProbeHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ReadinessProbeHandler checks database connectivity.
func ReadinessProbeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			jsonResponse(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unavailable",
				"error":  err.Error(),
			})
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}
