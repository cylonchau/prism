package models

import "time"

// TaskStatus represents task execution status.
type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusRunning
	TaskStatusSuccess
	TaskStatusFailed
	TaskStatusCancelled
)

// ExecutionTask stores task execution information.
type ExecutionTask struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID     string     `gorm:"size:64;uniqueIndex;not null" json:"task_id"`
	ResourceID int64      `gorm:"index;not null" json:"resource_id"`
	Action     string     `gorm:"size:20;not null" json:"action"`
	Status     TaskStatus `gorm:"not null;default:0" json:"status"`
	Output     string     `gorm:"type:text" json:"output"`
	Error      string     `gorm:"type:text" json:"error"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Duration   int64      `json:"duration"` // milliseconds
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ExecutionTask) TableName() string {
	return "execution_task"
}

// IsRetryable checks if task can be retried.
func (t *ExecutionTask) IsRetryable() bool {
	return t.Status == TaskStatusFailed || t.Status == TaskStatusCancelled
}
