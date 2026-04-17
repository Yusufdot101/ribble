package api

import (
	"net/http"
	"sync"

	"github.com/Yusufdot101/ripple/services/chat/config"
	"github.com/gorilla/websocket"
)

func wsError(conn *websocket.Conn, msg string) {
	conn.WriteJSON(map[string]string{
		"type":    "error",
		"message": msg,
	})
	conn.Close()
}

type hub struct {
	mu sync.RWMutex

	clients map[uint]map[*websocket.Conn]bool
}

func newHub() *hub {
	return &hub{
		clients: make(map[uint]map[*websocket.Conn]bool),
	}
}

func (h *hub) addClient(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*websocket.Conn]bool)
	}

	h.clients[userID][conn] = true
}

func (h *hub) removeClient(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[userID] != nil {
		delete(h.clients[userID], conn)
		if len(h.clients[userID]) == 0 {
			delete(h.clients, userID)
		}
	}
}

func (h *hub) SendToUser(userID uint, msg any) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns := h.clients[userID]

	for conn := range conns {
		conn.WriteJSON(msg)
	}
}

var upgrader = websocket.Upgrader{
	// tighten it later
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == config.GetFrontendURL()
	},
}
