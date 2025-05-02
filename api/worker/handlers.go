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
	router.Post("/fail/:id", h.FailJobHandler)
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

	return c.JSON(fiber.Map{
		"job_id":  job.ID,
		"payload": job.Payload,
		"headers": job.Headers,
	})
}

func (h *Handler) CompleteJobHandler(c *fiber.Ctx) error {
	jobID := c.Params("id")
	if jobID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Job ID is required")
	}

	// TODO: Update job status to completed

	return nil
}

func (h *Handler) FailJobHandler(c *fiber.Ctx) error {
	jobID := c.Params("id")
	if jobID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Job ID is required")
	}

	// TODO: Update job status to failed

	return nil
}
