package admin

import (
	"github.com/PAFFx/job-poll-queue/queue"
)

// QueueService provides additional functionality for queue operations
func (h *Handler) GetQueueStatistics() map[string]interface{} {
	// TODO: Add jobs status count separately
	return map[string]interface{}{
		"size":   h.jobQueue.Size(),
		"status": "operational",
	}
}

// ClearQueue clears both the job queue and results
func (h *Handler) ClearQueue() error {
	// Clear queue
	if err := h.jobQueue.Clear(); err != nil {
		return err
	}

	// Clear results
	if err := h.jobQueue.ClearResults(); err != nil {
		return err
	}

	return nil
}

// GetNextJob retrieves the next job from the queue
func (h *Handler) GetNextJob() (*queue.Message, error) {
	return h.jobQueue.Pop()
}
