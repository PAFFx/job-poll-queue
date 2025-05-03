package worker

import (
	"github.com/PAFFx/job-poll-queue/queue"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	jobQueue *queue.Queue
}

func NewHandler(jobQueue *queue.Queue) *Handler {
	return &Handler{jobQueue: jobQueue}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/poll", h.RequestJobHandler)
	router.Post("/complete/:id", h.CompleteJobHandler)
}

func (h *Handler) RequestJobHandler(c *fiber.Ctx) error {
	job, err := h.RequestJob()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to poll for job")
	}

	if job == nil {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"message": "No job available",
		})
	}

	return c.JSON(h.FormatJobResponse(job))
}

func (h *Handler) CompleteJobHandler(c *fiber.Ctx) error {
	jobID := c.Params("id")
	if err := h.ValidateJobID(jobID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Get the raw request body as the result
	result := string(c.Body())

	// If body is empty, set a default message
	if result == "" {
		result = "{}"
	}

	// Complete the job with the provided payload
	if err := h.CompleteJob(jobID, result); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to complete job: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Job completed",
		"job_id":  jobID,
	})
}
