package cmd

import (
	"context"
	"testing"
	"time"
)

func TestRunner_New(t *testing.T) {
	r := NewRunner(0)
	if r == nil {
		t.Fatal("NewRunner should not return nil")
	}
	if r.timeout != 30*time.Minute {
		t.Errorf("default timeout should be 30 minutes, got %v", r.timeout)
	}

	r2 := NewRunner(5 * time.Second)
	if r2.timeout != 5*time.Second {
		t.Errorf("timeout should be 5 seconds, got %v", r2.timeout)
	}
}

func TestRunner_Exec_EmptyCommand(t *testing.T) {
	r := NewRunner(5 * time.Second)
	result := r.Exec(context.Background(), []string{})

	if result.Error == nil {
		t.Error("empty command should return error")
	}
}

func TestRunner_Exec_SimpleCommand(t *testing.T) {
	r := NewRunner(5 * time.Second)
	result := r.Exec(context.Background(), []string{"echo", "hello"})

	if result.Error != nil {
		t.Fatalf("echo should succeed: %v", result.Error)
	}
	if result.ExitCode != 0 {
		t.Errorf("exit code should be 0, got %d", result.ExitCode)
	}
	if result.Output != "hello\n" {
		t.Errorf("output should be 'hello\\n', got %q", result.Output)
	}
}

func TestRunner_Exec_FailingCommand(t *testing.T) {
	r := NewRunner(5 * time.Second)
	result := r.Exec(context.Background(), []string{"false"})

	if result.Error == nil {
		t.Error("false command should return error")
	}
	if result.ExitCode == 0 {
		t.Error("exit code should not be 0")
	}
}

func TestRunner_Exec_NonExistentCommand(t *testing.T) {
	r := NewRunner(5 * time.Second)
	result := r.Exec(context.Background(), []string{"nonexistent_command_12345"})

	if result.Error == nil {
		t.Error("nonexistent command should return error")
	}
}

func TestRunner_ExecWithHandler_EmptyCommand(t *testing.T) {
	r := NewRunner(5 * time.Second)
	result := r.ExecWithHandler(context.Background(), []string{}, nil)

	if result.Error == nil {
		t.Error("empty command should return error")
	}
}

func TestRunner_ExecWithHandler(t *testing.T) {
	r := NewRunner(5 * time.Second)

	var lines []string
	handler := func(line string) {
		lines = append(lines, line)
	}

	result := r.ExecWithHandler(context.Background(), []string{"echo", "-e", "line1\nline2"}, handler)

	if result.Error != nil {
		t.Fatalf("echo should succeed: %v", result.Error)
	}

	if len(lines) < 1 {
		t.Error("handler should receive at least one line")
	}
}

func TestRunner_Exec_ContextCancelled(t *testing.T) {
	r := NewRunner(30 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	result := r.Exec(ctx, []string{"sleep", "10"})

	if result.Error == nil {
		t.Error("cancelled context should return error")
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "\x1b[32mgreen\x1b[0m",
			expected: "green",
		},
		{
			input:    "\x1b[1;31mred bold\x1b[0m",
			expected: "red bold",
		},
		{
			input:    "no ansi",
			expected: "no ansi",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		result := StripANSI(tt.input)
		if result != tt.expected {
			t.Errorf("StripANSI(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
