// Package workspace provides work directory management.
package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager 工作目录管理器
type Manager struct {
	basePath string
}

// NewManager 创建管理器
func NewManager(basePath string) *Manager {
	return &Manager{basePath: basePath}
}

// Create 创建工作目录
func (m *Manager) Create(provider, region, resourceID, taskID string) (string, error) {
	path := m.GenPath(provider, region, resourceID, taskID)
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("failed to create workspace: %w", err)
	}
	return path, nil
}

// GenPath 生成工作目录路径
func (m *Manager) GenPath(provider, region, resourceID, taskID string) string {
	return filepath.Join(m.basePath, provider, region, resourceID, taskID)
}

// Clean 清理工作目录
func (m *Manager) Clean(path string) error {
	if path == "" || path == "/" {
		return fmt.Errorf("invalid path")
	}
	return os.RemoveAll(path)
}

// WriteFile 写入文件
func (m *Manager) WriteFile(dir, filename string, content []byte) error {
	path := filepath.Join(dir, filename)
	return os.WriteFile(path, content, 0644)
}

// ReadFile 读取文件
func (m *Manager) ReadFile(dir, filename string) ([]byte, error) {
	path := filepath.Join(dir, filename)
	return os.ReadFile(path)
}

// Exists 检查路径是否存在
func (m *Manager) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
