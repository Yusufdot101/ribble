package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) deleteMessage(ctx *gin.Context) {
	currentUserID := context.UserIDFromContext(ctx)

	chatID, err := strconv.ParseUint(ctx.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}
	chatIDUint := uint(chatID)

	messageID, err := strconv.ParseUint(ctx.Param("messageId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid message id",
		})
		return
	}
	if messageID > uint64(^uint(0)) {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid message id",
		})
		return
	}
	messageIDUint := uint(messageID)

	message, err := h.csvc.DeleteMessage(chatIDUint, currentUserID, messageIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response[any]{
		Message: "message deleted successfully",
	})

	// broadcast to all the connections
	participants, err := h.csvc.GetChatParticipants(chatIDUint, currentUserID)
	if err != nil {
		// deletion succeeded; log and continue without broadcast
		log.Printf("deleteMessage: get participants for chat %d failed: %v", chatID, err)
		return
	}

	msg := outgoingMsg{
		Type:    "message",
		SubType: "messageDeleted",
		Message: "message deleted successfully",
		Payload: map[string]any{
			"id":        message.ID,
			"content":   message.Content,
			"updatedAt": message.UpdatedAt,
			"deleted":   true,
		},
	}
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}
