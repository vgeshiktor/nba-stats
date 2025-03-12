package integration_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"github.com/vgeshiktor/nba-stats/internal/api"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestLivenessProbe(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health/live", api.LivenessProbeHandler)

	req, _ := http.NewRequest("GET", "/health/live", nil)
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"status":"ok"`)
}

func TestReadinessProbe_Healthy(t *testing.T) {
	// Use an in-memory SQLite database for testing readiness check
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Printf("Error during db creation: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health/ready", api.ReadinessProbeHandler(db))

	req, _ := http.NewRequest("GET", "/health/ready", nil)
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"status":"ready"`)
}
