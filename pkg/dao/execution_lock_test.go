package dao

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestExecutionLockDAO_Acquire(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewExecutionLockDAO(db, 30*time.Minute)

	// Delete expired locks
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `execution_lock` WHERE resource_id = ? AND expires_at < ?")).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	// Create new lock
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `execution_lock`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := dao.Acquire(100, "task-1")
	assert.NoError(t, err)
}

func TestExecutionLockDAO_Release(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewExecutionLockDAO(db, 30*time.Minute)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `execution_lock` WHERE resource_id = ?")).
		WithArgs(int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Release(100)
	assert.NoError(t, err)
}

func TestExecutionLockDAO_Get(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewExecutionLockDAO(db, 30*time.Minute)

	rows := sqlmock.NewRows([]string{"id", "resource_id", "task_id", "status"}).
		AddRow(1, 100, "task-1", "running")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `execution_lock` WHERE resource_id = ?")).
		WithArgs(int64(100)).
		WillReturnRows(rows)

	lock, err := dao.Get(100)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), lock.ResourceID)
	assert.Equal(t, "task-1", lock.TaskID)
}

func TestExecutionLockDAO_IsLocked(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewExecutionLockDAO(db, 30*time.Minute)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `execution_lock`")).
		WillReturnRows(rows)

	isLocked := dao.IsLocked(100)
	assert.True(t, isLocked)
}

func TestExecutionLockDAO_UpdateStatus(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewExecutionLockDAO(db, 30*time.Minute)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `execution_lock` SET `status`=? WHERE resource_id = ?")).
		WithArgs("completed", int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.UpdateStatus(100, "completed")
	assert.NoError(t, err)
}

func TestExecutionLockDAO_Extend(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewExecutionLockDAO(db, 30*time.Minute)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `execution_lock` SET `expires_at`=? WHERE resource_id = ?")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Extend(100)
	assert.NoError(t, err)
}

func TestExecutionLockDAO_CleanExpired(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewExecutionLockDAO(db, 30*time.Minute)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `execution_lock` WHERE expires_at < ?")).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	err := dao.CleanExpired()
	assert.NoError(t, err)
}
