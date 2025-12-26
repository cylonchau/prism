package lock

import (
	"context"
	"fmt"
	"time"

	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// DBLocker implements database-based locking.
type DBLocker struct {
	db     *gorm.DB
	config *Config
}

// NewDBLocker creates a new database locker.
func NewDBLocker(db *gorm.DB, config *Config) *DBLocker {
	if config == nil {
		config = DefaultConfig()
	}

	db.AutoMigrate(&models.ExecutionLock{})
	return &DBLocker{
		db:     db,
		config: config,
	}
}

// Acquire acquires a lock for the resource.
func (d *DBLocker) Acquire(ctx context.Context, resourceID int64, taskID string) error {
	now := time.Now()

	d.db.Where("expires_at < ?", now).Delete(&models.ExecutionLock{})

	var existing models.ExecutionLock
	result := d.db.Where("resource_id = ? AND expires_at > ?", resourceID, now).First(&existing)
	if result.Error == nil {
		return fmt.Errorf("resource %d is locked by task %s", resourceID, existing.TaskID)
	}

	lock := models.ExecutionLock{
		ResourceID: resourceID,
		TaskID:     taskID,
		Status:     "running",
		LockedAt:   now,
		ExpiresAt:  now.Add(d.config.ExpireTime),
	}

	result = d.db.Where("resource_id = ?", resourceID).Assign(lock).FirstOrCreate(&lock)
	if result.Error != nil {
		return fmt.Errorf("failed to acquire lock: %w", result.Error)
	}
	return nil
}

// Release releases the lock.
func (d *DBLocker) Release(resourceID int64) error {
	result := d.db.Where("resource_id = ?", resourceID).Delete(&models.ExecutionLock{})
	return result.Error
}

// GetStatus returns the lock status.
func (d *DBLocker) GetStatus(resourceID int64) (*LockStatus, error) {
	var lock models.ExecutionLock
	result := d.db.Where("resource_id = ?", resourceID).First(&lock)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &LockStatus{
		ResourceID: lock.ResourceID,
		TaskID:     lock.TaskID,
		Status:     lock.Status,
		LockedAt:   lock.LockedAt,
		ExpiresAt:  lock.ExpiresAt,
	}, nil
}

// IsLocked checks if the resource is locked.
func (d *DBLocker) IsLocked(resourceID int64) bool {
	var count int64
	d.db.Model(&models.ExecutionLock{}).Where("resource_id = ? AND expires_at > ?", resourceID, time.Now()).Count(&count)
	return count > 0
}
