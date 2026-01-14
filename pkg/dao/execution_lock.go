package dao

import (
	"time"

	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// ExecutionLockDAO provides lock data access operations.
type ExecutionLockDAO struct {
	db         *gorm.DB
	expireTime time.Duration
}

// NewExecutionLockDAO creates a new lock DAO.
func NewExecutionLockDAO(db *gorm.DB, expireTime time.Duration) *ExecutionLockDAO {
	db.AutoMigrate(&models.ExecutionLock{})
	if expireTime == 0 {
		expireTime = 30 * time.Minute
	}
	return &ExecutionLockDAO{db: db, expireTime: expireTime}
}

// Acquire acquires a lock for the resource.
func (d *ExecutionLockDAO) Acquire(resourceID int64, taskID string) error {
	now := time.Now()

	// Delete expired locks
	d.db.Where("resource_id = ? AND expires_at < ?", resourceID, now).
		Delete(&models.ExecutionLock{})

	// Try to create lock
	lock := &models.ExecutionLock{
		ResourceID: resourceID,
		TaskID:     taskID,
		Status:     "running",
		LockedAt:   now,
		ExpiresAt:  now.Add(d.expireTime),
	}

	result := d.db.Create(lock)
	return result.Error
}

// Release releases the lock.
func (d *ExecutionLockDAO) Release(resourceID int64) error {
	return d.db.Where("resource_id = ?", resourceID).
		Delete(&models.ExecutionLock{}).Error
}

// Get retrieves lock by resource ID.
func (d *ExecutionLockDAO) Get(resourceID int64) (*models.ExecutionLock, error) {
	var lock models.ExecutionLock
	result := d.db.Where("resource_id = ?", resourceID).First(&lock)
	if result.Error != nil {
		return nil, result.Error
	}
	return &lock, nil
}

// IsLocked checks if resource is locked.
func (d *ExecutionLockDAO) IsLocked(resourceID int64) bool {
	var count int64
	d.db.Model(&models.ExecutionLock{}).
		Where("resource_id = ? AND expires_at > ?", resourceID, time.Now()).
		Count(&count)
	return count > 0
}

// UpdateStatus updates lock status.
func (d *ExecutionLockDAO) UpdateStatus(resourceID int64, status string) error {
	return d.db.Model(&models.ExecutionLock{}).
		Where("resource_id = ?", resourceID).
		Update("status", status).Error
}

// Extend extends the lock expiration.
func (d *ExecutionLockDAO) Extend(resourceID int64) error {
	return d.db.Model(&models.ExecutionLock{}).
		Where("resource_id = ?", resourceID).
		Update("expires_at", time.Now().Add(d.expireTime)).Error
}

// CleanExpired removes all expired locks.
func (d *ExecutionLockDAO) CleanExpired() error {
	return d.db.Where("expires_at < ?", time.Now()).
		Delete(&models.ExecutionLock{}).Error
}
