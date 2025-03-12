# NBA Stats API

## Table of Contents

- [NBA Stats API](#nba-stats-api)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Project Structure](#project-structure)
  - [Prerequisites](#prerequisites)
  - [Setup and Running Locally](#setup-and-running-locally)
  - [API Endpoints](#api-endpoints)
      - [Player Statistics:](#player-statistics)
      - [Player Management:](#player-management)
      - [Team Management:](#team-management)
      - [Game Management:](#game-management)
        - [To run all tests in the project, execute:](#to-run-all-tests-in-the-project-execute)
  - [Deployment](#deployment)
        - [Use the deployment script to build, containerize, and deploy the application:](#use-the-deployment-script-to-build-containerize-and-deploy-the-application)
  - [Cleanup](#cleanup)
        - [To clean up old Docker containers and images, run:](#to-clean-up-old-docker-containers-and-images-run)
  - [Logging, Middleware, and Configuration](#logging-middleware-and-configuration)
  - [Future Enhancements](#future-enhancements)

## Overview

The NBA Stats API is a backend system designed to log NBA player statistics and calculate season aggregates for players and teams. It is built in Go using a clean architecture that separates concerns across API, domain, repository, service, and utility layers. The application is containerized using Docker and can be deployed using Docker Compose.

## Project Structure
```
nba-stats/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── api/                 # API Layer
│   │   ├── handlers.go      # Defines HTTP handlers
│   │   ├── routes.go        # Registers API routes
│   │   ├── middleware.go    # Middleware (logging, auth, tracing)
│   ├── app/
│   │   └── app.go           # Application initialization (DB, services, router)
│   ├── domain/              # Domain Layer (Entities & Errors)
│   │   ├── models.go        # Defines core data models
│   │   ├── errors.go        # Custom application errors
│   ├── repository/          # Data Persistence Layer (Repositories)
│   │   ├── db.go            # Database connection & pooling
│   │   ├── player_repository.go
│   │   ├── team_repository.go
│   │   ├── game_repository.go
│   │   ├── player_stats_repository.go
│   ├── service/             # Business Logic Layer (Services)
│   │   ├── player_service.go
│   │   ├── team_service.go
│   │   ├── game_service.go
│   │   ├── player_stats_service.go
│   │   ├── aggregation_service.go
├── pkg/                     # Utility Packages (Reusable)
│   ├── logger/
│   │   └── logger.go        # Structured JSON Logger
│   ├── validator/
│   │   └── validator.go     # Input validation utilities
│   ├── errors/
│   │   └── error_handling.go # Standardized error responses
├── scripts/                 # Deployment & Automation Scripts
│   ├── deploy.sh            # Build, Dockerize, and Deploy
│   ├── cleanup.sh           # Cleanup old containers/images
├── test/                    # Integration & Unit Tests
│   ├── player_service_test.go
│   ├── game_service_test.go
│   ├── team_service_test.go
├── Dockerfile               # Multi-stage Docker Build
├── docker-compose.yml        # Local environment setup
├── go.mod                    # Go module file
├── go.sum                    # Dependencies lock file
└── README.md                 # Project Documentation
```

## Prerequisites

- Go 1.23+
- Docker & Docker Compose

## Setup and Running Locally

1. **Clone the repository:**
```sh
   git clone https://github.com/your-username/nba-stats.git
   cd nba-stats
```
2. **Set Environment Variables:**
   Create a .env file or export the following variables:
```
    PORT (default: 8080)
    DATABASE_URL (e.g., postgres://nba_user:nba_password@db:5432/nba_stats?sslmode=disable)
    DB_MAX_OPEN_CONNS (default: 25)
    DB_MAX_IDLE_CONNS (default: 25)
    DB_CONN_MAX_LIFETIME (default: "5m")
```
3. **Run the Application:**
- Using Docker Compose:
  ```sh
    docker-compose up -d --build
  ```
- Without Docker:
Build the binary and run it:
```sh
go build -o nba-stats cmd/server/main.go
./nba-stats
```
## API Endpoints
#### Player Statistics:
- POST /api/v1/player-stats
Log player statistics.

- GET /api/v1/player-stats/player/{playerId}
Retrieve season aggregate stats for a player.

- GET /api/v1/player-stats/team/{teamId}
Retrieve season aggregate stats for a team.
#### Player Management:
- POST /api/v1/players
Create a new player.

- GET /api/v1/players/{playerId}
Retrieve details for a specific player.

#### Team Management:
- POST /api/v1/teams
Create a new team.

- GET /api/v1/teams/{teamId}
Retrieve details for a specific team.

#### Game Management:
- POST /api/v1/games
Create a new game.

- GET /api/v1/games/{gameId}
Retrieve details for a specific game.
5. **Running Tests:**
##### To run all tests in the project, execute:
```sh
go test ./...
```
## Deployment
##### Use the deployment script to build, containerize, and deploy the application:
```sh
./scripts/deploy.sh
```
## Cleanup
##### To clean up old Docker containers and images, run:
```sh
./scripts/cleanup.sh
```
## Logging, Middleware, and Configuration
- Logging: Uses a JSON-structured logger for detailed, machine-readable logs.
  
- Middleware: Includes logging, authentication, and request tracing middleware.

- Configuration: Loads environment variables to configure the application, including database connection pool settings.

## Future Enhancements
- Implement more comprehensive integration tests.
  
- Enhance authentication and authorization mechanisms.
  
- Integrate with a centralized logging solution (e.g., ELK, Splunk).
  
