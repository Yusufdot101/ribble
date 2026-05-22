package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/parameter"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

type Permissions struct {
	Permissions []*domain.Permission `json:"permissions"`
}

func (h *handler) getUserPermissions(c *gin.Context) {
	currentUserID := context.UserIDFromContext(c)

	chatID, err := parameter.GetParameterValueUint(c, "chatId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}

	permissions, err := h.csvc.GetUserPermissions(chatID, currentUserID)
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

func (h *handler) updateRolePermission(c *gin.Context) {
	var req struct {
		Permission string `json:"permission" binding:"required"`
		Action     string `json:"action" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		log.Println("log1: ", err)
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid request",
		})
		return
	}

	chatID, err := parameter.GetParameterValueUint(c, "chatId")
	log.Println("log2: ", err)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	role := parameter.GetParameterValue(c, "role")

	if role == "" {
		log.Println("log3: ", role)
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: fmt.Sprintf("invalid role: %s", role),
		})
		return
	}

	currentUserID := context.UserIDFromContext(c)

	participants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	log.Println("log4: ", err)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}

	if !userIsInChat(currentUserID, participants) {
		log.Println("log5: ")
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: fmt.Sprintf("invalid chat id: %d", chatID),
		})
		return
	}

	switch req.Action {
	case "grant":
		err = h.csvc.GrantRolePermission(currentUserID, chatID, role, req.Permission)
	case "revoke":
		err = h.csvc.RevokeRolePermission(currentUserID, chatID, role, req.Permission)
	default:
		log.Println("log6: ")
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: fmt.Sprintf("invalid action: %s", req.Action),
		})
		return
	}

	log.Println("log7: ", err)
	if err != nil {
		status := http.StatusInternalServerError
		error := "an error occured"
		if errors.Is(err, domain.ErrInvalidAction) {
			status = http.StatusBadRequest
			error = err.Error()
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			status = http.StatusNotFound
			error = err.Error()
		} else if errors.Is(err, domain.ErrNotPermitted) {
			status = http.StatusForbidden
			error = err.Error()
		}
		c.JSON(status, response.Response[any]{
			Error: error,
		})
		return
	}

	c.JSON(http.StatusOK, response.Response[any]{
		Message: "role permission updated successfully",
	})

	payload := map[string]any{
		"message":    "role permission updated",
		"permission": req.Permission,
		"action":     req.Action,
	}
	msg := outgoingMsg{
		Type:    "updatedUserRole",
		Payload: payload,
	}
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}
