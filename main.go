package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	grpcServer "github.com/PAFFx/job-poll-queue/api/grpc"
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

	// Get port from environment variable or use default
	envVars, err := config.GetEnvVariables()
	if err != nil {
		log.Fatalf("Failed to get environment variables: %v", err)
	}

	// Create the HTTP API server
	httpServer := http.NewServer(jobQueue)

	// Create the gRPC server
	grpcSrv := grpcServer.NewServer(jobQueue, envVars.GrpcPort)

	// Start both servers in goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Start HTTP server
	go func() {
		defer wg.Done()
		log.Printf("Starting HTTP API server on port %s", envVars.Port)
		if err := httpServer.Start(fmt.Sprintf(":%s", envVars.Port)); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start gRPC server
	go func() {
		defer wg.Done()
		if err := grpcSrv.Start(); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Wait for both servers
	wg.Wait()
}
