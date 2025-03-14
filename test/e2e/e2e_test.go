package e2e_test

import (
	"context"
	"time"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"testing"

	"github.com/vgeshiktor/nba-stats/internal/app"
	"github.com/vgeshiktor/nba-stats/pkg/logger"
)

var testServer *http.Server

// Setup the test environment
func TestMain(m *testing.M) {
	os.Setenv("DATABASE_URL", "postgres://nba_user:nba_password@localhost:5432/nba_stats?sslmode=disable")
	os.Setenv("TEST_MODE", "true")

	// Start the app for testing
	testServer = app.Initialize()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if err := testServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP shutdown error: %v", err)
		}

		// Clean up (e.g., close DB connections)
	}()

	// Run all tests
	code := m.Run()

	os.Exit(code)
}
