package logger

import (
	"fmt"
	"sync"
)

// LoggerType 日志驱动类型
type LoggerType string

const (
	LoggerTypeZap     LoggerType = "zap"
	LoggerTypeZerolog LoggerType = "zerolog"
)

// OutputType 输出类型
type OutputType string

const (
	OutputStdout OutputType = "stdout"
	OutputFile   OutputType = "file"
	OutputBoth   OutputType = "both"
)

// Level 日志级别
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
)

// FileConfig 文件输出配置
type FileConfig struct {
	Filename   string // 日志文件路径
	MaxSize    int    // 单个文件最大大小(MB)
	MaxBackups int    // 保留的旧文件最大数量
	MaxAge     int    // 保留旧文件的最大天数
	Compress   bool   // 是否压缩旧文件
}

// Config 日志配置
type Config struct {
	Type         LoggerType // 日志驱动类型: zap 或 zerolog
	Level        Level      // 日志级别
	Output       OutputType // 输出类型: stdout, file, both
	File         FileConfig // 文件输出配置
	EnableCaller bool       // 是否显示调用者信息
	JSONFormat   bool       // 是否使用 JSON 格式
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// Logger 统一日志接口
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	With(fields ...Field) Logger // 添加字段
	Named(name string) Logger    // 添加命名空间
	Sync() error                 // 刷新缓冲区
}

// 字段构造函数
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Err(err error) Field {
	return Field{Key: "error", Value: err}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

var (
	defaultLogger Logger
	once          sync.Once
	mu            sync.RWMutex
)

// Initialize 初始化全局日志实例
func Initialize(config Config) error {
	mu.Lock()
	defer mu.Unlock()

	var logger Logger
	var err error

	switch config.Type {
	case LoggerTypeZap:
		logger, err = newZapLogger(config)
	case LoggerTypeZerolog:
		logger, err = newZerologLogger(config)
	default:
		return fmt.Errorf("unsupported logger type: %s", config.Type)
	}

	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	defaultLogger = logger
	return nil
}

// Default 获取默认日志实例
func Default() Logger {
	mu.RLock()
	defer mu.RUnlock()

	if defaultLogger == nil {
		// 如果未初始化，使用默认配置
		once.Do(func() {
			mu.RUnlock()
			_ = Initialize(Config{
				Type:   LoggerTypeZap,
				Level:  LevelInfo,
				Output: OutputStdout,
			})
			mu.RLock()
		})
	}

	return defaultLogger
}

// Debug 记录 debug 级别日志
func Debug(msg string, fields ...Field) {
	Default().Debug(msg, fields...)
}

// Info 记录 info 级别日志
func Info(msg string, fields ...Field) {
	Default().Info(msg, fields...)
}

// Warn 记录 warn 级别日志
func Warn(msg string, fields ...Field) {
	Default().Warn(msg, fields...)
}

// Error 记录 error 级别日志
func Error(msg string, fields ...Field) {
	Default().Error(msg, fields...)
}

// Fatal 记录 fatal 级别日志并退出程序
func Fatal(msg string, fields ...Field) {
	Default().Fatal(msg, fields...)
}

// With 添加字段
func With(fields ...Field) Logger {
	return Default().With(fields...)
}

// Named 添加命名空间
func Named(name string) Logger {
	return Default().Named(name)
}

// Sync 刷新缓冲区
func Sync() error {
	return Default().Sync()
}
