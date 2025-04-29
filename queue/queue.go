package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Message represents an item in the queue
type Message struct {
	ID      string            `json:"id"`
	Payload string            `json:"payload"`
	Headers map[string]string `json:"headers,omitempty"`
}

// Queue implements a FIFO queue with disk persistence
type Queue struct {
	name      string
	messages  []Message
	storePath string
	mutex     sync.Mutex
}

// NewQueue creates a new queue with the given name and storage directory
func NewQueue(name string, storageDir string) (*Queue, error) {
	q := &Queue{
		name:      name,
		messages:  []Message{},
		storePath: filepath.Join(storageDir, fmt.Sprintf("%s.json", name)),
		mutex:     sync.Mutex{},
	}

	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Load existing messages from disk if file exists
	if _, err := os.Stat(q.storePath); err == nil {
		if err := q.loadFromDisk(); err != nil {
			return nil, fmt.Errorf("failed to load queue from disk: %w", err)
		}
	}

	return q, nil
}

// Push adds a message to the end of the queue and persists to disk
func (q *Queue) Push(msg Message) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.messages = append(q.messages, msg)
	return q.saveToDisk()
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

	// Remove it from the queue
	q.messages = q.messages[1:]

	// Save the updated queue state
	if err := q.saveToDisk(); err != nil {
		return nil, fmt.Errorf("failed to persist queue after pop: %w", err)
	}

	return &msg, nil
}

// Size returns the current number of messages in the queue
func (q *Queue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.messages)
}

// saveToDisk persists the queue to disk
func (q *Queue) saveToDisk() error {
	data, err := json.Marshal(q.messages)
	if err != nil {
		return fmt.Errorf("failed to marshal queue data: %w", err)
	}

	// Write to temporary file first to avoid corruption
	tempFile := q.storePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write queue data to temp file: %w", err)
	}

	// Rename temp file to actual file (atomic operation)
	if err := os.Rename(tempFile, q.storePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// loadFromDisk loads the queue state from disk
func (q *Queue) loadFromDisk() error {
	data, err := os.ReadFile(q.storePath)
	if err != nil {
		return fmt.Errorf("failed to read queue data from disk: %w", err)
	}

	if err := json.Unmarshal(data, &q.messages); err != nil {
		return fmt.Errorf("failed to unmarshal queue data: %w", err)
	}

	return nil
}

// Clear removes all messages from the queue
func (q *Queue) Clear() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.messages = []Message{}
	return q.saveToDisk()
}
