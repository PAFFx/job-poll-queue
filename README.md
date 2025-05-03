# Job Queue API

A lightweight synchronous job processing queue with HTTP API.

## Features

- Synchronous job processing
- RESTful API using Gofiber framework
- Persistent storage
- Thread-safe operations

## Installation

```bash
git clone https://github.com/PAFFx/job-poll-queue.git
cd job-poll-queue
go mod download
go build
```

## API Endpoints

### Client Endpoints

#### Submit a job (synchronous)

```
POST /api/submit
```

Request: Send any JSON payload in the body  
Response: Returns the worker's processed result

```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "payload": {/* Worker's response payload */},
  "created_at": "2023-06-01T12:34:56Z",
  "completed_at": "2023-06-01T12:34:58Z"
}
```

### Worker Endpoints

#### Poll for a job

```
GET /api/worker/poll
```

Response (when job available):
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "payload": {/* Original job payload */},
  "headers": {/* Original request headers */}
}
```

#### Complete a job

```
POST /api/worker/complete/:id
```

Request: Send any payload in the body (can be JSON, text, or any format)  
Response:
```json
{
  "message": "Job completed",
  "job_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Admin Endpoints

#### Get queue stats

```
GET /api/admin/stats
```

#### Clear the queue

```
DELETE /api/admin/clear
```

#### Get next job (admin only)

```
GET /api/admin/next
```

### Health check

```
GET /health
```

## Architecture

The system consists of three main components:

1. **Submit API**: Handles job submissions from clients
2. **Worker API**: Allows workers to poll for jobs and submit results
3. **Job Status Manager**: Tracks job state throughout its lifecycle

Jobs submitted by clients are held in a synchronous request until a worker processes the job and returns a result, which is then returned to the client.

## Project Structure

```
job-poll-queue/
├── api/               # API implementation
│   ├── admin/         # Admin API endpoints
│   ├── submit/        # Client submission endpoints
│   ├── worker/        # Worker endpoints
│   └── server.go      # API server and routes
├── queue/             # Queue implementation
│   ├── jobstatus.go   # Job status tracking
│   ├── queue.go       # Main queue functionality
│   └── storage.go     # Persistence layer
└── main.go            # Application entry point
```

## License

MIT