package api

import (
	"github.com/PAFFx/job-poll-queue/api/admin"
	"github.com/PAFFx/job-poll-queue/api/submit"
	"github.com/PAFFx/job-poll-queue/api/worker"
	"github.com/PAFFx/job-poll-queue/queue"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Server represents the API server
type Server struct {
	app      *fiber.App
	jobQueue *queue.Queue
}

// NewServer creates a new API server with the provided job queue
func NewServer(jobQueue *queue.Queue) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	server := &Server{
		app:      app,
		jobQueue: jobQueue,
	}

	// Register routes
	server.registerRoutes()

	return server
}

// registerRoutes sets up all API routes
func (s *Server) registerRoutes() {
	api := s.app.Group("/api")

	adminGroup := api.Group("/admin")
	adminHandler := admin.NewHandler(s.jobQueue)
	adminHandler.RegisterRoutes(adminGroup)

	submitGroup := api.Group("/submit")
	submitHandler := submit.NewHandler(s.jobQueue)
	submitHandler.RegisterRoutes(submitGroup)

	workerGroup := api.Group("/worker")
	workerHandler := worker.NewHandler(s.jobQueue)
	workerHandler.RegisterRoutes(workerGroup)
}

// Start starts the API server on the given address
func (s *Server) Start(address string) error {
	return s.app.Listen(address)
}
