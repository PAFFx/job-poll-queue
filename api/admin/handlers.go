package admin

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
	router.Get("/stats", h.GetQueueStatsHandler)
	router.Delete("/clear", h.ClearQueueHandler)
}

func (h *Handler) GetQueueStatsHandler(c *fiber.Ctx) error {
	return c.JSON(h.GetQueueStatistics())
}

func (h *Handler) ClearQueueHandler(c *fiber.Ctx) error {
	if err := h.ClearQueue(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to clear queue")
	}
	return c.JSON(fiber.Map{"message": "Queue cleared successfully"})
}
