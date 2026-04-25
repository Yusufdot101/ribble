package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (h *handler) newMessage(ctx *gin.Context) {
	conn, userID, err := h.authenticateWS(ctx)
	if err != nil {
		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	h.hub.addClient(userID, conn)
	defer h.hub.removeClient(userID, conn)

	h.readLoop(conn, userID)
}

func (h *handler) readLoop(conn *websocket.Conn, userID uint) {
	for {
		var msg struct {
			Type    string `json:"type"`
			ChatID  uint   `json:"chatID"`
			Content string `json:"content"`
		}

		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		if err := h.handleMessage(conn, userID, msg); err != nil {
			log.Println("message error:", err)
		}
	}
}

func (h *handler) handleMessage(conn *websocket.Conn, userID uint, msg struct {
	Type    string `json:"type"`
	ChatID  uint   `json:"chatID"`
	Content string `json:"content"`
},
) error {
	if msg.ChatID == 0 || strings.TrimSpace(msg.Content) == "" {
		_ = conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": "invalid message",
		})
		return nil
	}

	userHasPermission, err := h.csvc.UserHasPermission(userID, msg.ChatID, domain.SendMessage)
	if err != nil {
		_ = conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": err.Error(),
		})
		return nil
	}

	if !userHasPermission {
		_ = conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": "not allowed to write messages",
		})
		return nil
	}

	participants, err := h.csvc.GetChatParticipants(msg.ChatID)
	if err != nil {
		_ = conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": "chat not found",
		})
		return nil
	}

	if !userIsInChat(userID, participants) {
		_ = conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": "not a participant of this chat",
		})
		return fmt.Errorf("user not in chat")
	}

	message, err := h.csvc.NewMessage(userID, msg.ChatID, msg.Content)
	if err != nil {
		_ = conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": "failed to send message",
		})
		return nil
	}

	for _, p := range participants {
		// TODO: add permission check
		h.hub.SendToUser(p.UserID, message)
	}

	return nil
}

func (h *handler) getMessages(ctx *gin.Context) {
	var GetMessagesRequest struct {
		ChatID uint `json:"chatID"`
	}
	if err := ctx.ShouldBind(&GetMessagesRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// make sure the user is in the chat
	currentUserID := context.UserIDFromContext(ctx)
	participants, err := h.csvc.GetChatParticipants(GetMessagesRequest.ChatID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !userIsInChat(currentUserID, participants) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "chat not found",
		})
		return
	}

	messages, err := h.csvc.GetMessages(GetMessagesRequest.ChatID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func userIsInChat(userID uint, participants []*domain.ChatParticipant) bool {
	userInChat := false
	for _, p := range participants {
		if p.UserID == uint(userID) {
			userInChat = true
			break
		}
	}
	return userInChat
}
