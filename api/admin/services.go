package admin

// QueueService provides additional functionality for queue operations
func (h *Handler) GetQueueStatistics() map[string]interface{} {
	return map[string]interface{}{
		"size":   h.jobQueue.Size(),
		"status": "operational",
	}
}

func (h *Handler) ClearQueue() error {
	return h.jobQueue.Clear()
}
