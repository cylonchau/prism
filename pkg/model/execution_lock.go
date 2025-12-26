package models

import "time"

// ExecutionLock 执行锁（db锁）
type ExecutionLock struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceID int64     `gorm:"uniqueIndex;not null" json:"resource_id"`
	TaskID     string    `gorm:"size:64;not null" json:"task_id"`
	Status     string    `gorm:"size:20;not null;default:'running'" json:"status"`
	LockedAt   time.Time `gorm:"not null" json:"locked_at"`
	ExpiresAt  time.Time `gorm:"not null;index" json:"expires_at"`
}

func (ExecutionLock) TableName() string {
	return "execution_lock"
}
