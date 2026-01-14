// Package lock provides execution lock management.
package lock

import (
	"context"
	"time"
)

// LockStatus represents the status of a lock.
type LockStatus struct {
	ResourceID int64
	TaskID     string
	Status     string // running/completed/failed
	LockedAt   time.Time
	ExpiresAt  time.Time
}

// LockManager defines the lock manager interface.
type LockManager interface {
	// Acquire acquires a lock.
	Acquire(ctx context.Context, resourceID int64, taskID string) error

	// Release releases a lock.
	Release(resourceID int64) error

	// GetStatus returns the lock status.
	GetStatus(resourceID int64) (*LockStatus, error)

	// IsLocked checks if locked.
	IsLocked(resourceID int64) bool
}

// LockType defines the type of lock.
type LockType int

const (
	LockTypeMemory LockType = iota
	LockTypeDB
)

// Config holds lock configuration.
type Config struct {
	Type       LockType
	ExpireTime time.Duration
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Type:       LockTypeMemory,
		ExpireTime: 30 * time.Minute,
	}
}
