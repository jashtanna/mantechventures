#!/bin/bash

echo "Starting Video Ad Tracker Service"
echo "================================"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running"
    exit 1
fi

# Start the services
echo "Starting services with Docker Compose..."
docker-compose up -d

# Wait for services to be ready
echo "Waiting for services to start..."
sleep 10

# Check if services are running
if docker-compose ps | grep -q "Up"; then
    echo "Services started successfully"
    echo "API available at: http://localhost:8080"
    echo "Health check: http://localhost:8080/health"
    echo ""
    echo "To test the API, run: ./scripts/test_api.sh"
    echo "To stop services, run: docker-compose down"
else
    echo "Error: Services failed to start"
    docker-compose logs
    exit 1
fi
