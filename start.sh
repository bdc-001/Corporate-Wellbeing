#!/bin/bash

# Start both backend and frontend
# This script starts the backend in the background and frontend in foreground

cd "$(dirname "$0")"

echo "=========================================="
echo "Starting Convin Revenue Attribution Engine"
echo "=========================================="
echo ""

# Start backend in background
echo "Starting backend server..."
cd backend
export DATABASE_URL="postgres://localhost/convin_crae?sslmode=disable"
export PORT=8080
export ENVIRONMENT=development

go run cmd/server/main.go &
BACKEND_PID=$!

echo "Backend started with PID: $BACKEND_PID"
echo "Backend URL: http://localhost:8080"
echo ""

# Wait a bit for backend to start
sleep 3

# Start frontend
echo "Starting frontend server..."
cd ../frontend
npm start

# Cleanup on exit
trap "kill $BACKEND_PID 2>/dev/null" EXIT

