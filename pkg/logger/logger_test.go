package logger

import (
	"testing"
)

func TestZapLogger(t *testing.T) {
	config := Config{
		Type:         LoggerTypeZap,
		Level:        LevelDebug,
		Output:       OutputStdout,
		EnableCaller: true,
		JSONFormat:   false,
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize zap logger: %v", err)
	}

	// 测试基本日志
	Debug("Debug message", String("key", "value"))
	Info("Info message", Int("count", 42))
	Warn("Warning message", Bool("flag", true))
	Error("Error message", String("error", "something went wrong"))

	// 测试 With
	logger := With(String("component", "test"))
	logger.Info("Message with component field")

	// 测试 Named
	namedLogger := Named("test-logger")
	namedLogger.Info("Message from named logger")

	// 刷新缓冲区
	if err := Sync(); err != nil {
		t.Logf("Sync error (expected on stdout): %v", err)
	}
}

func TestZerologLogger(t *testing.T) {
	config := Config{
		Type:         LoggerTypeZerolog,
		Level:        LevelDebug,
		Output:       OutputStdout,
		EnableCaller: true,
		JSONFormat:   false,
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize zerolog logger: %v", err)
	}

	// 测试基本日志
	Debug("Debug message", String("key", "value"))
	Info("Info message", Int("count", 42))
	Warn("Warning message", Bool("flag", true))
	Error("Error message", String("error", "something went wrong"))

	// 测试 With
	logger := With(String("component", "test"))
	logger.Info("Message with component field")

	// 测试 Named
	namedLogger := Named("test-logger")
	namedLogger.Info("Message from named logger")
}

func TestFileOutput(t *testing.T) {
	config := Config{
		Type:   LoggerTypeZap,
		Level:  LevelInfo,
		Output: OutputFile,
		File: FileConfig{
			Filename:   "/tmp/test-prism.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
		},
		JSONFormat: true,
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize logger with file output: %v", err)
	}

	Info("Test message to file", String("test", "value"))

	if err := Sync(); err != nil {
		t.Logf("Sync error: %v", err)
	}
}

func TestBothOutput(t *testing.T) {
	config := Config{
		Type:   LoggerTypeZerolog,
		Level:  LevelInfo,
		Output: OutputBoth,
		File: FileConfig{
			Filename:   "/tmp/test-prism-both.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   false,
		},
		JSONFormat: false,
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize logger with both output: %v", err)
	}

	Info("Test message to both stdout and file", String("test", "value"))
}
