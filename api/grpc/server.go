package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/PAFFx/job-poll-queue/api/grpc/worker"
	pb "github.com/PAFFx/job-poll-queue/proto/worker"
	"github.com/PAFFx/job-poll-queue/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the gRPC server
type Server struct {
	grpcServer *grpc.Server
	jobQueue   *queue.Queue
	port       string
}

// NewServer creates a new gRPC server with the provided job queue
func NewServer(jobQueue *queue.Queue, port string) *Server {
	grpcServer := grpc.NewServer()

	// Create and register the worker service
	workerService := worker.NewService(jobQueue)
	pb.RegisterWorkerServiceServer(grpcServer, workerService)

	// Register reflection service for development tools like grpcurl
	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		jobQueue:   jobQueue,
		port:       port,
	}
}

// Start starts the gRPC server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", s.port, err)
	}

	log.Printf("Starting gRPC server on port %s", s.port)
	return s.grpcServer.Serve(listener)
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
