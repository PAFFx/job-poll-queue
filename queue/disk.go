package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DiskStorage handles persistence operations for the queue
type DiskStorage struct {
	queuePath  string
	resultPath string
}

// NewDiskStorage creates a new disk storage manager
func NewDiskStorage(name string, storageDir string) (*DiskStorage, error) {
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &DiskStorage{
		queuePath:  filepath.Join(storageDir, fmt.Sprintf("%s.json", name)),
		resultPath: filepath.Join(storageDir, fmt.Sprintf("%s-results.json", name)),
	}, nil
}

// SaveQueue persists the queue messages to disk
func (d *DiskStorage) SaveQueue(messages []Message) error {
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

// LoadQueue loads the queue state from disk
func (d *DiskStorage) LoadQueue() ([]Message, error) {
	var messages []Message

	// Check if file exists
	_, err := os.Stat(d.queuePath)
	if os.IsNotExist(err) {
		// No file yet, return empty slice
		return []Message{}, nil
	}

	data, err := os.ReadFile(d.queuePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read queue data from disk: %w", err)
	}

	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal queue data: %w", err)
	}

	return messages, nil
}

// SaveResults persists the results to disk
func (d *DiskStorage) SaveResults(results map[string]Message) error {
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

// LoadResults loads the results from disk
func (d *DiskStorage) LoadResults() (map[string]Message, error) {
	results := make(map[string]Message)

	// Check if file exists
	_, err := os.Stat(d.resultPath)
	if os.IsNotExist(err) {
		// No file yet, return empty map
		return results, nil
	}

	data, err := os.ReadFile(d.resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read results data from disk: %w", err)
	}

	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal results data: %w", err)
	}

	return results, nil
}
