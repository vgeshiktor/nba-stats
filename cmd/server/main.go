// cmd/server/main.go
package main

import (
	"github.com/vgeshiktor/nba-stats/internal/app"
	"github.com/vgeshiktor/nba-stats/pkg/logger"
)

func main() {
	// Initialize the application via our abstraction layer.
	server := app.Initialize()

	// Log and start the server.
	logger.Info("Starting server on " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Server failed: " + err.Error())
	}
}
