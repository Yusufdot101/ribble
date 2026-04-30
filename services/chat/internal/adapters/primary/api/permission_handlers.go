package api

import (
	"net/http"
	"strconv"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/gin-gonic/gin"
)

func (h *handler) getUserPermissions(c *gin.Context) {
	currentUserID := context.UserIDFromContext(c)

	chatID, err := strconv.ParseUint(c.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid chat id",
		})
		return
	}
	chatIDUint := uint(chatID)

	permissions, err := h.csvc.GetUserPermissions(chatIDUint, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error getting permissions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
	})
}
