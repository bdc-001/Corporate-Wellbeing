#!/bin/bash

# Production Deployment Script
# Usage: ./deploy.sh [environment]

set -e

ENVIRONMENT=${1:-production}

echo "🚀 Deploying to $ENVIRONMENT environment..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found"
    echo "Please copy .env.example to .env and configure it"
    exit 1
fi

# Load environment variables
export $(cat .env | grep -v '^#' | xargs)

# Validate critical environment variables
if [ "$ENVIRONMENT" = "production" ]; then
    if [ "$JWT_SECRET" = "change-me-in-production" ]; then
        echo "❌ Error: JWT_SECRET must be changed for production"
        exit 1
    fi
    
    if [ "$DATABASE_URL" = "postgres://localhost/convin_crae?sslmode=disable" ]; then
        echo "❌ Error: DATABASE_URL must be configured for production"
        exit 1
    fi
fi

# Build Docker images
echo "📦 Building Docker images..."
docker-compose build

# Run database migrations
echo "🗄️  Running database migrations..."
docker-compose run --rm backend sh -c "psql \$DATABASE_URL -f /app/database/schema.sql || true"

# Start services
echo "🚀 Starting services..."
docker-compose up -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Health check
echo "🏥 Checking service health..."
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Backend is healthy"
else
    echo "❌ Backend health check failed"
    docker-compose logs backend
    exit 1
fi

if curl -f http://localhost/health > /dev/null 2>&1; then
    echo "✅ Frontend is healthy"
else
    echo "⚠️  Frontend health check failed (may still be starting)"
fi

echo "✅ Deployment complete!"
echo ""
echo "Services:"
echo "  - Backend: http://localhost:8080"
echo "  - Frontend: http://localhost"
echo "  - Health: http://localhost:8080/health"
echo ""
echo "View logs: docker-compose logs -f"

