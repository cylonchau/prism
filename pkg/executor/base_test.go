package executor

import (
	"context"
	"testing"
)

func TestBaseExecutor_New(t *testing.T) {
	b := NewBaseExecutor()
	if b == nil {
		t.Fatal("NewBaseExecutor should not return nil")
	}
	if b.progress == nil {
		t.Fatal("progress should be initialized")
	}
	if b.fsm == nil {
		t.Fatal("fsm should be initialized")
	}
}

func TestBaseExecutor_Status(t *testing.T) {
	b := NewBaseExecutor()

	// 初始状态应该是 pending
	if b.Status() != StatusPending {
		t.Errorf("initial status should be pending, got %s", b.Status())
	}
}

func TestBaseExecutor_Transition(t *testing.T) {
	b := NewBaseExecutor()

	// pending -> running
	err := b.Transition("start")
	if err != nil {
		t.Fatalf("transition to running should succeed: %v", err)
	}
	if b.Status() != StatusRunning {
		t.Errorf("status should be running, got %s", b.Status())
	}

	// running -> success
	err = b.Transition("success")
	if err != nil {
		t.Fatalf("transition to success should succeed: %v", err)
	}
	if b.Status() != StatusSuccess {
		t.Errorf("status should be success, got %s", b.Status())
	}
}

func TestBaseExecutor_TransitionFail(t *testing.T) {
	b := NewBaseExecutor()

	// pending -> running
	b.Transition("start")

	// running -> failed
	err := b.Transition("fail")
	if err != nil {
		t.Fatalf("transition to failed should succeed: %v", err)
	}
	if b.Status() != StatusFailed {
		t.Errorf("status should be failed, got %s", b.Status())
	}
}

func TestBaseExecutor_TransitionCancel(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*BaseExecutor)
		expect Status
	}{
		{
			name:   "cancel from pending",
			setup:  func(b *BaseExecutor) {},
			expect: StatusCancelled,
		},
		{
			name: "cancel from running",
			setup: func(b *BaseExecutor) {
				b.Transition("start")
			},
			expect: StatusCancelled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBaseExecutor()
			tt.setup(b)
			err := b.Transition("cancel")
			if err != nil {
				t.Fatalf("transition should succeed: %v", err)
			}
			if b.Status() != tt.expect {
				t.Errorf("status should be %s, got %s", tt.expect, b.Status())
			}
		})
	}
}

func TestBaseExecutor_UpdateProgress(t *testing.T) {
	b := NewBaseExecutor()

	b.UpdateProgress("init", 10, "initializing")

	p := b.GetProgress()
	if p.Phase != "init" {
		t.Errorf("phase should be init, got %s", p.Phase)
	}
	if p.Percent != 10 {
		t.Errorf("percent should be 10, got %d", p.Percent)
	}
	if p.Message != "initializing" {
		t.Errorf("message should be initializing, got %s", p.Message)
	}
}

func TestBaseExecutor_GetProgressCopy(t *testing.T) {
	b := NewBaseExecutor()
	b.UpdateProgress("test", 50, "testing")

	p1 := b.GetProgress()
	p2 := b.GetProgress()

	// 修改返回的 progress 不应影响原始数据
	p1.Phase = "modified"

	if p2.Phase == "modified" {
		t.Error("GetProgress should return a copy")
	}
}

func TestBaseExecutor_Cancel(t *testing.T) {
	b := NewBaseExecutor()

	// 设置 cancel 函数
	ctx, cancel := context.WithCancel(context.Background())
	b.SetCancel(cancel)

	// 调用 Cancel
	err := b.Cancel()
	if err != nil {
		t.Fatalf("cancel should succeed: %v", err)
	}

	// context 应该被取消
	select {
	case <-ctx.Done():
		// 预期行为
	default:
		t.Error("context should be cancelled")
	}

	if b.Status() != StatusCancelled {
		t.Errorf("status should be cancelled, got %s", b.Status())
	}
}

func TestBaseExecutor_CancelWithoutCancelFunc(t *testing.T) {
	b := NewBaseExecutor()

	// 没有设置 cancel 函数也不应 panic
	err := b.Cancel()
	if err != nil {
		t.Fatalf("cancel without cancel func should succeed: %v", err)
	}
}
