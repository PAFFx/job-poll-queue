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

	return c.JSON(h.FormatJobResponse(job))
}

func (h *Handler) CompleteJobHandler(c *fiber.Ctx) error {
	jobID := c.Params("id")
	if err := h.ValidateJobID(jobID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Parse result data
	type ResultRequest struct {
		Result string `json:"result"`
	}

	var req ResultRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid result data")
	}

	if err := h.CompleteJob(jobID, req.Result); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to complete job: "+err.Error())
	}

	return c.JSON(h.FormatJobResultResponse(jobID, true))
}

func (h *Handler) FailJobHandler(c *fiber.Ctx) error {
	jobID := c.Params("id")
	if err := h.ValidateJobID(jobID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Parse error data
	type ErrorRequest struct {
		Error string `json:"error"`
	}

	var req ErrorRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid error data")
	}

	if req.Error == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Error message is required")
	}

	if err := h.FailJob(jobID, req.Error); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to mark job as failed: "+err.Error())
	}

	return c.JSON(h.FormatJobResultResponse(jobID, false))
}
