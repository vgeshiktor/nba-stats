package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver

	config "github.com/vgeshiktor/nba-stats/config"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "NBA-STATS: ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting NBA Stats API...")
	
	// Load configuration from environment variables
	config := loadConfig()
	logger.Printf("Configuration loaded: Server port %d, DB host %s", config.ServerPort, config.DBHost)
	
	// Connect to database
	db, err := connectDB(config)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Println("Database connection established")
	
	// Set up application components
	app := setupApplication(db, logger)
	
	// Configure HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.ServerPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	
	// Start server in a goroutine to allow graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("Server listening on port %d", config.ServerPort)
		serverErrors <- server.ListenAndServe()
	}()
	
	// Wait for interrupt signal or server error
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	
	select {
	case err := <-serverErrors:
		logger.Fatalf("Server error: %v", err)
	case <-shutdown:
		logger.Println("Server shutdown initiated...")
		
		// Create a deadline for the shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		// Gracefully shutdown the server
		if err := server.Shutdown(ctx); err != nil {
			logger.Printf("Server forced to shutdown: %v", err)
			if err := server.Close(); err != nil {
				logger.Fatalf("Server close error: %v", err)
			}
		}
	}
	
	logger.Println("Server shutdown complete")
}

// Load configuration from environment variables with sensible defaults
func loadConfig() config.Config {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnvAsInt("DB_PORT", 5432)
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "nba_stats")
	serverPort := getEnvAsInt("SERVER_PORT", 8080)
	
	return config.Config {
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		ServerPort: serverPort,
	}
}