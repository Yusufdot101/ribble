package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) deleteMessage(ctx *gin.Context) {
	currentUserID := userIDFromContext(ctx)
	messageID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}

	chatID, err := h.csvc.DeleteMessage(currentUserID, uint(messageID))
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
	participants, err := h.csvc.GetChatParticipants(chatID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	msg := &struct {
		Type      string `json:"type"`
		MessageID uint   `json:"messageID"`
	}{
		Type:      "messageDeleted",
		MessageID: uint(messageID),
	}
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}
