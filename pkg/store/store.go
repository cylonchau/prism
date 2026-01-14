package store

import (
	"context"

	"gorm.io/gorm"
)

// Store defines the unified interface for all storage backends (RDB, Redis, etc.).
type Store interface {
	// Initialize initializes the store with given configuration
	Initialize(config DatabaseConfig) error

	// GetDB returns the underlying database connection (for RDB stores)
	// Returns nil for non-RDB stores like Redis
	GetDB() *gorm.DB

	// Close closes the store connection
	Close() error

	// HealthCheck performs a health check on the store
	HealthCheck() error

	// AutoMigrate runs database migrations (for RDB stores)
	AutoMigrate(models ...interface{}) error

	// GetDatabaseType returns the type of database
	GetDatabaseType() DBType

	// IsInitialized checks if the store is initialized
	IsInitialized() bool

	// MonitorConnectionPool monitors connection pool status (for RDB stores)
	MonitorConnectionPool(ctx context.Context)
}
