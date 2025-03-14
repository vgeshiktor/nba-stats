// test/e2e/team_e2e_test.go
package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/vgeshiktor/nba-stats/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetTeamE2E(t *testing.T) {
	// --- Step 1: Create a new team ---
	newTeam := domain.Team{
		ID:   "team1",
		Name: "Test Team",
	}
	teamBody, err := json.Marshal(newTeam)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/teams", bytes.NewBuffer(teamBody))
	req.Header.Set("Content-Type", "application/json")

	// For endpoints with middleware, add a dummy Authorization header.
	req.Header.Set("Authorization", "dummy-token")

	resp := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// --- Step 2: Retrieve the created team ---
	reqGet, _ := http.NewRequest("GET", "/api/v1/teams/team1", nil)

	// For endpoints with middleware, add a dummy Authorization header.
	reqGet.Header.Set("Authorization", "dummy-token")

	respGet := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(respGet, reqGet)
	assert.Equal(t, http.StatusOK, respGet.Code)

	var retrievedTeam domain.Team
	err = json.Unmarshal(respGet.Body.Bytes(), &retrievedTeam)
	assert.NoError(t, err)

	// --- Assert that the retrieved team matches the created team ---
	assert.Equal(t, newTeam.ID, retrievedTeam.ID)
	assert.Equal(t, newTeam.Name, retrievedTeam.Name)
}
