package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/PAFFx/job-poll-queue/api"
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
	server := api.NewServer(jobQueue)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start the server
	log.Printf("Starting job queue API server on port %s", port)
	if err := server.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
