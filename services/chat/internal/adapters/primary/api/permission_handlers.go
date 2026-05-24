package api

import (
	"errors"
	"fmt"
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
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid request",
		})
		return
	}

	chatID, err := parameter.GetParameterValueUint(c, "chatId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	role := parameter.GetParameterValue(c, "role")

	if role == "" {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: fmt.Sprintf("invalid role: %s", role),
		})
		return
	}

	currentUserID := context.UserIDFromContext(c)

	participants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) || errors.Is(err, domain.ErrNotPermitted) {
			c.JSON(http.StatusForbidden, response.Response[any]{
				Error: "not a participant of this chat",
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Response[any]{
				Error: "error getting chat participants",
			})
		}
		return
	}

	if !userIsInChat(currentUserID, participants) {
		c.JSON(http.StatusForbidden, response.Response[any]{
			Error: "not a participant of this chat",
		})
		return
	}

	switch req.Action {
	case "grant":
		err = h.csvc.GrantRolePermission(currentUserID, chatID, role, req.Permission)
	case "revoke":
		err = h.csvc.RevokeRolePermission(currentUserID, chatID, role, req.Permission)
	default:
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: fmt.Sprintf("invalid action: %s", req.Action),
		})
		return
	}

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
		"role":       role,
		"permission": req.Permission,
		"action":     req.Action,
	}
	msg := outgoingMsg{
		Type:    "message",
		SubType: "updatedRolePermissions",
		Payload: payload,
	}
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}

type GetRolePermissions struct {
	RolePermissions map[string][]*domain.Permission `json:"rolePermissions"`
}

func (h *handler) getRolePermissions(c *gin.Context) {
	var req struct {
		Roles []string `json:"roles" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid roles",
		})
		return
	}
	if len(req.Roles) == 0 {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "roles cannot be empty",
		})
		return
	}

	chatID, err := parameter.GetParameterValueUint(c, "chatId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}

	currentUserID := context.UserIDFromContext(c)
	participants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) || errors.Is(err, domain.ErrNotPermitted) {
			c.JSON(http.StatusForbidden, response.Response[any]{
				Error: "not a participant of this chat",
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.Response[any]{
				Error: "error getting chat participants",
			})
		}
		return
	}

	if !userIsInChat(currentUserID, participants) {
		c.JSON(http.StatusForbidden, response.Response[any]{
			Error: "not a participant of this chat",
		})
		return
	}

	rolePermissions := make(map[string][]*domain.Permission)
	for _, role := range req.Roles {
		permissions, err := h.csvc.GetRolePermissions(chatID, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Response[any]{
				Error: "error getting permissions",
			})
			return
		}
		rolePermissions[role] = permissions
	}

	c.JSON(http.StatusOK, response.Response[GetRolePermissions]{
		Data: GetRolePermissions{
			RolePermissions: rolePermissions,
		},
	})
}
