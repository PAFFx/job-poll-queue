package queue

import (
	"errors"
	"sync"
	"time"
)

// ResultManager handles job result tracking and retrieval
type ResultManager struct {
	results map[string]Message // Map job ID to message with results
	waiters map[string]chan Message
	storage *Storage
	mutex   sync.Mutex
}

// NewResultManager creates a new result manager
func NewResultManager(storage *Storage) (*ResultManager, error) {
	rm := &ResultManager{
		results: make(map[string]Message),
		waiters: make(map[string]chan Message),
		storage: storage,
		mutex:   sync.Mutex{},
	}

	// Load existing results from disk
	results, err := storage.LoadResults()
	if err != nil {
		return nil, err
	}
	rm.results = results

	return rm, nil
}

// RegisterWaiter creates a waiter channel for a job ID
func (rm *ResultManager) RegisterWaiter(jobID string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	rm.waiters[jobID] = make(chan Message, 1)
}

// SubmitResult stores the result of a processed job
func (rm *ResultManager) SubmitResult(jobID string, result string, err error) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Find if the job exists in the results map
	msg, exists := rm.results[jobID]
	if !exists {
		// It's possible the job was just popped but not yet in results
		// Let's create a new result entry
		now := time.Now()
		msg = Message{
			ID:          jobID,
			Status:      JobStatusCompleted,
			UpdatedAt:   now,
			CompletedAt: &now,
		}
	}

	// Update with results
	msg.Result = result
	if err != nil {
		msg.Status = JobStatusFailed
		msg.Error = err.Error()
	} else {
		msg.Status = JobStatusCompleted
	}

	now := time.Now()
	msg.UpdatedAt = now
	msg.CompletedAt = &now

	// Store the result
	rm.results[jobID] = msg

	// Notify any waiters for this job
	if waiter, ok := rm.waiters[jobID]; ok {
		// Non-blocking send to the channel
		select {
		case waiter <- msg:
			// Message sent successfully
		default:
			// No one is listening, which is fine
		}

		// Delete the waiter channel after use
		delete(rm.waiters, jobID)
	}

	// Save results to disk
	return rm.storage.SaveResults(rm.results)
}

// WaitForResult waits for a job result with a timeout
func (rm *ResultManager) WaitForResult(jobID string, timeout time.Duration) (*Message, error) {
	rm.mutex.Lock()

	// Check if we already have the result
	if result, exists := rm.results[jobID]; exists {
		rm.mutex.Unlock()
		return &result, nil
	}

	// Get the waiter channel for this job
	waiter, exists := rm.waiters[jobID]
	if !exists {
		rm.mutex.Unlock()
		return nil, errors.New("job not found")
	}

	rm.mutex.Unlock()

	// Wait for the result with timeout
	select {
	case result := <-waiter:
		return &result, nil
	case <-time.After(timeout):
		return nil, errors.New("timeout waiting for job result")
	}
}

// GetResult retrieves a job result if available
func (rm *ResultManager) GetResult(jobID string) (*Message, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	result, exists := rm.results[jobID]
	if !exists {
		return nil, errors.New("result not found")
	}

	return &result, nil
}

// Count returns the number of completed job results
func (rm *ResultManager) Count() int {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	return len(rm.results)
}

// Clear removes all stored results
func (rm *ResultManager) Clear() error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	rm.results = make(map[string]Message)
	return rm.storage.SaveResults(rm.results)
}
