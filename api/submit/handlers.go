package submit

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/PAFFx/job-poll-queue/queue"
)

type Handler struct {
	jobQueue *queue.Queue
}

func NewHandler(jobQueue *queue.Queue) *Handler {
	return &Handler{jobQueue: jobQueue}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/", h.SubmitJobSyncHandler) // Only synchronous job submission
}

// SubmitJobSyncHandler handles synchronous job submission
func (h *Handler) SubmitJobSyncHandler(c *fiber.Ctx) error {
	// Parse job data from request
	type JobRequest struct {
		Payload string            `json:"payload"`
		Headers map[string]string `json:"headers,omitempty"`
		Timeout int               `json:"timeout_seconds,omitempty"` // Optional timeout in seconds
	}

	var jobReq JobRequest
	if err := c.BodyParser(&jobReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid job data")
	}

	// Set timeout
	var timeout time.Duration
	if jobReq.Timeout > 0 {
		timeout = time.Duration(jobReq.Timeout) * time.Second
	}

	// Submit job and wait for result
	result, err := h.SubmitJobSync(jobReq.Payload, jobReq.Headers, timeout)
	if err != nil {
		if err.Error() == "timeout waiting for job result" {
			return fiber.NewError(fiber.StatusRequestTimeout, "Job processing timed out")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process job: "+err.Error())
	}

	// Check if job failed
	if result.Status == queue.JobStatusFailed {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"job_id": result.ID,
			"status": string(result.Status),
			"error":  result.Error,
		})
	}

	// Return the successful result
	return c.JSON(fiber.Map{
		"job_id":       result.ID,
		"status":       string(result.Status),
		"result":       result.Result,
		"created_at":   result.CreatedAt,
		"completed_at": result.CompletedAt,
	})
}
