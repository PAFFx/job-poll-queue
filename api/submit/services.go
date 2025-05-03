package submit

import (
	"time"

	"github.com/PAFFx/job-poll-queue/queue"
	"github.com/google/uuid"
)

// DefaultTimeout is the default time to wait for a synchronous job result
const DefaultTimeout = 30 * time.Second

// SubmitJobSync adds a job to the queue and waits for its result
func (h *Handler) SubmitJobSync(payload string, headers map[string]string, timeout time.Duration) (*queue.Message, error) {
	// Use default timeout if not specified
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

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

	// Wait for the result
	return h.jobQueue.GetStatusManager().WaitForCompletion(jobID, timeout)
}
