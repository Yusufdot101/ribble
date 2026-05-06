package api

import (
	"net/http"
	"strconv"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

type Permissions struct {
	Permissions []*domain.Permission `json:"permissions"`
}

func (h *handler) getUserPermissions(c *gin.Context) {
	currentUserID := context.UserIDFromContext(c)

	chatID, err := strconv.ParseUint(c.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}
	chatIDUint := uint(chatID)

	permissions, err := h.csvc.GetUserPermissions(chatIDUint, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: "error getting permissions",
		})
		return
	}

	c.JSON(http.StatusOK, response.Response[Permissions]{
		Data: Permissions{
			Permissions: permissions,
		},
	})
}
