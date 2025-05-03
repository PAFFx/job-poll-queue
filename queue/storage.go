package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Storage handles persistence operations for the queue
type Storage struct {
	queuePath     string
	jobStatusPath string
}

// NewStorage creates a new storage manager
func NewStorage(name string, storageDir string) (*Storage, error) {
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &Storage{
		queuePath:     filepath.Join(storageDir, fmt.Sprintf("%s.json", name)),
		jobStatusPath: filepath.Join(storageDir, fmt.Sprintf("%s-jobstatus.json", name)),
	}, nil
}

// SaveQueue persists the queue messages to storage
func (s *Storage) SaveQueue(messages []Message) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal queue data: %w", err)
	}

	// Write to temporary file first to avoid corruption
	tempFile := s.queuePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write queue data to temp file: %w", err)
	}

	// Rename temp file to actual file (atomic operation)
	if err := os.Rename(tempFile, s.queuePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// LoadQueue loads the queue state from storage
func (s *Storage) LoadQueue() ([]Message, error) {
	var messages []Message

	// Check if file exists
	_, err := os.Stat(s.queuePath)
	if os.IsNotExist(err) {
		// No file yet, return empty slice
		return []Message{}, nil
	}

	data, err := os.ReadFile(s.queuePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read queue data from storage: %w", err)
	}

	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal queue data: %w", err)
	}

	return messages, nil
}

// SaveJobStatus persists the job status records to storage
func (s *Storage) SaveJobStatus(jobStatus map[string]Message) error {
	data, err := json.Marshal(jobStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal job status data: %w", err)
	}

	// Write to temporary file first to avoid corruption
	tempFile := s.jobStatusPath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write job status data to temp file: %w", err)
	}

	// Rename temp file to actual file (atomic operation)
	if err := os.Rename(tempFile, s.jobStatusPath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// LoadJobStatus loads the job status records from storage
func (s *Storage) LoadJobStatus() (map[string]Message, error) {
	jobStatus := make(map[string]Message)

	// Check if file exists
	_, err := os.Stat(s.jobStatusPath)
	if os.IsNotExist(err) {
		// No file yet, return empty map
		return jobStatus, nil
	}

	data, err := os.ReadFile(s.jobStatusPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read job status data from storage: %w", err)
	}

	if err := json.Unmarshal(data, &jobStatus); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job status data: %w", err)
	}

	return jobStatus, nil
}
