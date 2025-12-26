// Package terraform provides Terraform executor implementation.
package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/cylonchau/prism/pkg/dao"
	"github.com/cylonchau/prism/pkg/executor"
	"github.com/cylonchau/prism/pkg/executor/cmd"
	"github.com/cylonchau/prism/pkg/executor/lock"
	"github.com/cylonchau/prism/pkg/executor/workspace"
	"github.com/cylonchau/prism/pkg/executor/ws"
)

// Config holds Terraform executor configuration.
type Config struct {
	BinaryPath string
	BasePath   string
	Timeout    time.Duration
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		BinaryPath: "terraform",
		BasePath:   "/opt/homebrew/bin/terraform",
		Timeout:    30 * time.Minute,
	}
}

// Executor implements Terraform execution.
type Executor struct {
	*executor.BaseExecutor
	config    *Config
	locker    lock.LockManager
	taskDAO   *dao.ExecutionTaskDAO
	workspace *workspace.Manager
	runner    *cmd.Runner
	hub       *ws.Hub
	parser    *Parser
	errors    []Diagnostic // Extracted errors from JSON output
}

// New creates a new Terraform executor.
func New(config *Config, locker lock.LockManager, taskDAO *dao.ExecutionTaskDAO, hub *ws.Hub) *Executor {
	if config == nil {
		config = DefaultConfig()
	}
	return &Executor{
		BaseExecutor: executor.NewBaseExecutor(),
		config:       config,
		locker:       locker,
		taskDAO:      taskDAO,
		workspace:    workspace.NewManager(config.BasePath),
		runner:       cmd.NewRunner(config.Timeout),
		hub:          hub,
		parser:       NewParser(),
		errors:       []Diagnostic{},
	}
}

// Type returns the executor type.
func (e *Executor) Type() string {
	return "terraform"
}

func (e *Executor) Execute(ctx context.Context, req *executor.ExecuteRequest) (*executor.ExecuteResult, error) {
	start := time.Now()
	e.errors = []Diagnostic{} // Reset errors

	result := &executor.ExecuteResult{
		TaskID: req.TaskID,
		Status: executor.StatusRunning,
	}

	// 1. Create task record
	if e.taskDAO != nil {
		_, err := e.taskDAO.Create(req.TaskID, req.ResourceID, string(req.Action))
		if err != nil {
			result.Status = executor.StatusFailed
			result.Error = err.Error()
			return result, err
		}
	}

	// 2. Acquire lock
	if e.locker != nil {
		if err := e.locker.Acquire(ctx, req.ResourceID, req.TaskID); err != nil {
			result.Status = executor.StatusFailed
			result.Error = err.Error()
			e.completeTask(req.TaskID, false, err.Error())
			return result, err
		}
		defer e.locker.Release(req.ResourceID)
	}

	// 3. Start task
	if e.taskDAO != nil {
		e.taskDAO.Start(req.TaskID)
	}
	if err := e.Transition("start"); err != nil {
		result.Status = executor.StatusFailed
		result.Error = err.Error()
		e.completeTask(req.TaskID, false, err.Error())
		return result, err
	}

	// 4. Setup cancellation
	ctx, cancel := context.WithCancel(ctx)
	e.SetCancel(cancel)

	// 5. Create work directory
	workDir := req.WorkDir
	if workDir == "" {
		var err error
		workDir, err = e.workspace.Create("default", "default", fmt.Sprintf("%d", req.ResourceID), req.TaskID)
		if err != nil {
			result.Status = executor.StatusFailed
			result.Error = err.Error()
			e.Transition("fail")
			e.completeTask(req.TaskID, false, err.Error())
			return result, err
		}
		defer e.workspace.Clean(workDir)
	}

	// 6. Execute action
	var err error
	switch req.Action {
	case executor.ActionInit:
		err = e.init(ctx, workDir, req)
	case executor.ActionPlan:
		err = e.plan(ctx, workDir, req)
	case executor.ActionApply:
		err = e.apply(ctx, workDir, req)
	case executor.ActionDestroy:
		err = e.destroy(ctx, workDir, req)
	default:
		err = fmt.Errorf("unsupported action: %s", req.Action)
	}

	// 7. Update result
	result.Duration = time.Since(start).Milliseconds()
	result.Output = e.getErrorSummary()

	if err != nil {
		result.Status = executor.StatusFailed
		result.Error = err.Error()
		e.Transition("fail")
		e.completeTask(req.TaskID, false, err.Error())
		e.sendComplete(req.TaskID, false, result)
		return result, err
	}

	result.Status = executor.StatusSuccess
	e.Transition("success")
	e.completeTask(req.TaskID, true, "")
	e.sendComplete(req.TaskID, true, result)
	return result, nil
}

// completeTask persists task completion.
func (e *Executor) completeTask(taskID string, success bool, errMsg string) {
	if e.taskDAO != nil {
		e.taskDAO.Complete(taskID, success, e.getErrorSummary(), errMsg)
	}
}

// Retry resets and re-executes a failed task.
func (e *Executor) Retry(ctx context.Context, taskID string) (*executor.ExecuteResult, error) {
	if e.taskDAO == nil {
		return nil, fmt.Errorf("task store not configured")
	}

	canRetry, err := e.taskDAO.CanRetry(taskID)
	if err != nil {
		return nil, err
	}
	if !canRetry {
		return nil, fmt.Errorf("task %s is not retryable", taskID)
	}

	// Get original task
	task, err := e.taskDAO.Get(taskID)
	if err != nil {
		return nil, err
	}

	// Reset task
	if err := e.taskDAO.Reset(taskID); err != nil {
		return nil, err
	}

	// Re-execute
	req := &executor.ExecuteRequest{
		TaskID:     taskID,
		ResourceID: task.ResourceID,
		Action:     executor.Action(task.Action),
	}
	return e.Execute(ctx, req)
}

// Validate 验证配置
func (e *Executor) Validate(config string) error {
	// TODO: 使用 terraform validate
	return nil
}

// init 执行 terraform init
func (e *Executor) init(ctx context.Context, workDir string, req *executor.ExecuteRequest) error {
	e.UpdateProgress("init", 10, "Running terraform init...")
	e.sendProgress(req.TaskID, "init", 10, "Running terraform init...")

	args := []string{
		e.config.BinaryPath,
		"-chdir=" + workDir,
		"init",
		"-no-color",
	}

	result := e.runner.ExecWithHandler(ctx, args, func(line string) {
		cleaned := cmd.StripANSI(line)
		e.sendLog(req.TaskID, cleaned)
	})

	if result.Error != nil {
		return fmt.Errorf("terraform init failed: %w", result.Error)
	}
	return nil
}

// plan 执行 terraform plan
func (e *Executor) plan(ctx context.Context, workDir string, req *executor.ExecuteRequest) error {
	// 先 init
	if err := e.init(ctx, workDir, req); err != nil {
		return err
	}

	e.UpdateProgress("plan", 30, "Running terraform plan...")
	e.sendProgress(req.TaskID, "plan", 30, "Running terraform plan...")

	args := []string{
		e.config.BinaryPath,
		"-chdir=" + workDir,
		"plan",
		"-input=false",
		"-json",
	}

	result := e.runner.ExecWithHandler(ctx, args, func(line string) {
		cleaned := cmd.StripANSI(line)
		e.sendLog(req.TaskID, cleaned)
	})

	if result.Error != nil {
		return fmt.Errorf("terraform plan failed: %w", result.Error)
	}

	// 解析 plan 输出
	planInfo := e.parser.ParsePlan(result.Output)
	e.sendLog(req.TaskID, fmt.Sprintf("Plan: %d to add, %d to change, %d to destroy",
		planInfo.ToAdd, planInfo.ToChange, planInfo.ToDestroy))

	return nil
}

// apply 执行 terraform apply
func (e *Executor) apply(ctx context.Context, workDir string, req *executor.ExecuteRequest) error {
	// 先 init
	if err := e.init(ctx, workDir, req); err != nil {
		return err
	}

	e.UpdateProgress("apply", 50, "Running terraform apply...")
	e.sendProgress(req.TaskID, "apply", 50, "Running terraform apply...")

	args := []string{
		e.config.BinaryPath,
		"-chdir=" + workDir,
		"apply",
		"-auto-approve",
		"-json",
	}

	result := e.runner.ExecWithHandler(ctx, args, func(line string) {
		cleaned := cmd.StripANSI(line)
		e.sendLog(req.TaskID, cleaned)
	})

	if result.Error != nil {
		return fmt.Errorf("terraform apply failed: %w", result.Error)
	}

	e.UpdateProgress("apply", 90, "Parsing tfstate...")

	// 解析 tfstate
	tfstatePath := filepath.Join(workDir, "terraform.tfstate")
	if e.workspace.Exists(tfstatePath) {
		data, err := e.workspace.ReadFile(workDir, "terraform.tfstate")
		if err == nil {
			attrs := e.parser.ParseTfstate(data)
			e.sendLog(req.TaskID, fmt.Sprintf("Extracted %d attributes from tfstate", len(attrs)))
		}
	}

	return nil
}

// destroy 执行 terraform destroy
func (e *Executor) destroy(ctx context.Context, workDir string, req *executor.ExecuteRequest) error {
	// 先 init
	if err := e.init(ctx, workDir, req); err != nil {
		return err
	}

	e.UpdateProgress("destroy", 50, "Running terraform destroy...")
	e.sendProgress(req.TaskID, "destroy", 50, "Running terraform destroy...")

	args := []string{
		e.config.BinaryPath,
		"-chdir=" + workDir,
		"destroy",
		"-auto-approve",
		"-json",
	}

	result := e.runner.ExecWithHandler(ctx, args, func(line string) {
		cleaned := cmd.StripANSI(line)
		e.sendLog(req.TaskID, cleaned)
	})

	if result.Error != nil {
		return fmt.Errorf("terraform destroy failed: %w", result.Error)
	}

	return nil
}

// sendLog processes a JSON line and extracts errors.
func (e *Executor) sendLog(taskID, line string) {
	// Parse JSON message
	var msg TerraformMessage
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		// Not JSON, just send as-is
		if e.hub != nil {
			e.hub.SendLog(taskID, line)
		}
		return
	}

	// Send human-readable message to websocket
	if e.hub != nil {
		e.hub.SendLog(taskID, msg.Message)
	}

	// Extract and store errors
	if msg.Type == "diagnostic" && msg.Diagnostic != nil && msg.Diagnostic.Severity == "error" {
		e.errors = append(e.errors, *msg.Diagnostic)
	}
}

// getErrorSummary returns formatted error summary for storage.
func (e *Executor) getErrorSummary() string {
	if len(e.errors) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, err := range e.errors {
		if i > 0 {
			sb.WriteString("\n---\n")
		}
		sb.WriteString(fmt.Sprintf("Error: %s\n", err.Summary))
		if err.Address != "" {
			sb.WriteString(fmt.Sprintf("Resource: %s\n", err.Address))
		}
		if err.Range != nil {
			sb.WriteString(fmt.Sprintf("File: %s:%d\n", err.Range.Filename, err.Range.Start.Line))
		}
		if err.Detail != "" {
			sb.WriteString(fmt.Sprintf("Detail: %s\n", err.Detail))
		}
	}
	return sb.String()
}

// GetErrors returns extracted errors.
func (e *Executor) GetErrors() []Diagnostic {
	return e.errors
}

// sendProgress sends progress update.
func (e *Executor) sendProgress(taskID, phase string, percent int, message string) {
	if e.hub != nil {
		e.hub.SendProgress(taskID, &ws.ProgressData{
			Phase:   phase,
			Percent: percent,
			Elapsed: 0,
			Message: message,
		})
	}
}

// sendComplete sends completion message.
func (e *Executor) sendComplete(taskID string, success bool, result interface{}) {
	if e.hub != nil {
		e.hub.SendComplete(taskID, success, result)
	}
}
