package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManager_New(t *testing.T) {
	m := NewManager("/tmp/test")
	if m == nil {
		t.Fatal("NewManager should not return nil")
	}
	if m.basePath != "/tmp/test" {
		t.Errorf("basePath should be /tmp/test, got %s", m.basePath)
	}
}

func TestManager_GenPath(t *testing.T) {
	m := NewManager("/tmp/terraform")

	path := m.GenPath("aws", "us-east-1", "resource-1", "task-123")
	expected := "/tmp/terraform/aws/us-east-1/resource-1/task-123"

	if path != expected {
		t.Errorf("GenPath = %s, want %s", path, expected)
	}
}

func TestManager_CreateAndClean(t *testing.T) {
	m := NewManager("/tmp/test-workspace")
	defer os.RemoveAll("/tmp/test-workspace")

	// Create
	path, err := m.Create("test", "region", "res-1", "task-1")
	if err != nil {
		t.Fatalf("Create should succeed: %v", err)
	}

	if !m.Exists(path) {
		t.Error("directory should exist after Create")
	}

	// Clean
	err = m.Clean(path)
	if err != nil {
		t.Fatalf("Clean should succeed: %v", err)
	}

	if m.Exists(path) {
		t.Error("directory should not exist after Clean")
	}
}

func TestManager_CleanInvalidPath(t *testing.T) {
	m := NewManager("/tmp/test")

	err := m.Clean("")
	if err == nil {
		t.Error("Clean empty path should return error")
	}

	err = m.Clean("/")
	if err == nil {
		t.Error("Clean root path should return error")
	}
}

func TestManager_WriteAndReadFile(t *testing.T) {
	m := NewManager("/tmp/test-workspace")
	defer os.RemoveAll("/tmp/test-workspace")

	dir, _ := m.Create("test", "region", "res-1", "task-1")

	// Write
	content := []byte("test content")
	err := m.WriteFile(dir, "test.txt", content)
	if err != nil {
		t.Fatalf("WriteFile should succeed: %v", err)
	}

	// Read
	data, err := m.ReadFile(dir, "test.txt")
	if err != nil {
		t.Fatalf("ReadFile should succeed: %v", err)
	}

	if string(data) != "test content" {
		t.Errorf("content should be 'test content', got %s", string(data))
	}
}

func TestManager_ReadNonExistentFile(t *testing.T) {
	m := NewManager("/tmp/test-workspace")

	_, err := m.ReadFile("/tmp/nonexistent", "file.txt")
	if err == nil {
		t.Error("ReadFile nonexistent should return error")
	}
}

func TestManager_Exists(t *testing.T) {
	m := NewManager("/tmp")

	// 存在的目录
	if !m.Exists("/tmp") {
		t.Error("/tmp should exist")
	}

	// 不存在的目录
	if m.Exists("/nonexistent_path_12345") {
		t.Error("nonexistent path should not exist")
	}
}

func TestManager_WriteFileCreatesDir(t *testing.T) {
	m := NewManager("/tmp/test-workspace")
	defer os.RemoveAll("/tmp/test-workspace")

	dir := filepath.Join("/tmp/test-workspace", "nested", "dir")
	os.MkdirAll(dir, 0755)

	err := m.WriteFile(dir, "test.txt", []byte("content"))
	if err != nil {
		t.Fatalf("WriteFile with nested dir should succeed: %v", err)
	}
}
