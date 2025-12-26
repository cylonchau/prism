package terraform

import (
	"context"
	"testing"

	"github.com/cylonchau/prism/pkg/executor"
	"github.com/cylonchau/prism/pkg/executor/lock"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig should not return nil")
	}
	if cfg.BinaryPath != "terraform" {
		t.Errorf("BinaryPath should be 'terraform', got %s", cfg.BinaryPath)
	}
	if cfg.BasePath != "/tmp/terraform" {
		t.Errorf("BasePath should be '/tmp/terraform', got %s", cfg.BasePath)
	}
}

func TestNew(t *testing.T) {
	locker := lock.NewMemoryLocker(nil)
	exec := New(nil, locker, nil, nil)

	if exec == nil {
		t.Fatal("New should not return nil")
	}
	if exec.config == nil {
		t.Fatal("config should be initialized")
	}
	if exec.locker == nil {
		t.Fatal("locker should be set")
	}
}

func TestNew_WithConfig(t *testing.T) {
	cfg := &Config{
		BinaryPath: "/usr/local/bin/terraform",
		BasePath:   "/var/terraform",
	}
	exec := New(cfg, nil, nil, nil)

	if exec.config.BinaryPath != "/usr/local/bin/terraform" {
		t.Errorf("BinaryPath should be set, got %s", exec.config.BinaryPath)
	}
}

func TestExecutor_Type(t *testing.T) {
	exec := New(nil, nil, nil, nil)
	if exec.Type() != "terraform" {
		t.Errorf("Type should be 'terraform', got %s", exec.Type())
	}
}

func TestExecutor_Validate(t *testing.T) {
	exec := New(nil, nil, nil, nil)
	err := exec.Validate("test config")
	if err != nil {
		t.Errorf("Validate should return nil, got %v", err)
	}
}

func TestExecutor_Execute_UnsupportedAction(t *testing.T) {
	exec := New(nil, nil, nil, nil)

	req := &executor.ExecuteRequest{
		TaskID:     "test-task",
		ResourceID: 1,
		Action:     "invalid",
	}

	result, err := exec.Execute(context.Background(), req)

	if err == nil {
		t.Error("unsupported action should return error")
	}
	if result.Status != executor.StatusFailed {
		t.Errorf("status should be failed, got %s", result.Status)
	}
}

func TestExecutor_Execute_LockFailed(t *testing.T) {
	locker := lock.NewMemoryLocker(nil)
	exec := New(nil, locker, nil, nil)

	locker.Acquire(context.Background(), 1, "other-task")

	req := &executor.ExecuteRequest{
		TaskID:     "test-task",
		ResourceID: 1,
		Action:     executor.ActionApply,
	}

	result, err := exec.Execute(context.Background(), req)

	if err == nil {
		t.Error("locked resource should return error")
	}
	if result.Status != executor.StatusFailed {
		t.Errorf("status should be failed, got %s", result.Status)
	}
}

func TestExecutor_sendMethods(t *testing.T) {
	exec := New(nil, nil, nil, nil)

	exec.sendLog("task-1", "test")
	exec.sendProgress("task-1", "init", 10, "testing")
	exec.sendComplete("task-1", true, nil)
}

func TestExecutor_GetProgress(t *testing.T) {
	exec := New(nil, nil, nil, nil)

	exec.UpdateProgress("test-phase", 50, "testing")

	p := exec.GetProgress()
	if p.Phase != "test-phase" {
		t.Errorf("phase should be 'test-phase', got %s", p.Phase)
	}
	if p.Percent != 50 {
		t.Errorf("percent should be 50, got %d", p.Percent)
	}
}

func TestExecutor_Cancel(t *testing.T) {
	exec := New(nil, nil, nil, nil)

	exec.Transition("start")
	err := exec.Cancel()
	if err != nil {
		t.Errorf("cancel should succeed: %v", err)
	}
}
