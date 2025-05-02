package submit

import (
	"github.com/PAFFx/job-poll-queue/queue"
	"github.com/google/uuid"
)

func (h *Handler) SubmitJob(jobData string) error {
	// Add job to queue
	err := h.jobQueue.Push(queue.Message{
		ID:      uuid.New().String(),
		Payload: jobData,
		Headers: map[string]string{},
	})

	if err != nil {
		return err
	}

	return nil
}
