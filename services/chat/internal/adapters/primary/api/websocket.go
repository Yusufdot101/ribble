package api

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/config"
	"github.com/Yusufdot101/ripple/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

func wsError(conn *websocket.Conn, msg string) {
	err := conn.WriteJSON(map[string]string{
		"type":    "error",
		"message": msg,
	})
	if err != nil {
		log.Println("error writing json: ", err)
		return
	}
	err = conn.Close()
	if err != nil {
		log.Println("error writing json: ", err)
		return
	}
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
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Println("error writing json: ", err)
		}
	}
}

var upgrader = websocket.Upgrader{
	// tighten it later
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == config.GetFrontendURL()
	},
}

func (h *handler) authenticateWS(ctx *gin.Context) (*websocket.Conn, uint, error) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, 0, err
	}

	var authMsg struct {
		Token string `json:"token"`
	}

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	if err := conn.ReadJSON(&authMsg); err != nil {
		return nil, 0, err
	}

	_ = conn.SetReadDeadline(time.Time{})

	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// validate JWT
	token, err := middleware.ValidateJWT(authMsg.Token)
	if err != nil {
		wsError(conn, middleware.ErrInvalidJWT.Error())
		return nil, 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		wsError(conn, middleware.ErrInvalidJWT.Error())
		return nil, 0, middleware.ErrInvalidJWT
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok || userIDStr == "" {
		wsError(conn, middleware.ErrInvalidJWT.Error())
		return nil, 0, middleware.ErrInvalidJWT
	}

	userIDInt, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		wsError(conn, middleware.ErrInvalidJWT.Error())
		return nil, 0, middleware.ErrInvalidJWT
	}

	return conn, uint(userIDInt), nil
}
