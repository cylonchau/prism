package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger GORM 日志适配器
type GormLogger struct {
	logger                    Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// NewGormLogger 创建 GORM 日志适配器
func NewGormLogger(logger Logger) *GormLogger {
	return &GormLogger{
		logger:                    logger,
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 记录 info 级别日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.logger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 记录 warn 级别日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 记录 error 级别日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 记录 SQL 执行日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.logger.Error("SQL execution error",
			String("sql", sql),
			Int64("rows", rows),
			String("elapsed", elapsed.String()),
			Err(err),
		)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		l.logger.Warn("Slow SQL query",
			String("sql", sql),
			Int64("rows", rows),
			String("elapsed", elapsed.String()),
			String("threshold", l.SlowThreshold.String()),
		)
	case l.LogLevel >= gormlogger.Info:
		l.logger.Debug("SQL query",
			String("sql", sql),
			Int64("rows", rows),
			String("elapsed", elapsed.String()),
		)
	}
}
