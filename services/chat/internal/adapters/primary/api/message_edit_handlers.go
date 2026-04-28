package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

func (h *handler) editMessage(ctx *gin.Context) {
	currentUserID := context.UserIDFromContext(ctx)

	messageID, err := strconv.ParseUint(ctx.Param("messageId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}
	if messageID > uint64(^uint(0)) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}
	messageIDUint := uint(messageID)

	var editMessageRequest struct {
		NewContent string `json:"newContent"`
	}
	if err := ctx.ShouldBind(&editMessageRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
		ctx.JSON(statusCode, gin.H{
			"error": error,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "message edit successfully",
	})

	// broadcast to all the connections
	participants, err := h.csvc.GetChatParticipants(message.ChatID, currentUserID)
	if err != nil {
		// edit succeeded; log and continue without broadcast
		log.Printf("editMessage: get participants for chat %d failed: %v", message.ChatID, err)
		return
	}

	msg := &struct {
		Type       string    `json:"type"`
		MessageID  uint      `json:"messageID"`
		NewContent string    `json:"newContent"`
		UpdatedAt  time.Time `json:"updatedAt"`
	}{
		Type:       "messageEdited",
		MessageID:  message.ID,
		NewContent: message.Content,
		UpdatedAt:  message.UpdatedAt,
	}
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}
