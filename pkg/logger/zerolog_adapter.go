package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// zerologLogger zerolog 日志适配器
type zerologLogger struct {
	logger zerolog.Logger
}

// newZerologLogger 创建 zerolog 日志实例
func newZerologLogger(config Config) (Logger, error) {
	// 配置日志级别
	level := parseZerologLevel(config.Level)
	zerolog.SetGlobalLevel(level)

	// 配置输出
	var writers []io.Writer

	switch config.Output {
	case OutputStdout:
		if config.JSONFormat {
			writers = append(writers, os.Stdout)
		} else {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			})
		}
	case OutputFile:
		writers = append(writers, newLumberjackWriter(config.File))
	case OutputBoth:
		if config.JSONFormat {
			writers = append(writers, os.Stdout)
		} else {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			})
		}
		writers = append(writers, newLumberjackWriter(config.File))
	}

	multi := zerolog.MultiLevelWriter(writers...)

	// 创建 logger
	logger := zerolog.New(multi).With().Timestamp()

	if config.EnableCaller {
		logger = logger.Caller()
	}

	return &zerologLogger{
		logger: logger.Logger(),
	}, nil
}

// parseZerologLevel 解析日志级别
func parseZerologLevel(level Level) zerolog.Level {
	switch level {
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelError:
		return zerolog.ErrorLevel
	case LevelFatal:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// fieldsToZerolog 转换字段为 zerolog 事件
func (z *zerologLogger) fieldsToZerolog(event *zerolog.Event, fields []Field) *zerolog.Event {
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}
	return event
}

// Debug 记录 debug 级别日志
func (z *zerologLogger) Debug(msg string, fields ...Field) {
	event := z.logger.Debug()
	z.fieldsToZerolog(event, fields).Msg(msg)
}

// Info 记录 info 级别日志
func (z *zerologLogger) Info(msg string, fields ...Field) {
	event := z.logger.Info()
	z.fieldsToZerolog(event, fields).Msg(msg)
}

// Warn 记录 warn 级别日志
func (z *zerologLogger) Warn(msg string, fields ...Field) {
	event := z.logger.Warn()
	z.fieldsToZerolog(event, fields).Msg(msg)
}

// Error 记录 error 级别日志
func (z *zerologLogger) Error(msg string, fields ...Field) {
	event := z.logger.Error()
	z.fieldsToZerolog(event, fields).Msg(msg)
}

// Fatal 记录 fatal 级别日志并退出程序
func (z *zerologLogger) Fatal(msg string, fields ...Field) {
	event := z.logger.Fatal()
	z.fieldsToZerolog(event, fields).Msg(msg)
}

// With 添加字段
func (z *zerologLogger) With(fields ...Field) Logger {
	ctx := z.logger.With()
	for _, f := range fields {
		ctx = ctx.Interface(f.Key, f.Value)
	}
	return &zerologLogger{
		logger: ctx.Logger(),
	}
}

// Named 添加命名空间
func (z *zerologLogger) Named(name string) Logger {
	return &zerologLogger{
		logger: z.logger.With().Str("logger", name).Logger(),
	}
}

// Sync 刷新缓冲区
func (z *zerologLogger) Sync() error {
	// zerolog 不需要显式刷新
	return nil
}
