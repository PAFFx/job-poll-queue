package api

import (
	"github.com/PAFFx/job-poll-queue/queue"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
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

	// Job queue endpoints
	jobs := api.Group("/jobs")
	jobs.Post("/", s.pushJob)           // Submit a new job
	jobs.Get("/", s.getQueueStatus)     // Get queue status
	jobs.Get("/next", s.getNextJob)     // Get the next job from the queue
	jobs.Delete("/clear", s.clearQueue) // Clear the queue

	// Health check
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}

// Start starts the API server on the given address
func (s *Server) Start(address string) error {
	return s.app.Listen(address)
}

// pushJob handles job submission
func (s *Server) pushJob(c *fiber.Ctx) error {
	// Parse job data from request
	type JobRequest struct {
		Payload string            `json:"payload"`
		Headers map[string]string `json:"headers,omitempty"`
	}

	var jobReq JobRequest
	if err := c.BodyParser(&jobReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid job data")
	}

	// Create a new job
	job := queue.Message{
		ID:      uuid.New().String(),
		Payload: jobReq.Payload,
		Headers: jobReq.Headers,
	}

	// Push to queue
	if err := s.jobQueue.Push(job); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to push job to queue")
	}

	// Return job info
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Job submitted successfully",
		"job_id":  job.ID,
	})
}

// getQueueStatus returns the current status of the job queue
func (s *Server) getQueueStatus(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"queue_size": s.jobQueue.Size(),
	})
}

// getNextJob gets the next job from the queue
func (s *Server) getNextJob(c *fiber.Ctx) error {
	job, err := s.jobQueue.Pop()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to pop job from queue")
	}

	if job == nil {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"message": "Queue is empty",
		})
	}

	return c.JSON(job)
}

// clearQueue clears all jobs from the queue
func (s *Server) clearQueue(c *fiber.Ctx) error {
	if err := s.jobQueue.Clear(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to clear queue")
	}

	return c.JSON(fiber.Map{
		"message": "Queue cleared successfully",
	})
}
