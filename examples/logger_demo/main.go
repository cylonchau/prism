package main

import (
	"context"
	"time"

	"github.com/cylonchau/prism/pkg/logger"
	"github.com/cylonchau/prism/pkg/store"
)

func main() {
	// 示例 1: 使用 zap + stdout
	zapConfig := logger.Config{
		Type:         logger.LoggerTypeZap,
		Level:        logger.LevelInfo,
		Output:       logger.OutputStdout,
		EnableCaller: true,
		JSONFormat:   false,
	}

	if err := logger.Initialize(zapConfig); err != nil {
		panic(err)
	}

	logger.Info("Application started with zap logger")
	logger.Info("Database configuration",
		logger.String("type", "mysql"),
		logger.Int("port", 3306),
	)

	// 示例 2: 切换到 zerolog + 文件输出
	zerologConfig := logger.Config{
		Type:   logger.LoggerTypeZerolog,
		Level:  logger.LevelDebug,
		Output: logger.OutputBoth, // 同时输出到 stdout 和文件
		File: logger.FileConfig{
			Filename:   "/tmp/prism.log",
			MaxSize:    100, // MB
			MaxBackups: 5,
			MaxAge:     30, // days
			Compress:   true,
		},
		EnableCaller: true,
		JSONFormat:   false,
	}

	if err := logger.Initialize(zerologConfig); err != nil {
		panic(err)
	}

	logger.Info("Switched to zerolog logger")
	logger.Debug("This is a debug message")

	// 示例 3: 使用命名空间
	dbLogger := logger.Named("database")
	dbLogger.Info("Database logger initialized")

	// 示例 4: 使用字段
	requestLogger := logger.With(
		logger.String("request_id", "req-12345"),
		logger.String("user", "admin"),
	)
	requestLogger.Info("Processing request")
	requestLogger.Warn("Request took too long", logger.Int("duration_ms", 1500))

	// 示例 5: 集成到数据库
	dbConfig := store.DatabaseConfig{
		Type:              store.SQLite,
		File:              "/tmp/test",
		MaxOpenConnection: "10",
		MaxIdleConnection: "5",
	}

	db := store.GetInstance()
	if err := db.Initialize(dbConfig); err != nil {
		logger.Error("Failed to initialize database", logger.Err(err))
		return
	}

	logger.Info("Database initialized successfully")

	// 启动连接池监控
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go db.MonitorConnectionPool(ctx)

	// 等待一段时间以查看监控日志
	time.Sleep(2 * time.Second)

	// 清理
	if err := db.Close(); err != nil {
		logger.Error("Failed to close database", logger.Err(err))
	}

	logger.Info("Application shutdown")
	logger.Sync()
}
