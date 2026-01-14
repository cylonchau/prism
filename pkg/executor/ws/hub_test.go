package ws

import (
	"encoding/json"
	"testing"
)

func TestHub_New(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub should not return nil")
	}
	if hub.clients == nil {
		t.Fatal("clients map should be initialized")
	}
}

func TestHub_RegisterUnregister(t *testing.T) {
	hub := NewHub()

	client := &Client{
		taskID: "task-1",
		send:   make(chan []byte, 10),
	}

	// Register
	hub.Register(client)

	hub.mu.RLock()
	if len(hub.clients["task-1"]) != 1 {
		t.Error("client should be registered")
	}
	hub.mu.RUnlock()

	// Unregister
	hub.Unregister(client)

	hub.mu.RLock()
	if len(hub.clients["task-1"]) != 0 {
		t.Error("client should be unregistered")
	}
	hub.mu.RUnlock()
}

func TestHub_BroadcastNoClients(t *testing.T) {
	hub := NewHub()

	// 没有客户端时发送不应 panic
	hub.Broadcast("task-1", &Message{
		Type:   TypeLog,
		TaskID: "task-1",
		Data:   "test",
	})
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()

	client := &Client{
		taskID: "task-1",
		send:   make(chan []byte, 10),
	}
	hub.Register(client)

	// 发送消息
	hub.Broadcast("task-1", &Message{
		Type:   TypeLog,
		TaskID: "task-1",
		Data:   "test message",
	})

	// 检查客户端是否收到
	select {
	case data := <-client.send:
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("should unmarshal message: %v", err)
		}
		if msg.Type != TypeLog {
			t.Errorf("type should be log, got %s", msg.Type)
		}
	default:
		t.Error("client should receive message")
	}

	hub.Unregister(client)
}

func TestHub_SendLog(t *testing.T) {
	hub := NewHub()

	client := &Client{
		taskID: "task-1",
		send:   make(chan []byte, 10),
	}
	hub.Register(client)

	hub.SendLog("task-1", "log message")

	select {
	case data := <-client.send:
		var msg Message
		json.Unmarshal(data, &msg)
		if msg.Type != TypeLog {
			t.Errorf("type should be log, got %s", msg.Type)
		}
		if msg.Data != "log message" {
			t.Errorf("data should be 'log message', got %v", msg.Data)
		}
	default:
		t.Error("client should receive log message")
	}

	hub.Unregister(client)
}

func TestHub_SendProgress(t *testing.T) {
	hub := NewHub()

	client := &Client{
		taskID: "task-1",
		send:   make(chan []byte, 10),
	}
	hub.Register(client)

	hub.SendProgress("task-1", &ProgressData{
		Phase:   "apply",
		Percent: 50,
		Message: "applying",
	})

	select {
	case data := <-client.send:
		var msg Message
		json.Unmarshal(data, &msg)
		if msg.Type != TypeProgress {
			t.Errorf("type should be progress, got %s", msg.Type)
		}
	default:
		t.Error("client should receive progress message")
	}

	hub.Unregister(client)
}

func TestHub_SendComplete(t *testing.T) {
	hub := NewHub()

	client := &Client{
		taskID: "task-1",
		send:   make(chan []byte, 10),
	}
	hub.Register(client)

	hub.SendComplete("task-1", true, map[string]string{"id": "123"})

	select {
	case data := <-client.send:
		var msg Message
		json.Unmarshal(data, &msg)
		if msg.Type != TypeComplete {
			t.Errorf("type should be complete, got %s", msg.Type)
		}
	default:
		t.Error("client should receive complete message")
	}

	hub.Unregister(client)
}

func TestHub_MultipleClients(t *testing.T) {
	hub := NewHub()

	client1 := &Client{taskID: "task-1", send: make(chan []byte, 10)}
	client2 := &Client{taskID: "task-1", send: make(chan []byte, 10)}
	client3 := &Client{taskID: "task-2", send: make(chan []byte, 10)}

	hub.Register(client1)
	hub.Register(client2)
	hub.Register(client3)

	// 发送到 task-1
	hub.SendLog("task-1", "message")

	// client1 和 client2 应该收到
	select {
	case <-client1.send:
	default:
		t.Error("client1 should receive message")
	}
	select {
	case <-client2.send:
	default:
		t.Error("client2 should receive message")
	}
	// client3 不应收到
	select {
	case <-client3.send:
		t.Error("client3 should not receive message for task-1")
	default:
	}

	hub.Unregister(client1)
	hub.Unregister(client2)
	hub.Unregister(client3)
}
