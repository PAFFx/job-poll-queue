package worker

import "github.com/PAFFx/job-poll-queue/queue"

func (h *Handler) RequestJob() (*queue.Message, error) {
	job, err := h.jobQueue.Pop()
	if err != nil {
		return nil, err
	}

	return job, nil
}
