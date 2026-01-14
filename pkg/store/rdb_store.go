// pkg/store/store.go
package store

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cylonchau/prism/pkg/logger"
)

var (
	instance Store
	initOnce sync.Once
)

// RDBStore implements the Store interface for relational databases.
type RDBStore struct {
	db     *gorm.DB
	config DatabaseConfig
	once   sync.Once
	mu     sync.RWMutex
}

// GetInstance returns the singleton store instance.
func GetInstance() Store {
	initOnce.Do(func() {
		instance = &RDBStore{}
	})
	return instance
}

// NewRDBStore creates a new RDBStore instance (useful for testing).
func NewRDBStore() *RDBStore {
	return &RDBStore{}
}

// Initialize 初始化数据库连接
func (m *RDBStore) Initialize(config DatabaseConfig) error {
	var initErr error

	m.once.Do(func() {
		m.config = config

		// 验证配置
		if err := m.validateConfig(); err != nil {
			initErr = fmt.Errorf("config validation failed: %w", err)
			return
		}

		// 初始化数据库连接
		if err := m.initDatabase(); err != nil {
			initErr = fmt.Errorf("database initialization failed: %w", err)
			return
		}

		// 配置连接池
		if err := m.configureConnectionPool(); err != nil {
			initErr = fmt.Errorf("connection pool configuration failed: %w", err)
			return
		}

		logger.Info("Database connection initialized", logger.Any("type", config.Type))
	})

	return initErr
}

// GetDB 获取数据库实例
func (m *RDBStore) GetDB() *gorm.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.db
}

// initDatabase 初始化数据库连接
func (m *RDBStore) initDatabase() error {
	gormLogger := logger.NewGormLogger(logger.Default())
	gormConfig := &gorm.Config{
		Logger: gormLogger,
	}

	var err error

	switch m.config.Type {
	case MySQL:
		err = m.initMySQL(gormConfig)
	case SQLite:
		err = m.initSQLite(gormConfig)
	case PostgreSQL:
		err = m.initPostgreSQL(gormConfig)
	default:
		return fmt.Errorf("unsupported database type: %v", m.config.Type)
	}

	return err
}

// initMySQL 初始化MySQL连接
func (m *RDBStore) initMySQL(config *gorm.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.config.Username,
		m.config.Password,
		m.config.Host,
		m.config.Port,
		m.config.Database,
	)

	var err error
	m.db, err = gorm.Open(mysql.Open(dsn), config)
	return err
}

// initSQLite 初始化SQLite连接
func (m *RDBStore) initSQLite(config *gorm.Config) error {
	dbPath := m.config.File + ".db"
	var err error
	m.db, err = gorm.Open(sqlite.Open(dbPath), config)
	return err
}

// initPostgreSQL 初始化PostgreSQL连接
func (m *RDBStore) initPostgreSQL(config *gorm.Config) error {
	sslMode := m.config.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		m.config.Host,
		m.config.Username,
		m.config.Password,
		m.config.Database,
		m.config.Port,
		sslMode,
	)

	var err error
	m.db, err = gorm.Open(postgres.Open(dsn), config)
	return err
}

// configureConnectionPool 配置数据库连接池参数
func (m *RDBStore) configureConnectionPool() error {
	dbConn, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	// 设置最大打开连接数
	maxOpen, _ := strconv.Atoi(m.config.MaxOpenConnection)
	if maxOpen <= 0 {
		maxOpen = 25
	}
	dbConn.SetMaxOpenConns(maxOpen)

	// 设置最大空闲连接数
	maxIdle, _ := strconv.Atoi(m.config.MaxIdleConnection)
	if maxIdle <= 0 {
		maxIdle = 10
	}
	dbConn.SetMaxIdleConns(maxIdle)

	// 设置连接最大生存时间
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	// 设置连接最大空闲时间
	dbConn.SetConnMaxIdleTime(1 * time.Minute)

	// 测试连接
	if err := dbConn.Ping(); err != nil {
		return fmt.Errorf("database ping test failed: %w", err)
	}

	// 打印连接池统计信息
	stats := dbConn.Stats()
	logger.Info("Database connection pool stats",
		logger.Int("max_open", stats.MaxOpenConnections),
		logger.Int("open", stats.OpenConnections),
		logger.Int("in_use", stats.InUse),
		logger.Int("idle", stats.Idle),
	)

	return nil
}

// Close 关闭数据库连接
func (m *RDBStore) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		dbConn, err := m.db.DB()
		if err != nil {
			return err
		}
		return dbConn.Close()
	}
	return nil
}

// HealthCheck 执行数据库健康检查
func (m *RDBStore) HealthCheck() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	dbConn, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return dbConn.PingContext(ctx)
}

// MonitorConnectionPool 监控数据库连接池状态
func (m *RDBStore) MonitorConnectionPool(ctx context.Context) {
	if m.db == nil {
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dbConn, err := m.db.DB()
			if err != nil {
				logger.Error("Failed to get database connection for monitoring", logger.Err(err))
				continue
			}

			stats := dbConn.Stats()
			logger.Debug("Connection pool status",
				logger.Int("open", stats.OpenConnections),
				logger.Int("in_use", stats.InUse),
				logger.Int("idle", stats.Idle),
				logger.Int64("wait_count", stats.WaitCount),
			)

			// 如果等待队列过长则记录警告
			if stats.WaitCount > 10 {
				logger.Warn("High database connection pool wait queue", logger.Int64("wait_count", stats.WaitCount))
			}
		}
	}
}

// validateConfig 验证数据库配置
func (m *RDBStore) validateConfig() error {
	config := m.config

	switch config.Type {
	case MySQL, PostgreSQL:
		if config.Host == "" {
			return fmt.Errorf("database host is required for %s", m.getDBTypeName(config.Type))
		}
		if config.Port <= 0 {
			return fmt.Errorf("valid database port is required for %s", m.getDBTypeName(config.Type))
		}
		if config.Database == "" {
			return fmt.Errorf("database name is required for %s", m.getDBTypeName(config.Type))
		}
		if config.Username == "" {
			return fmt.Errorf("database username is required for %s", m.getDBTypeName(config.Type))
		}
	case SQLite:
		if config.File == "" {
			return fmt.Errorf("database file path is required for SQLite")
		}
	default:
		return fmt.Errorf("unsupported rdb type: %v", config.Type)
	}

	return nil
}

// IsInitialized 检查数据库是否已初始化
func (m *RDBStore) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.db != nil
}

// AutoMigrate 自动迁移表结构
func (m *RDBStore) AutoMigrate(models ...interface{}) error {
	if m.db == nil {
		return fmt.Errorf("database not initialized")
	}

	return m.db.AutoMigrate(models...)
}

// GetDatabaseType 获取数据库类型
func (m *RDBStore) GetDatabaseType() DBType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.Type
}

// GetConfigInfo 获取配置信息（隐藏敏感信息）
func (m *RDBStore) GetConfigInfo() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	info := map[string]interface{}{
		"type":                m.getDBTypeName(m.config.Type),
		"host":                m.config.Host,
		"port":                m.config.Port,
		"database":            m.config.Database,
		"username":            m.config.Username,
		"file":                m.config.File,
		"max_open_connection": m.config.MaxOpenConnection,
		"max_idle_connection": m.config.MaxIdleConnection,
		"ssl_mode":            m.config.SSLMode,
	}

	// 不暴露密码
	return info
}

// getDBTypeName 获取数据库类型名称
func (m *RDBStore) getDBTypeName(dbType DBType) string {
	switch dbType {
	case MySQL:
		return "mysql"
	case PostgreSQL:
		return "postgresql"
	case SQLite:
		return "sqlite"
	default:
		return "unknown"
	}
}
