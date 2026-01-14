// Package cmd provides command execution utilities.
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"syscall"
	"time"
)

// Result 命令执行结果
type Result struct {
	Output   string
	ExitCode int
	Duration time.Duration
	Error    error
}

// Runner 命令执行器
type Runner struct {
	timeout time.Duration
}

// NewRunner 创建命令执行器
func NewRunner(timeout time.Duration) *Runner {
	if timeout == 0 {
		timeout = 30 * time.Minute
	}
	return &Runner{timeout: timeout}
}

// Exec 执行命令
func (r *Runner) Exec(ctx context.Context, args []string) *Result {
	if len(args) == 0 {
		return &Result{Error: fmt.Errorf("empty command")}
	}

	start := time.Now()
	result := &Result{}

	// 创建带超时的 context
	timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	cmd := exec.CommandContext(timeoutCtx, args[0], args[1:]...)
	// 设置进程组，便于杀死子进程
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()
	result.Duration = time.Since(start)

	// 合并输出
	result.Output = stdout.String() + stderr.String()

	if err != nil {
		// 检查是否超时
		if timeoutCtx.Err() == context.DeadlineExceeded {
			// 杀死进程组
			if cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			}
			result.Error = fmt.Errorf("command timed out after %v", r.timeout)
			result.ExitCode = -1
			return result
		}

		// 获取退出码
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		result.Error = err
	}

	return result
}

// ExecWithHandler 执行命令并实时处理输出
func (r *Runner) ExecWithHandler(ctx context.Context, args []string, handler func(line string)) *Result {
	if len(args) == 0 {
		return &Result{Error: fmt.Errorf("empty command")}
	}

	start := time.Now()
	result := &Result{}

	timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	cmd := exec.CommandContext(timeoutCtx, args[0], args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 使用自定义 writer 处理输出
	output := &outputWriter{
		handler: handler,
		buffer:  &bytes.Buffer{},
	}
	cmd.Stdout = output
	cmd.Stderr = output

	err := cmd.Run()
	result.Duration = time.Since(start)
	result.Output = output.buffer.String()

	if err != nil {
		if timeoutCtx.Err() == context.DeadlineExceeded {
			if cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			}
			result.Error = fmt.Errorf("command timed out after %v", r.timeout)
			result.ExitCode = -1
			return result
		}

		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		result.Error = err
	}

	return result
}

// outputWriter 输出处理器
type outputWriter struct {
	handler func(line string)
	buffer  *bytes.Buffer
	line    bytes.Buffer
}

func (w *outputWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	w.buffer.Write(p)

	for _, b := range p {
		if b == '\n' {
			if w.handler != nil {
				w.handler(w.line.String())
			}
			w.line.Reset()
		} else {
			w.line.WriteByte(b)
		}
	}
	return n, nil
}

// StripANSI 移除 ANSI 转义序列
func StripANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}
