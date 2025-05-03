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

// Queue implements a FIFO queue with storage persistence
type Queue struct {
	name      string
	messages  []Message
	storage   *Storage
	mutex     sync.Mutex
	statusMgr *JobStatusManager
}

// NewQueue creates a new queue with the given name and storage directory
func NewQueue(name string, storageDir string) (*Queue, error) {
	// Initialize storage storage
	storage, err := NewStorage(name, storageDir)
	if err != nil {
		return nil, err
	}

	// Initialize job status manager
	statusMgr, err := NewJobStatusManager(storage)
	if err != nil {
		return nil, err
	}

	q := &Queue{
		name:      name,
		messages:  []Message{},
		storage:   storage,
		mutex:     sync.Mutex{},
		statusMgr: statusMgr,
	}

	// Load existing messages from storage
	messages, err := storage.LoadQueue()
	if err != nil {
		return nil, err
	}
	q.messages = messages

	return q, nil
}

func (q *Queue) GetStatusManager() *JobStatusManager {
	return q.statusMgr
}

func (q *Queue) GetStorage() *Storage {
	return q.storage
}

// Push adds a message to the end of the queue and persists to storage
func (q *Queue) Push(msg Message) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Initialize job fields
	now := time.Now()
	msg.Status = JobStatusPending
	msg.CreatedAt = now
	msg.UpdatedAt = now

	q.messages = append(q.messages, msg)

	// Register job in status manager
	q.statusMgr.RegisterJob(&msg)

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

	// Update status in status manager
	q.statusMgr.UpdateStatus(msg.ID, JobStatusProcessing)

	// Remove it from the queue
	q.messages = q.messages[1:]

	// Save the updated queue state
	if err := q.storage.SaveQueue(q.messages); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Clear removes all messages from the queue
func (q *Queue) Clear() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.messages = []Message{}
	return q.storage.SaveQueue(q.messages)
}
