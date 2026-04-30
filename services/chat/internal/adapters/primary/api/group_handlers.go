package api

import (
	"net/http"
	"strconv"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

func (h *handler) addToGroup(c *gin.Context) {
	var req struct {
		UserIDs []uint `json:"userIDs"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}
	currentUserID := context.UserIDFromContext(c)

	chatID, err := strconv.ParseUint(c.Param("chatId"), 10, strconv.IntSize)
	chatIDUint := uint(chatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid chat id",
		})
		return
	}

	userHasPermission, err := h.csvc.UserHasPermission(currentUserID, chatIDUint, domain.AddToGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "not permitted",
		})
		return
	}

	if !userHasPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not permitted",
		})
		return
	}

	err = h.csvc.AddUsersToGroup(chatIDUint, req.UserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "users added to group",
	})
}
