#!/bin/bash

echo "🧹 Cleaning up old Docker containers and images..."

# Stop and remove running containers
docker-compose down

# Remove unused images, containers, and networks
docker system prune -af

echo "✅ Cleanup completed!"
