package submit

import (
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
	router.Post("/", h.SubmitJobHandler)
}

func (h *Handler) SubmitJobHandler(c *fiber.Ctx) error {
	// Parse job data from request
	// Add job to queue
	// check if job was process and return response successfully
	// response with error if job fails

	jobData := c.Body()

	// Parse job data from request
	if err := c.BodyParser(jobData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Job submitted successfully",
	})
}
