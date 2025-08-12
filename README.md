# Video Advertisement Tracking Service

A Go backend service for managing and tracking video advertisements with real-time analytics and high throughput processing.

## Features

- Ad management and retrieval
- Asynchronous click event processing
- Real-time performance analytics with granular time periods
- Fault-tolerant data handling
- Monitoring and logging
- Scalable architecture

## Architecture

### Key Components

1. Async Processing: Click events processed in background for fast response
2. Database: PostgreSQL with optimized indexes for analytics
3. Service Layer: Clean separation of business logic
4. Monitoring: Metrics and structured logging

### Data Flow

1. Click Recording: Client request → Immediate response → Background processing
2. Analytics: Real-time aggregation of click events
3. Recovery: Automatic processing of unprocessed events

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Option 1: Docker Setup (Recommended)

The `run.sh` script automates the entire setup process:

```bash
# Clone the repository
git clone <repository-url>
cd video-ad-tracker

# Start all services (database + application)
./run.sh

# Test the API
curl http://localhost:8080/api/v1/ads
```

**What the run.sh script does:**
1. Checks if Docker is running
2. Starts PostgreSQL database and Go application using Docker Compose
3. Waits for services to be ready
4. Verifies all services are running
5. Provides API endpoints and testing instructions

### Option 2: Local Development

For development without Docker:

```bash
# Install dependencies
go mod download

# Start PostgreSQL database (using Docker)
docker run -d --name postgres-db \
  -e POSTGRES_DB=video_ads \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 postgres:15-alpine

# Set environment variables
export DATABASE_URL="postgres://postgres:password@localhost:5432/video_ads?sslmode=disable"
export PORT="8080"
export LOG_LEVEL="info"

# Run application
go run main.go
```

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/ads` | List all advertisements |
| `POST` | `/ads/click` | Record click event (async) |
| `GET` | `/ads/analytics` | Get performance metrics |
| `GET` | `/ads/analytics/hourly` | Get hourly breakdown for last 24h |

| `GET` | `/metrics` | Prometheus metrics |

### Testing with cURL



**Get Advertisements:**
```bash
curl http://localhost:8080/api/v1/ads
```

**Record Click:**
```bash
curl -X POST http://localhost:8080/api/v1/ads/click \
  -H "Content-Type: application/json" \
  -d '{
    "ad_id": 1,
    "video_playback_time": 15.5,
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0 (Test Browser)"
  }'
```

**Get Analytics:**
```bash
# Basic analytics
curl "http://localhost:8080/api/v1/ads/analytics?timeframe=24h"

# Hourly breakdown for last 24 hours
curl "http://localhost:8080/api/v1/ads/analytics/hourly"
```

**Get Metrics:**
```bash
curl http://localhost:8080/metrics
```

### Testing with Postman

**Import Postman Collection:**
1. Download the `Video_Ad_Tracker_API.postman_collection.json` file
2. Open Postman and click "Import"
3. Drag and drop the JSON file or click "Upload Files"
4. The collection includes all 13 API endpoints with pre-configured requests

**Collection Features:**
- ✅ All endpoints pre-configured with correct headers and body data
- ✅ Environment variable `{{base_url}}` set to `http://localhost:8080`
- ✅ Descriptive names and documentation for each endpoint
- ✅ Ready to use - no manual setup required
- ✅ Easy to switch between environments (local, staging, production)

**Included Endpoints:**
- Get All Advertisements
- Record Click Event
- Analytics (24h timeframe)
- Hourly Analytics Breakdown
- Prometheus Metrics

**Quick Start:**
1. Start the application: `./run.sh`
2. Import the Postman collection
3. Click "Send" on any request to test the API
4. Update the `base_url` variable for different environments

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:password@localhost:5432/video_ads?sslmode=disable` |
| `LOG_LEVEL` | Logging level | `info` |

## Database Schema

### Tables

#### ads
- `id` (SERIAL PRIMARY KEY)
- `image_url` (VARCHAR(500))
- `target_url` (VARCHAR(500))
- `title` (VARCHAR(200))
- `description` (TEXT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

#### click_events
- `id` (SERIAL PRIMARY KEY)
- `ad_id` (INTEGER REFERENCES ads(id))
- `timestamp` (TIMESTAMP)
- `ip_address` (VARCHAR(45))
- `video_playback_time` (DECIMAL(10,2))
- `user_agent` (TEXT)
- `processed` (BOOLEAN)

### Indexes
- `idx_click_events_ad_id` on `click_events(ad_id)`
- `idx_click_events_timestamp` on `click_events(timestamp)`
- `idx_click_events_processed` on `click_events(processed)`

## Monitoring

### Prometheus Metrics

- `http_requests_total`: Total HTTP requests by method, endpoint, and status
- `http_request_duration_seconds`: Request duration histogram

### Logging

Structured JSON logging with Logrus including:
- Request timestamps
- HTTP status codes
- Request latency
- Client IP addresses
- Error messages

## Scalability Features

1. Async Processing: Click events processed in background goroutines
2. Database Optimization: Proper indexing for analytics queries
3. Connection Pooling: Efficient database connection management
4. Graceful Shutdown: Proper cleanup on service termination

## Testing



### Manual Testing
After starting the application, you can test each endpoint manually:

1. **Health Check**: Verify the service is running
2. **Get Ads**: View available advertisements
3. **Record Click**: Test async click processing
4. **Get Analytics**: View performance metrics
5. **Get Metrics**: View Prometheus metrics

### Using Makefile
```bash
# See all available commands
make help

# Build application
make build

# Run with Docker
make docker-run
```









## CI/CD Pipeline

The project includes automated testing and build pipelines using GitHub Actions:

### Automated Building
- **Docker Build**: Automated Docker image building
- **Security Scan**: Vulnerability scanning with Trivy

### Pipeline Stages
1. **Build**: Build and test Docker image
2. **Security**: Vulnerability scanning with Trivy

### Local CI/CD
```bash
# Run CI build locally
make ci-build
```

## Production Deployment

```bash
# Build and run with Docker
docker build -t video-ad-tracker .
docker run -p 8080:8080 video-ad-tracker
```

## Project Structure

```
video-ad-tracker/
├── run.sh                 # Automated startup script
├── main.go                # Application entry point
├── Makefile               # Development and CI commands
├── docker-compose.yml     # Development environment
├── Dockerfile             # Production container
├── Video_Ad_Tracker_API.postman_collection.json  # Postman collection
├── .github/workflows/     # CI/CD pipelines
│   ├── ci-cd.yml         # Full CI/CD pipeline
│   └── test-and-build.yml # Basic test and build

└── internal/
    ├── config/            # Configuration management
    ├── database/          # Database setup and schema
    ├── models/            # Data structures
    ├── services/          # Business logic layer
    ├── handlers/          # HTTP request handlers
    └── middleware/        # Logging and metrics
```





## License

MIT License
