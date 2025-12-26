package lock

import (
	"context"
	"testing"
	"time"

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

func TestDBLocker_NewDBLocker(t *testing.T) {
	db := setupTestDB(t)
	locker := NewDBLocker(db, nil)

	if locker == nil {
		t.Fatal("NewDBLocker should not return nil")
	}
	if locker.db == nil {
		t.Fatal("db should be set")
	}
	if locker.config == nil {
		t.Fatal("config should use default when nil passed")
	}
}

func TestDBLocker_Acquire(t *testing.T) {
	db := setupTestDB(t)
	locker := NewDBLocker(db, &Config{ExpireTime: 5 * time.Second})
	ctx := context.Background()

	// 第一次获取锁应该成功
	err := locker.Acquire(ctx, 1, "task-1")
	if err != nil {
		t.Fatalf("first acquire should succeed: %v", err)
	}

	// 第二次获取同一资源的锁应该失败
	err = locker.Acquire(ctx, 1, "task-2")
	if err == nil {
		t.Fatal("second acquire should fail")
	}

	// 不同资源应该成功
	err = locker.Acquire(ctx, 2, "task-3")
	if err != nil {
		t.Fatalf("different resource should succeed: %v", err)
	}
}

func TestDBLocker_AcquireExpired(t *testing.T) {
	db := setupTestDB(t)
	locker := NewDBLocker(db, &Config{ExpireTime: 50 * time.Millisecond})
	ctx := context.Background()

	// 获取锁
	err := locker.Acquire(ctx, 1, "task-1")
	if err != nil {
		t.Fatalf("acquire should succeed: %v", err)
	}

	// 等待锁过期
	time.Sleep(100 * time.Millisecond)

	// 再次获取应该成功（因为过期了）
	err = locker.Acquire(ctx, 1, "task-2")
	if err != nil {
		t.Fatalf("acquire after expiry should succeed: %v", err)
	}
}

func TestDBLocker_Release(t *testing.T) {
	db := setupTestDB(t)
	locker := NewDBLocker(db, nil)
	ctx := context.Background()

	// 获取锁
	err := locker.Acquire(ctx, 1, "task-1")
	if err != nil {
		t.Fatalf("acquire should succeed: %v", err)
	}

	// 释放锁
	err = locker.Release(1)
	if err != nil {
		t.Fatalf("release should succeed: %v", err)
	}

	// 再次获取应该成功
	err = locker.Acquire(ctx, 1, "task-2")
	if err != nil {
		t.Fatalf("acquire after release should succeed: %v", err)
	}
}

func TestDBLocker_IsLocked(t *testing.T) {
	db := setupTestDB(t)
	locker := NewDBLocker(db, nil)
	ctx := context.Background()

	// 未锁定
	if locker.IsLocked(1) {
		t.Fatal("should not be locked initially")
	}

	// 获取锁
	locker.Acquire(ctx, 1, "task-1")

	// 已锁定
	if !locker.IsLocked(1) {
		t.Fatal("should be locked after acquire")
	}

	// 释放锁
	locker.Release(1)

	// 未锁定
	if locker.IsLocked(1) {
		t.Fatal("should not be locked after release")
	}
}

func TestDBLocker_IsLockedExpired(t *testing.T) {
	db := setupTestDB(t)
	locker := NewDBLocker(db, &Config{ExpireTime: 50 * time.Millisecond})
	ctx := context.Background()

	locker.Acquire(ctx, 1, "task-1")

	if !locker.IsLocked(1) {
		t.Fatal("should be locked")
	}

	time.Sleep(100 * time.Millisecond)

	if locker.IsLocked(1) {
		t.Fatal("should not be locked after expiry")
	}
}

func TestDBLocker_GetStatus(t *testing.T) {
	db := setupTestDB(t)
	locker := NewDBLocker(db, nil)
	ctx := context.Background()

	// 未锁定时返回 nil
	status, err := locker.GetStatus(1)
	if err != nil {
		t.Fatalf("get status should not error: %v", err)
	}
	if status != nil {
		t.Fatal("should return nil for unlocked resource")
	}

	// 获取锁后返回状态
	locker.Acquire(ctx, 1, "task-1")
	status, err = locker.GetStatus(1)
	if err != nil {
		t.Fatalf("get status should succeed: %v", err)
	}
	if status == nil {
		t.Fatal("status should not be nil")
	}
	if status.TaskID != "task-1" {
		t.Fatalf("task id should be 'task-1', got %s", status.TaskID)
	}
	if status.Status != "running" {
		t.Fatalf("status should be 'running', got %s", status.Status)
	}
}

func TestExecutionLock_TableName(t *testing.T) {
	lock := models.ExecutionLock{}
	if lock.TableName() != "execution_lock" {
		t.Errorf("table name should be 'execution_lock', got %s", lock.TableName())
	}
}
