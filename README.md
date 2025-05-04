# Job Queue API

A lightweight job processing queue with HTTP and gRPC APIs.

## Features

- Synchronous job processing
- RESTful API using Gofiber framework
- gRPC API for efficient worker communication
- Persistent storage
- Thread-safe operations

## Installation

```bash
git clone https://github.com/PAFFx/job-poll-queue.git
cd job-poll-queue
go mod download
go build
```

## HTTP API Endpoints

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

## gRPC API

Worker service is also available via gRPC on port 50051 (configurable).

### Methods

- `RequestJob`: Retrieves the next available job
- `CompleteJob`: Submits results for a processed job

### Testing with grpcurl

```bash
# List available services
grpcurl -plaintext localhost:50051 list

# Request a job
grpcurl -plaintext -d '{}' localhost:50051 worker.WorkerService/RequestJob

# Complete a job
grpcurl -plaintext -d '{"job_id":"JOB_ID","payload":"result"}' \
  localhost:50051 worker.WorkerService/CompleteJob
```

## Architecture

The system now supports dual communication methods:

1. **HTTP API**: Traditional RESTful endpoints for client, worker, and admin operations
2. **gRPC API**: Efficient binary protocol for worker communication

The core components remain:
1. **Submit API**: Handles job submissions from clients
2. **Worker API**: Allows workers to poll for jobs and submit results (via HTTP or gRPC)
3. **Job Status Manager**: Tracks job state throughout its lifecycle

Jobs submitted by clients are held in a synchronous request until a worker processes the job and returns a result, which is then returned to the client.

## Project Structure

```
job-poll-queue/
├── api/              # API implementations
│   ├── http/         # HTTP API endpoints
│   │   ├── admin/    # Admin HTTP endpoints  
│   │   ├── submit/   # Client submission endpoints
│   │   ├── worker/   # Worker HTTP endpoints
│   │   └── server.go # HTTP server and routes
│   └── grpc/         # gRPC API endpoints
│       ├── worker/   # Worker gRPC service
│       └── server.go # gRPC server
├── proto/            # Protocol buffer definitions
│   └── worker/       # Worker service proto definitions
├── queue/            # Core queue implementation
│   ├── jobstatus.go  # Job status tracking
│   ├── queue.go      # Main queue functionality
│   └── storage.go    # Persistence layer
├── config/           # Configuration
│   └── env.go        # Environment variables
└── main.go           # Application entry point (runs both servers)
```

## License

MIT