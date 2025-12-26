// Package executor defines the executor interface and common types.
package executor

import (
	"context"
)

// Action 执行动作
type Action string

const (
	ActionInit    Action = "init"
	ActionPlan    Action = "plan"
	ActionApply   Action = "apply"
	ActionDestroy Action = "destroy"
	ActionImport  Action = "import"
)

// Status 执行状态
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusSuccess   Status = "success"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// ExecuteRequest 执行请求
type ExecuteRequest struct {
	TaskID     string            // 任务ID
	ResourceID int64             // 资源ID
	Action     Action            // 执行动作
	WorkDir    string            // 工作目录
	Config     string            // 配置内容
	Params     map[string]string // 额外参数
}

// ExecuteResult 执行结果
type ExecuteResult struct {
	TaskID     string
	Status     Status
	Output     string            // 完整输出
	Error      string            // 错误信息
	Duration   int64             // 执行时长(ms)
	Attributes map[string]string // 提取的属性
}

// Progress 进度信息
type Progress struct {
	Phase   string // 当前阶段
	Percent int    // 百分比 0-100
	Elapsed int64  // 已用时间(ms)
	Message string // 当前消息
}

// Executor 执行器接口
type Executor interface {
	// Type 返回执行器类型
	Type() string

	// Execute 执行任务
	Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResult, error)

	// Validate 验证配置
	Validate(config string) error

	// GetProgress 获取进度
	GetProgress() *Progress

	// Cancel 取消执行
	Cancel() error
}
