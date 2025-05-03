package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Storage handles persistence operations for the queue
type Storage struct {
	queuePath  string
	resultPath string
}

// NewStorage creates a new storage manager
func NewStorage(name string, storageDir string) (*Storage, error) {
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &Storage{
		queuePath:  filepath.Join(storageDir, fmt.Sprintf("%s.json", name)),
		resultPath: filepath.Join(storageDir, fmt.Sprintf("%s-results.json", name)),
	}, nil
}

// SaveQueue persists the queue messages to storage
func (d *Storage) SaveQueue(messages []Message) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal queue data: %w", err)
	}

	// Write to temporary file first to avoid corruption
	tempFile := d.queuePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write queue data to temp file: %w", err)
	}

	// Rename temp file to actual file (atomic operation)
	if err := os.Rename(tempFile, d.queuePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// LoadQueue loads the queue state from storage
func (d *Storage) LoadQueue() ([]Message, error) {
	var messages []Message

	// Check if file exists
	_, err := os.Stat(d.queuePath)
	if os.IsNotExist(err) {
		// No file yet, return empty slice
		return []Message{}, nil
	}

	data, err := os.ReadFile(d.queuePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read queue data from storage: %w", err)
	}

	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal queue data: %w", err)
	}

	return messages, nil
}

// SaveResults persists the results to storage
func (d *Storage) SaveResults(results map[string]Message) error {
	data, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal results data: %w", err)
	}

	// Write to temporary file first to avoid corruption
	tempFile := d.resultPath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write results data to temp file: %w", err)
	}

	// Rename temp file to actual file (atomic operation)
	if err := os.Rename(tempFile, d.resultPath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// LoadResults loads the results from storage
func (d *Storage) LoadResults() (map[string]Message, error) {
	results := make(map[string]Message)

	// Check if file exists
	_, err := os.Stat(d.resultPath)
	if os.IsNotExist(err) {
		// No file yet, return empty map
		return results, nil
	}

	data, err := os.ReadFile(d.resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read results data from storage: %w", err)
	}

	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal results data: %w", err)
	}

	return results, nil
}
