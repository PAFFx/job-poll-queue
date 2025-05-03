package submit

import (
	"github.com/PAFFx/job-poll-queue/queue"
	"github.com/google/uuid"
)

// SubmitJobSync adds a job to the queue and waits for its result
func (h *Handler) SubmitJobSync(payload string, headers map[string]string) (*queue.Message, error) {
	// Create a job
	jobID := uuid.New().String()

	// Create message with provided payload and headers
	msg := queue.Message{
		ID:      jobID,
		Payload: payload,
		Headers: headers,
	}

	// Add job to queue
	if err := h.jobQueue.Push(msg); err != nil {
		return nil, err
	}

	// Wait indefinitely for the result
	return h.jobQueue.GetStatusManager().WaitForCompletionWithoutTimeout(jobID)
}
