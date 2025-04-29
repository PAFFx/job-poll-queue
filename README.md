# Job Queue API

A simple job queue API built with Go and Gofiber that allows you to push, pop, and manage jobs via HTTP endpoints.

## Features

- FIFO (First In, First Out) job queue
- RESTful API using Gofiber framework
- Persistent storage of jobs to disk
- Thread-safe operations

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/PAFFx/job-poll-queue.git
   cd job-poll-queue
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Build the application:
   ```
   go build -o job-queue-api
   ```

## Usage

### Starting the server

Run the application:
```
./job-queue-api
```

By default, the server listens on port 3000. You can change the port by setting the `PORT` environment variable:
```
PORT=8080 ./job-queue-api
```

### API Endpoints

#### Submit a job

```
POST /api/jobs
```

Request body:
```json
{
  "payload": "Job data goes here",
  "headers": {
    "priority": "high",
    "customField": "value"
  }
}
```

Response:
```json
{
  "message": "Job submitted successfully",
  "job_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Get queue status

```
GET /api/jobs
```

Response:
```json
{
  "queue_size": 5
}
```

#### Get the next job

```
GET /api/jobs/next
```

Response (if job exists):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "payload": "Job data goes here",
  "headers": {
    "priority": "high",
    "customField": "value"
  }
}
```

Response (if queue is empty):
```json
{
  "message": "Queue is empty"
}
```

#### Clear the queue

```
DELETE /api/jobs/clear
```

Response:
```json
{
  "message": "Queue cleared successfully"
}
```

#### Health check

```
GET /health
```

Response:
```json
{
  "status": "ok"
}
```

## Development

### Project Structure

```
job-poll-queue/
├── data/               # Queue storage directory (created at runtime)
├── src/
│   ├── api/            # API implementation
│   │   └── server.go   # API server and route handlers
│   └── queue/          # Queue implementation
│       └── queue.go    # FIFO queue with disk persistence
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── main.go             # Application entry point
└── README.md           # This file
```

### Running in Development Mode

For development, you can use tools like Air for live reloading:

```
go install github.com/cosmtrek/air@latest
air
```

## License

MIT