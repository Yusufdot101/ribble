package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/gin-gonic/gin"
)

func (h *handler) deleteMessage(ctx *gin.Context) {
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

	chatID, err := h.csvc.DeleteMessage(currentUserID, messageIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "message deleted successfully",
	})

	// broadcast to all the connections
	participants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		// deletion succeeded; log and continue without broadcast
		log.Printf("deleteMessage: get participants for chat %d failed: %v", chatID, err)
		return
	}

	msg := &struct {
		Type      string `json:"type"`
		MessageID uint   `json:"messageID"`
	}{
		Type:      "messageDeleted",
		MessageID: messageIDUint,
	}
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}
