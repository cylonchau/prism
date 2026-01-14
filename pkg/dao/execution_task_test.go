package dao

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	models "github.com/cylonchau/prism/pkg/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}

	cleanup := func() {
		sqlDB.Close()
	}

	return gormDB, mock, cleanup
}

func TestExecutionTaskDAO_Create(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := &ExecutionTaskDAO{db: db}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `execution_task`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	task, err := dao.Create("task-1", 100, "apply")
	assert.NoError(t, err)
	assert.Equal(t, "task-1", task.TaskID)
	assert.Equal(t, int64(100), task.ResourceID)
	assert.Equal(t, models.TaskStatusPending, task.Status)
}

func TestExecutionTaskDAO_Get(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := &ExecutionTaskDAO{db: db}

	rows := sqlmock.NewRows([]string{"id", "task_id", "resource_id", "action", "status"}).
		AddRow(1, "task-1", 100, "apply", models.TaskStatusPending)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `execution_task` WHERE task_id = ?")).
		WithArgs("task-1").
		WillReturnRows(rows)

	task, err := dao.Get("task-1")
	assert.NoError(t, err)
	assert.Equal(t, "task-1", task.TaskID)
}

func TestExecutionTaskDAO_Start(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := &ExecutionTaskDAO{db: db}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `execution_task` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Start("task-1")
	assert.NoError(t, err)
}

func TestExecutionTaskDAO_Complete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := &ExecutionTaskDAO{db: db}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `execution_task` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Complete("task-1", true, "output", "")
	assert.NoError(t, err)
}

func TestExecutionTaskDAO_Reset(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := &ExecutionTaskDAO{db: db}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `execution_task` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Reset("task-1")
	assert.NoError(t, err)
}

func TestExecutionTaskDAO_CanRetry(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := &ExecutionTaskDAO{db: db}

	rows := sqlmock.NewRows([]string{"id", "task_id", "status"}).
		AddRow(1, "task-1", models.TaskStatusFailed)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `execution_task` WHERE task_id = ?")).
		WithArgs("task-1").
		WillReturnRows(rows)

	canRetry, err := dao.CanRetry("task-1")
	assert.NoError(t, err)
	assert.True(t, canRetry)
}

func TestExecutionTaskDAO_ListByResource(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := &ExecutionTaskDAO{db: db}

	rows := sqlmock.NewRows([]string{"id", "task_id", "resource_id"}).
		AddRow(1, "task-1", 100).
		AddRow(2, "task-2", 100)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `execution_task` WHERE resource_id = ?")).
		WithArgs(int64(100)).
		WillReturnRows(rows)

	tasks, err := dao.ListByResource(100)
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
}

func TestExecutionTaskDAO_Delete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := &ExecutionTaskDAO{db: db}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `execution_task` WHERE task_id = ?")).
		WithArgs("task-1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Delete("task-1")
	assert.NoError(t, err)
}

func TestExecutionTaskDAO_Complete_Error(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := &ExecutionTaskDAO{db: db}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `execution_task` SET")).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err := dao.Complete("task-1", true, "output", "")
	assert.Error(t, err)
}
