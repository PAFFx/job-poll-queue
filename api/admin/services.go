package admin

import (
	"github.com/PAFFx/job-poll-queue/queue"
)

// QueueService provides additional functionality for queue operations
func (h *Handler) GetQueueStatistics() map[string]interface{} {
	return map[string]interface{}{
		"total":     h.jobQueue.GetStatusManager().CountTotalJobs(),
		"pending":   h.jobQueue.GetStatusManager().CountPendingJobs(),
		"running":   h.jobQueue.GetStatusManager().CountProcessingJobs(),
		"completed": h.jobQueue.GetStatusManager().CountCompletedJobs(),
		"failed":    h.jobQueue.GetStatusManager().CountFailedJobs(),
	}
}

// ClearQueue clears both the job queue and results
func (h *Handler) ClearQueue() error {
	// Clear queue
	if err := h.jobQueue.Clear(); err != nil {
		return err
	}

	// Clear job status
	if err := h.jobQueue.GetStatusManager().Clear(); err != nil {
		return err
	}

	return nil
}

// GetNextJob retrieves the next job from the queue
func (h *Handler) GetNextJob() (*queue.Message, error) {
	return h.jobQueue.Pop()
}
