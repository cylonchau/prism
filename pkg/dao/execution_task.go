// Package dao provides data access objects.
package dao

import (
	"fmt"
	"time"

	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// ExecutionTaskDAO provides task data access operations.
type ExecutionTaskDAO struct {
	db *gorm.DB
}

// NewExecutionTaskDAO creates a new task DAO.
func NewExecutionTaskDAO(db *gorm.DB) *ExecutionTaskDAO {
	db.AutoMigrate(&models.ExecutionTask{})
	return &ExecutionTaskDAO{db: db}
}

// Create creates a new task.
func (d *ExecutionTaskDAO) Create(taskID string, resourceID int64, action string) (*models.ExecutionTask, error) {
	task := &models.ExecutionTask{
		TaskID:     taskID,
		ResourceID: resourceID,
		Action:     action,
		Status:     models.TaskStatusPending,
	}
	result := d.db.Create(task)
	return task, result.Error
}

// Get retrieves a task by task ID.
func (d *ExecutionTaskDAO) Get(taskID string) (*models.ExecutionTask, error) {
	var task models.ExecutionTask
	result := d.db.Where("task_id = ?", taskID).First(&task)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

// UpdateStatus updates task status.
func (d *ExecutionTaskDAO) UpdateStatus(taskID string, status models.TaskStatus) error {
	return d.db.Model(&models.ExecutionTask{}).
		Where("task_id = ?", taskID).
		Update("status", status).Error
}

// Start marks task as running.
func (d *ExecutionTaskDAO) Start(taskID string) error {
	now := time.Now()
	return d.db.Model(&models.ExecutionTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":     models.TaskStatusRunning,
			"started_at": now,
		}).Error
}

// Complete marks task as completed.
func (d *ExecutionTaskDAO) Complete(taskID string, success bool, output, errMsg string) error {
	now := time.Now()
	status := models.TaskStatusSuccess
	if !success {
		status = models.TaskStatusFailed
	}

	var task models.ExecutionTask
	d.db.Where("task_id = ?", taskID).First(&task)

	var duration int64
	if task.StartedAt != nil {
		duration = now.Sub(*task.StartedAt).Milliseconds()
	}

	return d.db.Model(&models.ExecutionTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":      status,
			"output":      output,
			"error":       errMsg,
			"finished_at": now,
			"duration":    duration,
		}).Error
}

// Reset resets a failed task for retry.
func (d *ExecutionTaskDAO) Reset(taskID string) error {
	return d.db.Model(&models.ExecutionTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":      models.TaskStatusPending,
			"output":      "",
			"error":       "",
			"started_at":  nil,
			"finished_at": nil,
			"duration":    0,
		}).Error
}

// ListByResource lists tasks for a resource.
func (d *ExecutionTaskDAO) ListByResource(resourceID int64) ([]models.ExecutionTask, error) {
	var tasks []models.ExecutionTask
	result := d.db.Where("resource_id = ?", resourceID).Order("created_at DESC").Find(&tasks)
	return tasks, result.Error
}

// ListFailed lists all failed tasks.
func (d *ExecutionTaskDAO) ListFailed() ([]models.ExecutionTask, error) {
	var tasks []models.ExecutionTask
	result := d.db.Where("status = ?", models.TaskStatusFailed).Order("created_at DESC").Find(&tasks)
	return tasks, result.Error
}

// Delete deletes a task.
func (d *ExecutionTaskDAO) Delete(taskID string) error {
	return d.db.Where("task_id = ?", taskID).Delete(&models.ExecutionTask{}).Error
}

// CanRetry checks if task can be retried.
func (d *ExecutionTaskDAO) CanRetry(taskID string) (bool, error) {
	task, err := d.Get(taskID)
	if err != nil {
		return false, fmt.Errorf("task not found: %w", err)
	}
	return task.IsRetryable(), nil
}
