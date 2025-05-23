package ws

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*websocket.Conn
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*websocket.Conn),
	}
}

func (m *ConnectionManager) Add(userID string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[userID] = conn
}

func (m *ConnectionManager) Remove(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.connections, userID)
}

func (m *ConnectionManager) Send(userID, message string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, ok := m.connections[userID]
	if !ok {
		return errors.New("не удалось найти подключение")
	}

	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}
