package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/PAFFx/job-poll-queue/api/http"
	"github.com/PAFFx/job-poll-queue/config"
	"github.com/PAFFx/job-poll-queue/queue"
)

func main() {
	// Setup the queue storage directory (using a 'data' folder in the project root)
	storageDir := filepath.Join(".", "data")

	// Create the job queue
	jobQueue, err := queue.NewQueue("jobs", storageDir)
	if err != nil {
		log.Fatalf("Failed to create job queue: %v", err)
	}

	// Create the API server
	server := http.NewServer(jobQueue)

	// Get port from environment variable or use default
	envVars, err := config.GetEnvVariables()
	if err != nil {
		log.Fatalf("Failed to get environment variables: %v", err)
	}
	// Start the server
	log.Printf("Starting job queue API server on port %s", envVars.Port)
	if err := server.Start(fmt.Sprintf(":%s", envVars.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
