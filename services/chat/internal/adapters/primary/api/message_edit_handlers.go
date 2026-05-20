package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

func (h *handler) editMessage(ctx *gin.Context) {
	currentUserID := context.UserIDFromContext(ctx)

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

	var editMessageRequest struct {
		NewContent string `json:"newContent"`
	}
	if err := ctx.ShouldBind(&editMessageRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	message, err := h.csvc.EditMessage(currentUserID, messageIDUint, editMessageRequest.NewContent)
	if err != nil {
		statusCode := http.StatusInternalServerError
		error := "error editing message"
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			statusCode = http.StatusNotFound
			error = err.Error()
		case errors.Is(err, domain.ErrInvalidMessageContent):
			statusCode = http.StatusBadRequest
			error = err.Error()
		case errors.Is(err, domain.ErrUpdateWindowOver):
			statusCode = http.StatusForbidden
			error = err.Error()
		}
		ctx.JSON(statusCode, response.Response[any]{
			Error: error,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response[any]{
		Message: "message edited successfully",
	})

	// broadcast to all the connections
	participants, err := h.csvc.GetChatParticipants(message.ChatID, currentUserID)
	if err != nil {
		// edit succeeded; log and continue without broadcast
		log.Printf("editMessage: get participants for chat %d failed: %v", message.ChatID, err)
		return
	}

	msg := outgoingMsg{
		Type:    "message",
		SubType: "messageEdited",
		Message: "message edited successfully",
		Payload: map[string]any{
			"id":        message.ID,
			"content":   message.Content,
			"updatedAt": message.UpdatedAt,
		},
	}
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}
