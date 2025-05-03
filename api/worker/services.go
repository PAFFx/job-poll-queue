package worker

import (
	"errors"

	"github.com/PAFFx/job-poll-queue/queue"
	"github.com/gofiber/fiber/v2"
)

// RequestJob pops a job from the queue for worker processing
func (h *Handler) RequestJob() (*queue.Message, error) {
	job, err := h.jobQueue.Pop()
	if err != nil {
		return nil, err
	}

	return job, nil
}

// CompleteJob marks a job as completed with the given payload
func (h *Handler) CompleteJob(jobID string, payload string) error {
	// Just mark the job as completed with the provided payload
	// Let the payload itself contain any error information if needed
	return h.jobQueue.GetStatusManager().SubmitResult(jobID, payload, nil)
}

// FormatJobResponse formats a job for response to workers
func (h *Handler) FormatJobResponse(job *queue.Message) fiber.Map {
	return fiber.Map{
		"job_id":  job.ID,
		"payload": job.Payload,
		"headers": job.Headers,
	}
}

// ValidateJobID checks if a job ID is valid
func (h *Handler) ValidateJobID(jobID string) error {
	if jobID == "" {
		return errors.New("job ID is required")
	}
	return nil
}
