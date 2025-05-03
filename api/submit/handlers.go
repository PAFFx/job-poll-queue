package submit

import (
	"encoding/json"

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
	// Get the raw request body as payload
	rawBody := c.Body()

	// Convert raw body to JSON string for the payload
	payload := string(rawBody)

	// Use standard headers from the request
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	// Submit job and wait for result
	result, err := h.SubmitJobSync(payload, headers)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process job: "+err.Error())
	}

	// Parse the result payload as JSON if possible
	var resultPayload interface{}
	if err := json.Unmarshal([]byte(result.Result), &resultPayload); err != nil {
		// If it's not valid JSON, use the raw string
		resultPayload = result.Result
	}

	// Return the job result
	return c.JSON(fiber.Map{
		"job_id":       result.ID,
		"payload":      resultPayload,
		"created_at":   result.CreatedAt,
		"completed_at": result.CompletedAt,
	})
}
