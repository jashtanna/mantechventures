.PHONY: help test build run clean docker-build docker-run docker-test

# Default target
help:
	@echo "Available commands:"

	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"


	@echo "  ci-build     - Run CI build"



# Build the application
build:
	go build -o video-ad-tracker main.go

# Run the application locally
run:
	go run main.go

# Clean build artifacts
clean:
	rm -f video-ad-tracker coverage.out coverage.html
	go clean

# Build Docker image
docker-build:
	docker build -t video-ad-tracker .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose
docker-stop:
	docker-compose down





# CI build (for GitHub Actions)
ci-build:
	docker build -t video-ad-tracker .
	docker run -d --name test-app -p 8080:8080 \
		-e DATABASE_URL="postgres://postgres:password@host.docker.internal:5432/video_ads_test?sslmode=disable" \
		video-ad-tracker
	sleep 10
	curl -f http://localhost:8080/health || exit 1
	docker stop test-app
	docker rm test-app


