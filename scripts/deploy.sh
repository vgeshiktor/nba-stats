#!/bin/bash

set -e  # Exit immediately if a command fails

APP_NAME="nba-stats"
DOCKER_IMAGE="nba-stats:latest"
DOCKER_COMPOSE_FILE="docker-compose.yml"

echo "🚀 Starting deployment process..."

# Step 1: Build the Go application
echo "🔨 Building Go application..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/$APP_NAME cmd/server/main.go

# Step 2: Build the Docker image
echo "📦 Building Docker image: $DOCKER_IMAGE..."
docker build -t $DOCKER_IMAGE .

# Step 3: Deploy using Docker Compose
echo "🚢 Deploying application with Docker Compose..."
docker-compose -f $DOCKER_COMPOSE_FILE up -d --build

echo "✅ Deployment completed successfully!"
