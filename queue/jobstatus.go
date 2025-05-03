package queue

import (
	"errors"
	"sync"
	"time"
)

// JobStatusManager handles job status tracking throughout the job lifecycle
type JobStatusManager struct {
	statusMap map[string]Message // Map job ID to job with current status
	waiters   map[string]chan Message
	storage   *Storage
	mutex     sync.Mutex
}

// NewJobStatusManager creates a new job status manager
func NewJobStatusManager(storage *Storage) (*JobStatusManager, error) {
	jsm := &JobStatusManager{
		statusMap: make(map[string]Message),
		waiters:   make(map[string]chan Message),
		storage:   storage,
		mutex:     sync.Mutex{},
	}

	// Load existing job status information from storage
	jobStatus, err := storage.LoadJobStatus()
	if err != nil {
		return nil, err
	}
	jsm.statusMap = jobStatus

	return jsm, nil
}

// RegisterJob registers a new job for status tracking
func (jsm *JobStatusManager) RegisterJob(job *Message) {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()

	// Create a waiter channel for the job
	jsm.waiters[job.ID] = make(chan Message, 1)

	// Store initial status
	jsm.statusMap[job.ID] = *job
}

// UpdateStatus updates a job's status
func (jsm *JobStatusManager) UpdateStatus(jobID string, status JobStatus) error {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()

	// Find if the job exists
	job, exists := jsm.statusMap[jobID]
	if !exists {
		return errors.New("job not found")
	}

	// Update status
	job.Status = status
	job.UpdatedAt = time.Now()

	// Store the updated job
	jsm.statusMap[jobID] = job

	return jsm.storage.SaveJobStatus(jsm.statusMap)
}

// SubmitResult stores the result of a processed job
func (jsm *JobStatusManager) SubmitResult(jobID string, result string, err error) error {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()

	// Find if the job exists in the status map
	job, exists := jsm.statusMap[jobID]
	if !exists {
		// It's possible the job was just processed but not in statusMap
		// Let's create a new status entry
		now := time.Now()
		job = Message{
			ID:          jobID,
			UpdatedAt:   now,
			CompletedAt: &now,
		}
	}

	// Update with results
	job.Result = result
	now := time.Now()
	job.UpdatedAt = now
	job.CompletedAt = &now

	if err != nil {
		job.Status = JobStatusFailed
		job.Error = err.Error()
	} else {
		job.Status = JobStatusCompleted
	}

	// Store the updated job
	jsm.statusMap[jobID] = job

	// Notify any waiters for this job
	if waiter, ok := jsm.waiters[jobID]; ok {
		// Non-blocking send to the channel
		select {
		case waiter <- job:
			// Message sent successfully
		default:
			// No one is listening, which is fine
		}

		// Delete the waiter channel after use
		delete(jsm.waiters, jobID)
	}

	// Save job status to storage
	return jsm.storage.SaveJobStatus(jsm.statusMap)
}

// WaitForCompletion waits for a job to reach completion with a timeout
func (jsm *JobStatusManager) WaitForCompletion(jobID string, timeout time.Duration) (*Message, error) {
	jsm.mutex.Lock()

	// Check if we already have the job in completed/failed state
	if job, exists := jsm.statusMap[jobID]; exists {
		if job.Status == JobStatusCompleted || job.Status == JobStatusFailed {
			jsm.mutex.Unlock()
			return &job, nil
		}
	}

	// Get the waiter channel for this job
	waiter, exists := jsm.waiters[jobID]
	if !exists {
		jsm.mutex.Unlock()
		return nil, errors.New("job not found")
	}

	jsm.mutex.Unlock()

	// Wait for the job to complete with timeout
	select {
	case job := <-waiter:
		return &job, nil
	case <-time.After(timeout):
		return nil, errors.New("timeout waiting for job completion")
	}
}

// WaitForCompletionWithoutTimeout waits indefinitely for a job to reach completion
func (jsm *JobStatusManager) WaitForCompletionWithoutTimeout(jobID string) (*Message, error) {
	jsm.mutex.Lock()

	// Check if we already have the job in completed/failed state
	if job, exists := jsm.statusMap[jobID]; exists {
		if job.Status == JobStatusCompleted || job.Status == JobStatusFailed {
			jsm.mutex.Unlock()
			return &job, nil
		}
	}

	// Get the waiter channel for this job
	waiter, exists := jsm.waiters[jobID]
	if !exists {
		jsm.mutex.Unlock()
		return nil, errors.New("job not found")
	}

	jsm.mutex.Unlock()

	// Wait indefinitely for the job to complete
	job := <-waiter
	return &job, nil
}

// GetJobStatus retrieves a job's current status
func (jsm *JobStatusManager) GetJobStatus(jobID string) (*Message, error) {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()

	job, exists := jsm.statusMap[jobID]
	if !exists {
		return nil, errors.New("job not found")
	}

	return &job, nil
}

// Count returns the number of tracked jobs
func (jsm *JobStatusManager) CountTotalJobs() int {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()
	return len(jsm.statusMap)
}

func (jsm *JobStatusManager) CountPendingJobs() int {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()
	count := 0
	for _, job := range jsm.statusMap {
		if job.Status == JobStatusPending {
			count++
		}
	}
	return count
}

func (jsm *JobStatusManager) CountProcessingJobs() int {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()
	count := 0
	for _, job := range jsm.statusMap {
		if job.Status == JobStatusProcessing {
			count++
		}
	}
	return count
}

func (jsm *JobStatusManager) CountCompletedJobs() int {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()
	count := 0
	for _, job := range jsm.statusMap {
		if job.Status == JobStatusCompleted {
			count++
		}
	}
	return count
}

func (jsm *JobStatusManager) CountFailedJobs() int {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()
	count := 0
	for _, job := range jsm.statusMap {
		if job.Status == JobStatusFailed {
			count++
		}
	}
	return count
}

// Clear removes all job status records
func (jsm *JobStatusManager) Clear() error {
	jsm.mutex.Lock()
	defer jsm.mutex.Unlock()

	jsm.statusMap = make(map[string]Message)
	return jsm.storage.SaveJobStatus(jsm.statusMap)
}
