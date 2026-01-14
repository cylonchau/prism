package store

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	models "github.com/cylonchau/prism/pkg/model"
)

func TestRDBStore_SQLite(t *testing.T) {
	dbFile := "test_rdb_store.db"
	defer os.Remove(dbFile)

	store := NewRDBStore()

	config := DatabaseConfig{
		Type: SQLite,
		File: "test_rdb_store",
	}

	// Test Initialize
	if err := store.Initialize(config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Test IsInitialized
	if !store.IsInitialized() {
		t.Error("Store should be initialized")
	}

	// Test GetDB
	db := store.GetDB()
	if db == nil {
		t.Error("GetDB should return non-nil")
	}

	// Test GetDatabaseType
	if store.GetDatabaseType() != SQLite {
		t.Errorf("Expected SQLite, got %v", store.GetDatabaseType())
	}

	// Test AutoMigrate
	if err := store.AutoMigrate(&models.ExecutionTask{}); err != nil {
		t.Errorf("AutoMigrate failed: %v", err)
	}

	// Test HealthCheck
	if err := store.HealthCheck(); err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}

	// Test Close
	if err := store.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestRDBStore_Singleton(t *testing.T) {
	dbFile := "test_singleton.db"
	defer os.Remove(dbFile)

	// Reset singleton for testing
	instance = nil
	initOnce = *new(sync.Once)

	store1 := GetInstance()
	store2 := GetInstance()

	if store1 != store2 {
		t.Error("GetInstance should return same instance")
	}
}

func TestRDBStore_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config DatabaseConfig
	}{
		{
			name: "MySQL without host",
			config: DatabaseConfig{
				Type:     MySQL,
				Database: "test",
			},
		},
		{
			name: "PostgreSQL without database",
			config: DatabaseConfig{
				Type: PostgreSQL,
				Host: "localhost",
				Port: 5432,
			},
		},
		{
			name: "SQLite without file",
			config: DatabaseConfig{
				Type: SQLite,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new instance for each test to avoid sync.Once interference
			store := NewRDBStore()
			err := store.Initialize(tt.config)
			if err == nil {
				t.Error("Expected error for invalid config")
			}
		})
	}
}

func TestRDBStore_MonitorConnectionPool(t *testing.T) {
	dbFile := "test_monitor.db"
	defer os.Remove(dbFile)

	store := NewRDBStore()
	config := DatabaseConfig{
		Type: SQLite,
		File: "test_monitor",
	}

	if err := store.Initialize(config); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	defer store.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Should not panic
	go store.MonitorConnectionPool(ctx)

	<-ctx.Done()
}

func TestNewRDBStore(t *testing.T) {
	store := NewRDBStore()
	if store == nil {
		t.Error("NewRDBStore should return non-nil")
	}

	if store.IsInitialized() {
		t.Error("New store should not be initialized")
	}
}

func TestRDBStore_MultipleInitialize(t *testing.T) {
	dbFile := "test_multi_init.db"
	defer os.Remove(dbFile)

	store := NewRDBStore()
	config := DatabaseConfig{
		Type: SQLite,
		File: "test_multi_init",
	}

	// First initialize
	if err := store.Initialize(config); err != nil {
		t.Fatalf("First initialize failed: %v", err)
	}

	// Second initialize should be no-op (due to sync.Once)
	if err := store.Initialize(config); err != nil {
		t.Errorf("Second initialize failed: %v", err)
	}

	store.Close()
}
