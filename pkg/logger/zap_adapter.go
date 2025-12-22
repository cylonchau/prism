package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zapLogger zap 日志适配器
type zapLogger struct {
	logger *zap.Logger
}

// newZapLogger 创建 zap 日志实例
func newZapLogger(config Config) (Logger, error) {
	// 配置日志级别
	level := parseZapLevel(config.Level)

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if config.JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置输出
	var writers []zapcore.WriteSyncer

	switch config.Output {
	case OutputStdout:
		writers = append(writers, zapcore.AddSync(os.Stdout))
	case OutputFile:
		writers = append(writers, zapcore.AddSync(newLumberjackWriter(config.File)))
	case OutputBoth:
		writers = append(writers, zapcore.AddSync(os.Stdout))
		writers = append(writers, zapcore.AddSync(newLumberjackWriter(config.File)))
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writers...),
		level,
	)

	// 配置选项
	opts := []zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	if config.EnableCaller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	logger := zap.New(core, opts...)

	return &zapLogger{
		logger: logger,
	}, nil
}

// parseZapLevel 解析日志级别
func parseZapLevel(level Level) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// newLumberjackWriter 创建文件轮转 writer
func newLumberjackWriter(config FileConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
}

// fieldsToZap 转换字段为 zap 字段
func fieldsToZap(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}

// Debug 记录 debug 级别日志
func (z *zapLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, fieldsToZap(fields)...)
}

// Info 记录 info 级别日志
func (z *zapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, fieldsToZap(fields)...)
}

// Warn 记录 warn 级别日志
func (z *zapLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, fieldsToZap(fields)...)
}

// Error 记录 error 级别日志
func (z *zapLogger) Error(msg string, fields ...Field) {
	z.logger.Error(msg, fieldsToZap(fields)...)
}

// Fatal 记录 fatal 级别日志并退出程序
func (z *zapLogger) Fatal(msg string, fields ...Field) {
	z.logger.Fatal(msg, fieldsToZap(fields)...)
}

// With 添加字段
func (z *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: z.logger.With(fieldsToZap(fields)...),
	}
}

// Named 添加命名空间
func (z *zapLogger) Named(name string) Logger {
	return &zapLogger{
		logger: z.logger.Named(name),
	}
}

// Sync 刷新缓冲区
func (z *zapLogger) Sync() error {
	return z.logger.Sync()
}
