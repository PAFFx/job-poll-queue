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

// CompleteJob marks a job as completed with the given result
func (h *Handler) CompleteJob(jobID string, result string) error {
	return h.jobQueue.GetStatusManager().SubmitResult(jobID, result, nil)
}

// FailJob marks a job as failed with the given error
func (h *Handler) FailJob(jobID string, errorMessage string) error {
	return h.jobQueue.GetStatusManager().SubmitResult(jobID, "", errors.New(errorMessage))
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

// FormatJobResult formats a result response
func (h *Handler) FormatJobResultResponse(jobID string, success bool) fiber.Map {
	message := "Job completed successfully"
	if !success {
		message = "Job marked as failed"
	}

	return fiber.Map{
		"message": message,
		"job_id":  jobID,
	}
}
