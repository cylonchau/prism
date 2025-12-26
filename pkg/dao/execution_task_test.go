package dao

import (
	"testing"

	models "github.com/cylonchau/prism/pkg/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	return db
}

func TestExecutionTaskDAO_Create(t *testing.T) {
	db := setupTestDB(t)
	dao := NewExecutionTaskDAO(db)

	task, err := dao.Create("task-1", 100, "apply")
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}
	if task.TaskID != "task-1" {
		t.Errorf("taskID should be 'task-1', got %s", task.TaskID)
	}
	if task.Status != models.TaskStatusPending {
		t.Errorf("status should be pending, got %d", task.Status)
	}
}

func TestExecutionTaskDAO_Get(t *testing.T) {
	db := setupTestDB(t)
	dao := NewExecutionTaskDAO(db)

	dao.Create("task-1", 100, "apply")

	task, err := dao.Get("task-1")
	if err != nil {
		t.Fatalf("get should succeed: %v", err)
	}
	if task.TaskID != "task-1" {
		t.Errorf("taskID should be 'task-1', got %s", task.TaskID)
	}
}

func TestExecutionTaskDAO_Start(t *testing.T) {
	db := setupTestDB(t)
	dao := NewExecutionTaskDAO(db)

	dao.Create("task-1", 100, "apply")
	err := dao.Start("task-1")
	if err != nil {
		t.Fatalf("start should succeed: %v", err)
	}

	task, _ := dao.Get("task-1")
	if task.Status != models.TaskStatusRunning {
		t.Errorf("status should be running, got %d", task.Status)
	}
}

func TestExecutionTaskDAO_Complete(t *testing.T) {
	db := setupTestDB(t)
	dao := NewExecutionTaskDAO(db)

	dao.Create("task-1", 100, "apply")
	dao.Start("task-1")
	err := dao.Complete("task-1", true, "output", "")
	if err != nil {
		t.Fatalf("complete should succeed: %v", err)
	}

	task, _ := dao.Get("task-1")
	if task.Status != models.TaskStatusSuccess {
		t.Errorf("status should be success, got %d", task.Status)
	}
}

func TestExecutionTaskDAO_Reset(t *testing.T) {
	db := setupTestDB(t)
	dao := NewExecutionTaskDAO(db)

	dao.Create("task-1", 100, "apply")
	dao.Start("task-1")
	dao.Complete("task-1", false, "output", "error")

	err := dao.Reset("task-1")
	if err != nil {
		t.Fatalf("reset should succeed: %v", err)
	}

	task, _ := dao.Get("task-1")
	if task.Status != models.TaskStatusPending {
		t.Errorf("status should be pending after reset, got %d", task.Status)
	}
}

func TestExecutionTaskDAO_CanRetry(t *testing.T) {
	db := setupTestDB(t)
	dao := NewExecutionTaskDAO(db)

	dao.Create("task-1", 100, "apply")
	dao.Start("task-1")
	dao.Complete("task-1", false, "", "error")

	canRetry, _ := dao.CanRetry("task-1")
	if !canRetry {
		t.Error("failed task should be retryable")
	}
}

func TestExecutionTaskDAO_ListByResource(t *testing.T) {
	db := setupTestDB(t)
	dao := NewExecutionTaskDAO(db)

	dao.Create("task-1", 100, "apply")
	dao.Create("task-2", 100, "destroy")
	dao.Create("task-3", 200, "apply")

	tasks, err := dao.ListByResource(100)
	if err != nil {
		t.Fatalf("list should succeed: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("should have 2 tasks, got %d", len(tasks))
	}
}

func TestExecutionTaskDAO_Delete(t *testing.T) {
	db := setupTestDB(t)
	dao := NewExecutionTaskDAO(db)

	dao.Create("task-1", 100, "apply")
	err := dao.Delete("task-1")
	if err != nil {
		t.Fatalf("delete should succeed: %v", err)
	}

	_, err = dao.Get("task-1")
	if err == nil {
		t.Error("task should be deleted")
	}
}
