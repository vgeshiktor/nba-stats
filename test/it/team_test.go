package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/app"
	"github.com/vgeshiktor/nba-stats/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetTeam(t *testing.T) {
	os.Setenv("DATABASE_URL", ":memory:")

	server := app.Initialize()

	// Create a new team.
	newTeam := domain.Team{
		ID:   "team1",
		Name: "Team One",
	}
	teamBody, err := json.Marshal(newTeam)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/teams", bytes.NewBuffer(teamBody))
	req.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	resp := httptest.NewRecorder()
	server.Handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Retrieve the team.
	reqGet, _ := http.NewRequest("GET", "/api/v1/teams/team1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqGet.Header.Set("Authorization", "dummy-token")

	respGet := httptest.NewRecorder()
	server.Handler.ServeHTTP(respGet, reqGet)
	assert.Equal(t, http.StatusOK, respGet.Code)

	var retrievedTeam domain.Team
	err = json.Unmarshal(respGet.Body.Bytes(), &retrievedTeam)
	assert.NoError(t, err)
	assert.Equal(t, newTeam.ID, retrievedTeam.ID)
	assert.Equal(t, newTeam.Name, retrievedTeam.Name)
}
