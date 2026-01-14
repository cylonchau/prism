package lock

import (
	"context"
	"testing"
	"time"
)

func TestMemoryLocker_Acquire(t *testing.T) {
	locker := NewMemoryLocker(&Config{
		Type:       LockTypeMemory,
		ExpireTime: 5 * time.Second,
	})

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

func TestMemoryLocker_AcquireExpired(t *testing.T) {
	locker := NewMemoryLocker(&Config{
		Type:       LockTypeMemory,
		ExpireTime: 50 * time.Millisecond,
	})

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

func TestMemoryLocker_Release(t *testing.T) {
	locker := NewMemoryLocker(nil)
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

func TestMemoryLocker_IsLocked(t *testing.T) {
	locker := NewMemoryLocker(nil)
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

func TestMemoryLocker_IsLockedExpired(t *testing.T) {
	locker := NewMemoryLocker(&Config{
		ExpireTime: 50 * time.Millisecond,
	})
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

func TestMemoryLocker_GetStatus(t *testing.T) {
	locker := NewMemoryLocker(nil)
	ctx := context.Background()

	// 未锁定时返回 nil
	status, err := locker.GetStatus(1)
	if err != nil || status != nil {
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

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Type != LockTypeMemory {
		t.Errorf("default type should be memory, got %d", cfg.Type)
	}
	if cfg.ExpireTime != 30*time.Minute {
		t.Errorf("default expire time should be 30 minutes, got %v", cfg.ExpireTime)
	}
}

func TestNewMemoryLocker_NilConfig(t *testing.T) {
	locker := NewMemoryLocker(nil)
	if locker.config == nil {
		t.Fatal("config should be default when nil passed")
	}
}
