package lock

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryLocker implements in-memory locking.
type MemoryLocker struct {
	mu     sync.Mutex
	locks  sync.Map // map[int64]*LockStatus
	config *Config
}

// NewMemoryLocker creates a new memory locker.
func NewMemoryLocker(config *Config) *MemoryLocker {
	if config == nil {
		config = DefaultConfig()
	}
	return &MemoryLocker{
		config: config,
	}
}

// Acquire acquires a lock for the resource.
func (m *MemoryLocker) Acquire(ctx context.Context, resourceID int64, taskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.locks.Load(resourceID); ok {
		status := existing.(*LockStatus)

		if time.Now().Before(status.ExpiresAt) {
			return fmt.Errorf("resource %d is locked by task %s", resourceID, status.TaskID)
		}
	}

	now := time.Now()
	status := &LockStatus{
		ResourceID: resourceID,
		TaskID:     taskID,
		Status:     "running",
		LockedAt:   now,
		ExpiresAt:  now.Add(m.config.ExpireTime),
	}
	m.locks.Store(resourceID, status)
	return nil
}

// Release releases the lock.
func (m *MemoryLocker) Release(resourceID int64) error {
	m.locks.Delete(resourceID)
	return nil
}

// GetStatus returns the lock status.
func (m *MemoryLocker) GetStatus(resourceID int64) (*LockStatus, error) {
	if existing, ok := m.locks.Load(resourceID); ok {
		return existing.(*LockStatus), nil
	}
	return nil, nil
}

// IsLocked checks if the resource is locked.
func (m *MemoryLocker) IsLocked(resourceID int64) bool {
	if existing, ok := m.locks.Load(resourceID); ok {
		status := existing.(*LockStatus)
		return time.Now().Before(status.ExpiresAt)
	}
	return false
}
