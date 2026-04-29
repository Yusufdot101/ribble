package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

type websocketMsg struct {
	Type     string `json:"type"`
	ChatID   uint   `json:"chatID"`
	Content  string `json:"content"`
	ClientID string `json:"clientID"`
	SenderID uint   `json:"senderID"`
}

func (h *handler) readLoop(conn *websocket.Conn, userID uint) {
	for {
		var msg websocketMsg
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("error: ", err)
			break
		}

		if err := h.handleMessage(conn, userID, msg); err != nil {
			log.Println("message error:", err)
		}
	}
}

func (h *handler) handleMessage(conn *websocket.Conn, userID uint, msg websocketMsg) error {
	if msg.ChatID == 0 || strings.TrimSpace(msg.Content) == "" {
		_ = conn.WriteJSON(map[string]string{
			"type":     "nack",
			"message":  "invalid message",
			"clientID": msg.ClientID,
		})
		return nil
	}

	userHasPermission, err := h.csvc.UserHasPermission(userID, msg.ChatID, domain.SendMessage)
	if err != nil {
		_ = conn.WriteJSON(map[string]string{
			"type":     "nack",
			"message":  err.Error(),
			"clientID": msg.ClientID,
		})
		return nil
	}

	if !userHasPermission {
		_ = conn.WriteJSON(map[string]string{
			"type":     "nack",
			"message":  "not allowed to write messages",
			"clientID": msg.ClientID,
		})
		return nil
	}

	participants, err := h.csvc.GetChatParticipants(msg.ChatID, userID)
	if err != nil {
		_ = conn.WriteJSON(map[string]string{
			"type":     "nack",
			"message":  "chat not found",
			"clientID": msg.ClientID,
		})
		return nil
	}

	if !userIsInChat(userID, participants) {
		_ = conn.WriteJSON(map[string]string{
			"type":     "nack",
			"message":  "not a participant of this chat",
			"clientID": msg.ClientID,
		})
		return fmt.Errorf("user not in chat")
	}

	message, err := h.csvc.NewMessage(userID, msg.ChatID, msg.Content)
	if err != nil {
		_ = conn.WriteJSON(map[string]string{
			"type":     "nack",
			"clientID": msg.ClientID,
			"message":  "failed to send message",
		})
		return nil
	}

	_ = conn.WriteJSON(map[string]string{
		"type":     "ack",
		"clientID": msg.ClientID,
		"message":  "message delivered",
	})

	for _, p := range participants {
		if p.UserID == userID {
			continue
		}
		h.hub.SendToUser(p.UserID, message)
	}

	return nil
}

func (h *handler) getMessages(ctx *gin.Context) {
	chatID, err := strconv.ParseUint(ctx.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid chat id",
		})
		return
	}

	// make sure the user is in the chat
	currentUserID := context.UserIDFromContext(ctx)
	participants, err := h.csvc.GetChatParticipants(uint(chatID), currentUserID)
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

	messages, err := h.csvc.GetMessages(uint(chatID), domain.GetMessageFilter{})
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

func (h *handler) syncMessages(ctx *gin.Context) {
	chatID, err := strconv.ParseUint(ctx.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid chat id",
		})
		return
	}

	lastMessageID, err := strconv.ParseUint(ctx.Query("lastMessageId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}

	// make sure the user is in the chat
	currentUserID := context.UserIDFromContext(ctx)
	participants, err := h.csvc.GetChatParticipants(uint(chatID), currentUserID)
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

	messages, err := h.csvc.GetMessages(uint(chatID), domain.GetMessageFilter{
		LastMessageID: uint(lastMessageID),
	})
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
