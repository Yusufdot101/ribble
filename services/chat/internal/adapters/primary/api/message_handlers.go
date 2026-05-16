package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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
			"message":  "not permitted",
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

	message, err := h.csvc.NewMessage(userID, msg.ChatID, msg.Content, domain.StandardMessage)
	if err != nil {
		_ = conn.WriteJSON(map[string]string{
			"type":     "nack",
			"clientID": msg.ClientID,
			"message":  "failed to send message",
		})
		return nil
	}

	_ = conn.WriteJSON(map[string]any{
		"type":     "ack",
		"clientID": msg.ClientID,
		"message":  "message delivered",
		"id":       message.ID,
	})

	message.ClientID = msg.ClientID

	for _, p := range participants {
		h.hub.SendToUser(p.UserID, message)
	}

	return nil
}

type Messages struct {
	Messages []*domain.Message `json:"messages"`
}

func (h *handler) getMessages(ctx *gin.Context) {
	chatID, err := strconv.ParseUint(ctx.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}

	// make sure the user is in the chat
	currentUserID := context.UserIDFromContext(ctx)
	participants, err := h.csvc.GetChatParticipants(uint(chatID), currentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	if !userIsInChat(currentUserID, participants) {
		ctx.JSON(http.StatusNotFound, response.Response[any]{
			Error: "chat not found",
		})
		return
	}

	messages, err := h.csvc.GetMessages(uint(chatID), domain.GetMessageFilter{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response[Messages]{
		Data: Messages{
			Messages: messages,
		},
	})
}

func (h *handler) syncMessages(ctx *gin.Context) {
	chatID, err := strconv.ParseUint(ctx.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}

	lastMessageID, err := strconv.ParseUint(ctx.Query("lastMessageId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid message id",
		})
		return
	}

	// make sure the user is in the chat
	currentUserID := context.UserIDFromContext(ctx)
	participants, err := h.csvc.GetChatParticipants(uint(chatID), currentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	if !userIsInChat(currentUserID, participants) {
		ctx.JSON(http.StatusNotFound, response.Response[any]{
			Error: "chat not found",
		})
		return
	}

	messages, err := h.csvc.GetMessages(uint(chatID), domain.GetMessageFilter{
		LastMessageID: uint(lastMessageID),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response[Messages]{
		Data: Messages{
			Messages: messages,
		},
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
