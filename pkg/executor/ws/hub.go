// Package ws provides WebSocket real-time messaging.
package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageType 消息类型
type MessageType string

const (
	TypeLog      MessageType = "log"
	TypeProgress MessageType = "progress"
	TypeComplete MessageType = "complete"
	TypeError    MessageType = "error"
)

// Message WebSocket 消息
type Message struct {
	Type   MessageType `json:"type"`
	TaskID string      `json:"task_id"`
	Data   interface{} `json:"data"`
	Time   time.Time   `json:"time"`
}

// ProgressData 进度数据
type ProgressData struct {
	Phase   string `json:"phase"`
	Percent int    `json:"percent"`
	Elapsed int64  `json:"elapsed"`
	Message string `json:"message"`
}

// Client WebSocket 客户端
type Client struct {
	conn   *websocket.Conn
	taskID string
	send   chan []byte
}

// Hub 连接管理中心
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]bool // taskID -> clients
}

// NewHub 创建 Hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*Client]bool),
	}
}

// Register 注册客户端
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.taskID] == nil {
		h.clients[client.taskID] = make(map[*Client]bool)
	}
	h.clients[client.taskID][client] = true
}

// Unregister 注销客户端
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[client.taskID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.clients, client.taskID)
		}
	}
	close(client.send)
}

// Broadcast 广播消息
func (h *Hub) Broadcast(taskID string, msg *Message) {
	h.mu.RLock()
	clients := h.clients[taskID]
	h.mu.RUnlock()

	if clients == nil {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for client := range clients {
		select {
		case client.send <- data:
		default:
			// 缓冲区满，跳过
		}
	}
}

// SendLog 发送日志消息
func (h *Hub) SendLog(taskID, message string) {
	h.Broadcast(taskID, &Message{
		Type:   TypeLog,
		TaskID: taskID,
		Data:   message,
		Time:   time.Now(),
	})
}

// SendProgress 发送进度消息
func (h *Hub) SendProgress(taskID string, data *ProgressData) {
	h.Broadcast(taskID, &Message{
		Type:   TypeProgress,
		TaskID: taskID,
		Data:   data,
		Time:   time.Now(),
	})
}

// SendComplete 发送完成消息
func (h *Hub) SendComplete(taskID string, success bool, result interface{}) {
	h.Broadcast(taskID, &Message{
		Type:   TypeComplete,
		TaskID: taskID,
		Data: map[string]interface{}{
			"success": success,
			"result":  result,
		},
		Time: time.Now(),
	})
}
