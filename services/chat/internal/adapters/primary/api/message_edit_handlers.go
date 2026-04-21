package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) editMessage(ctx *gin.Context) {
	currentUserID := userIDFromContext(ctx)

	messageID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
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

	chatID, err := h.csvc.EditMessage(currentUserID, messageIDUint, editMessageRequest.NewContent)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "message edit successfully",
	})

	// broadcast to all the connections
	participants, err := h.csvc.GetChatParticipants(chatID)
	if err != nil {
		// deletion succeeded; log and continue without broadcast
		log.Printf("deleteMessage: get participants for chat %d failed: %v", chatID, err)
		ctx.JSON(http.StatusOK, gin.H{"message": "message deleted successfully"})
		return
	}

	msg := &struct {
		Type       string `json:"type"`
		MessageID  uint   `json:"messageID"`
		NewContent string `json:"newContent"`
	}{
		Type:       "messageEdited",
		MessageID:  messageIDUint,
		NewContent: editMessageRequest.NewContent,
	}
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}
