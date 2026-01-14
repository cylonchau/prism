package executor

import (
	"context"
	"sync"

	"github.com/looplab/fsm"
)

// BaseExecutor 基础执行器，提供公共功能
type BaseExecutor struct {
	mu       sync.RWMutex
	taskID   string
	progress *Progress
	fsm      *fsm.FSM
	cancel   context.CancelFunc
}

// NewBaseExecutor 创建基础执行器
func NewBaseExecutor() *BaseExecutor {
	b := &BaseExecutor{
		progress: &Progress{},
	}
	b.initFSM()
	return b
}

// initFSM 初始化状态机
func (b *BaseExecutor) initFSM() {
	b.fsm = fsm.NewFSM(
		string(StatusPending),
		fsm.Events{
			{Name: "start", Src: []string{string(StatusPending)}, Dst: string(StatusRunning)},
			{Name: "success", Src: []string{string(StatusRunning)}, Dst: string(StatusSuccess)},
			{Name: "fail", Src: []string{string(StatusRunning)}, Dst: string(StatusFailed)},
			{Name: "cancel", Src: []string{string(StatusPending), string(StatusRunning)}, Dst: string(StatusCancelled)},
		},
		fsm.Callbacks{},
	)
}

// GetProgress 获取进度
func (b *BaseExecutor) GetProgress() *Progress {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return &Progress{
		Phase:   b.progress.Phase,
		Percent: b.progress.Percent,
		Elapsed: b.progress.Elapsed,
		Message: b.progress.Message,
	}
}

// UpdateProgress 更新进度
func (b *BaseExecutor) UpdateProgress(phase string, percent int, message string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.progress.Phase = phase
	b.progress.Percent = percent
	b.progress.Message = message
}

// Status 获取当前状态
func (b *BaseExecutor) Status() Status {
	return Status(b.fsm.Current())
}

// Transition 状态转换
func (b *BaseExecutor) Transition(event string) error {
	return b.fsm.Event(context.Background(), event)
}

// Cancel 取消执行
func (b *BaseExecutor) Cancel() error {
	if b.cancel != nil {
		b.cancel()
	}
	return b.Transition("cancel")
}

// SetCancel 设置取消函数
func (b *BaseExecutor) SetCancel(cancel context.CancelFunc) {
	b.cancel = cancel
}
