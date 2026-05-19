package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MessagePayload struct {
	ClientID string `json:"clientId"`
	Content  string `json:"content"`
	ChatID   uint   `json:"chatId"`
}

func (h *handler) handleMessage(conn *websocket.Conn, userID uint, msg incomingMsg) error {
	var p MessagePayload
	err := json.Unmarshal(msg.Payload, &p)
	if err != nil {
		return domain.ErrInvalidPayload
	}
	log.Println("msg: ", msg)

	payload := map[string]any{
		"clientId": p.ClientID,
	}
	if p.ChatID == 0 || strings.TrimSpace(p.Content) == "" || msg.Type != "newMessage" {
		res := outgoingMsg{
			Type:    "nack",
			Payload: payload,
			Message: "invalid request",
		}
		_ = conn.WriteJSON(res)
		return nil
	}

	userHasPermission, err := h.csvc.UserHasPermission(userID, p.ChatID, domain.SendMessage)
	if err != nil {
		res := outgoingMsg{
			Type:    "nack",
			Payload: payload,
			Message: "not permitted",
		}
		_ = conn.WriteJSON(res)
		return nil
	}

	if !userHasPermission {
		res := outgoingMsg{
			Type:    "nack",
			Payload: payload,
			Message: "not allowed to write messages",
		}
		_ = conn.WriteJSON(res)
		return nil
	}

	participants, err := h.csvc.GetChatParticipants(p.ChatID, userID)
	if err != nil {
		res := outgoingMsg{
			Type:    "nack",
			Payload: payload,
			Message: "chat not found",
		}
		_ = conn.WriteJSON(res)
		return nil
	}

	if !userIsInChat(userID, participants) {
		res := outgoingMsg{
			Type:    "nack",
			Payload: payload,
			Message: "not a participant of this chat",
		}
		_ = conn.WriteJSON(res)
		return fmt.Errorf("user not in chat: %w", err)
	}

	message, err := h.csvc.NewMessage(userID, p.ChatID, p.Content, domain.StandardMessage)
	if err != nil {
		res := outgoingMsg{
			Type:    "nack",
			Payload: payload,
			Message: "failed to send message",
		}
		_ = conn.WriteJSON(res)
		return nil
	}

	payload["id"] = message.ID
	res := outgoingMsg{
		Type:    "ack",
		Payload: payload,
		Message: "message delivered",
	}

	_ = conn.WriteJSON(res)

	res.Type = "message"
	res.Payload["message"] = message
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, res)
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
