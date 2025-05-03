package queue

import (
	"sync"
	"time"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

// Message represents an item in the queue
type Message struct {
	ID          string            `json:"id"`
	Payload     string            `json:"payload"`
	Headers     map[string]string `json:"headers,omitempty"`
	Status      JobStatus         `json:"status"`
	Result      string            `json:"result,omitempty"`
	Error       string            `json:"error,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
}

// Queue implements a FIFO queue with disk persistence
type Queue struct {
	name      string
	messages  []Message
	storage   *DiskStorage
	mutex     sync.Mutex
	resultMgr *ResultManager
}

// NewQueue creates a new queue with the given name and storage directory
func NewQueue(name string, storageDir string) (*Queue, error) {
	// Initialize disk storage
	storage, err := NewDiskStorage(name, storageDir)
	if err != nil {
		return nil, err
	}

	// Initialize result manager
	resultMgr, err := NewResultManager(storage)
	if err != nil {
		return nil, err
	}

	q := &Queue{
		name:      name,
		messages:  []Message{},
		storage:   storage,
		mutex:     sync.Mutex{},
		resultMgr: resultMgr,
	}

	// Load existing messages from disk
	messages, err := storage.LoadQueue()
	if err != nil {
		return nil, err
	}
	q.messages = messages

	return q, nil
}

// Push adds a message to the end of the queue and persists to disk
func (q *Queue) Push(msg Message) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Initialize job fields
	now := time.Now()
	msg.Status = JobStatusPending
	msg.CreatedAt = now
	msg.UpdatedAt = now

	q.messages = append(q.messages, msg)

	// Register a waiter for this job
	q.resultMgr.RegisterWaiter(msg.ID)

	return q.storage.SaveQueue(q.messages)
}

// Pop removes and returns the first message from the queue
// Returns nil if the queue is empty
func (q *Queue) Pop() (*Message, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.messages) == 0 {
		return nil, nil
	}

	// Get the first message
	msg := q.messages[0]

	// Update status to processing
	msg.Status = JobStatusProcessing
	msg.UpdatedAt = time.Now()

	// Remove it from the queue
	q.messages = q.messages[1:]

	// Save the updated queue state
	if err := q.storage.SaveQueue(q.messages); err != nil {
		return nil, err
	}

	return &msg, nil
}

// SubmitResult stores the result of a processed job
func (q *Queue) SubmitResult(jobID string, result string, err error) error {
	return q.resultMgr.SubmitResult(jobID, result, err)
}

// WaitForResult waits for a job result with a timeout
func (q *Queue) WaitForResult(jobID string, timeout time.Duration) (*Message, error) {
	return q.resultMgr.WaitForResult(jobID, timeout)
}

// GetResult retrieves a job result if available
func (q *Queue) GetResult(jobID string) (*Message, error) {
	return q.resultMgr.GetResult(jobID)
}

// Size returns the current number of messages in the queue
func (q *Queue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.messages)
}

// ResultsCount returns the number of completed job results
func (q *Queue) ResultsCount() int {
	return q.resultMgr.Count()
}

// Clear removes all messages from the queue
func (q *Queue) Clear() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.messages = []Message{}
	return q.storage.SaveQueue(q.messages)
}

// ClearResults removes all stored results
func (q *Queue) ClearResults() error {
	return q.resultMgr.Clear()
}
