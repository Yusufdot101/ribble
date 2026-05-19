package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/parameter"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

func (h *handler) addToGroup(c *gin.Context) {
	var req struct {
		UserIDs []uint `json:"userIDs"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid user id",
		})
		return
	}
	currentUserID := context.UserIDFromContext(c)

	chatID, err := parameter.GetParameterValueUint(c, "chatId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	err = h.csvc.AddUsersToGroup(chatID, currentUserID, req.UserIDs)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotPermitted) {
			status = http.StatusForbidden
		}
		c.JSON(status, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.Response[any]{
		Message: "users added to group",
	})

	users, err := h.csvc.SearchUsers("", []uint{currentUserID})
	if err != nil {
		log.Printf("error getting current user: %v\n", err)
		return
	}
	if len(users) == 0 {
		log.Printf("current user not found: %d\n", currentUserID)
		return
	}
	currentUser := users[0]

	addedUsers, err := h.csvc.SearchUsers("", req.UserIDs)
	if err != nil {
		log.Printf("error getting added users: %v\n", err)
		return
	}

	names := make([]string, 0, len(addedUsers))
	for _, u := range addedUsers {
		names = append(names, u.Name)
	}
	usernames := strings.Join(names, ", ")

	message, err := h.csvc.NewMessage(currentUserID, chatID, fmt.Sprintf("%s added %s", currentUser.Name, usernames), domain.SystemMessage)
	if err != nil {
		log.Printf("error sending system message: %v\n", err)
		return
	}

	participants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		log.Printf("error getting chat participants: %v\n", err)
		return
	}

	msg := outgoingMsg{
		Type:    "messageDeleted",
		Message: "users added successfully",
		Payload: map[string]any{
			"chatId":    message.ChatID,
			"id":        message.ID,
			"content":   message.Content,
			"updatedAt": message.UpdatedAt,
			"deleted":   true,
		},
	}

	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}

func (h *handler) removeFromGroup(c *gin.Context) {
	currentUserID := context.UserIDFromContext(c)

	chatID, err := parameter.GetParameterValueUint(c, "chatId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	userID, err := parameter.GetParameterValueUint(c, "userId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	// get the chat members before removing the user to avoid not found error, as the user wont be allowed if he is not in the chat
	participants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: "failed to get chat participants",
		})
		log.Printf("error getting chat participants: %v\n", err)
		return
	}

	err = h.csvc.RemoveUserFromGroup(chatID, currentUserID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotPermitted) {
			status = http.StatusForbidden
		}
		c.JSON(status, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response[any]{
		Message: "user removed from group",
	})

	users, err := h.csvc.SearchUsers("", []uint{currentUserID, userID})
	if err != nil {
		log.Printf("error getting current user: %v\n", err)
		return
	}
	if (len(users) != 2 && currentUserID != userID) || (len(users) != 1 && currentUserID == userID) {
		log.Printf("user not found: %d\n", currentUserID)
		return
	}

	var content string
	if currentUserID == userID {
		content = fmt.Sprintf("%s left the group", users[0].Name)
	} else {
		var actorName, targetName string
		for _, user := range users {
			if user.Id == uint32(currentUserID) {
				actorName = user.Name
			}
			if user.Id == uint32(userID) {
				targetName = user.Name
			}
		}
		content = fmt.Sprintf("%s removed %s", actorName, targetName)
	}

	message, err := h.csvc.NewMessage(currentUserID, chatID, content, domain.SystemMessage)
	if err != nil {
		log.Printf("error sending system message: %v\n", err)
		return
	}

	for _, p := range participants {
		h.hub.SendToUser(p.UserID, message)
	}
}

func (h *handler) banFromGroup(c *gin.Context) {
	var req struct {
		UserID    uint       `json:"userId"`
		Reason    string     `json:"reason"`
		ExpiresAt *time.Time `json:"expiresAt"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid request",
		})
		return
	}

	currentUserID := context.UserIDFromContext(c)

	chatID, err := parameter.GetParameterValueUint(c, "chatId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	err = h.csvc.BanUser(chatID, currentUserID, req.UserID, req.Reason, req.ExpiresAt)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotPermitted) {
			status = http.StatusForbidden
		}
		c.JSON(status, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response[any]{
		Message: "user banned from group",
	})

	users, err := h.csvc.SearchUsers("", []uint{currentUserID, req.UserID})
	if err != nil {
		log.Printf("error getting current user: %v\n", err)
		return
	}
	if len(users) != 2 {
		log.Printf("users not found: %d\n", currentUserID)
		return
	}

	var actorName, targetName string
	for _, user := range users {
		if user.Id == uint32(currentUserID) {
			actorName = user.Name
		}
		if user.Id == uint32(req.UserID) {
			targetName = user.Name
		}
	}
	content := fmt.Sprintf("%s banned %s for %s", actorName, targetName, req.Reason)

	message, err := h.csvc.NewMessage(currentUserID, chatID, content, domain.SystemMessage)
	if err != nil {
		log.Printf("error sending system message: %v\n", err)
		return
	}

	participants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		log.Printf("error getting chat participants: %v\n", err)
		return
	}

	for _, p := range participants {
		h.hub.SendToUser(p.UserID, message)
	}
}

func (h *handler) unbanFromGroup(c *gin.Context) {
	currentUserID := context.UserIDFromContext(c)

	chatID, err := parameter.GetParameterValueUint(c, "chatId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	userID, err := parameter.GetParameterValueUint(c, "userId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	err = h.csvc.UnbanUser(chatID, currentUserID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotPermitted) {
			status = http.StatusForbidden
		}
		c.JSON(status, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.Response[any]{
		Message: "user unbanned from group",
	})

	users, err := h.csvc.SearchUsers("", []uint{currentUserID, userID})
	if err != nil {
		log.Printf("error getting current user: %v\n", err)
		return
	}
	if len(users) != 2 {
		log.Printf("users not found: %d\n", currentUserID)
		return
	}

	var actorName, targetName string
	for _, user := range users {
		if user.Id == uint32(currentUserID) {
			actorName = user.Name
		}
		if user.Id == uint32(userID) {
			targetName = user.Name
		}
	}
	content := fmt.Sprintf("%s unbanned %s", actorName, targetName)

	message, err := h.csvc.NewMessage(currentUserID, chatID, content, domain.SystemMessage)
	if err != nil {
		log.Printf("error sending system message: %v\n", err)
		return
	}

	participants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		log.Printf("error getting chat participants: %v\n", err)
		return
	}

	for _, p := range participants {
		h.hub.SendToUser(p.UserID, message)
	}
}

func (h *handler) getBannedUsers(c *gin.Context) {
	chatID, err := parameter.GetParameterValueUint(c, "chatId")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	q := c.Query("q")

	bannedUsers, err := h.csvc.GetBannedUsers(chatID, q)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response[ChatUsers]{
		Data: ChatUsers{
			Users: bannedUsers,
		},
	})
}

func (h *handler) updateUserRole(c *gin.Context) {
	var req struct {
		NewRole string `json:"newRole" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid request",
		})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid user id",
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

	currentUserID := context.UserIDFromContext(c)

	err = h.csvc.UpdateChatUser(chatID, uint(userID), currentUserID, req.NewRole)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrInvalidAction) {
			status = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrNotPermitted) {
			status = http.StatusForbidden
		}
		c.JSON(status, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response[any]{
		Message: "user role updated successfully",
	})

	users, err := h.csvc.SearchUsers("", []uint{currentUserID})
	if err != nil {
		log.Printf("error getting current user: %v\n", err)
		return
	}
	if len(users) == 0 {
		log.Printf("current user not found: %d\n", currentUserID)
		return
	}
	currentUser := users[0]

	updatedUsers, err := h.csvc.SearchUsers("", []uint{uint(userID)})
	if err != nil {
		log.Printf("error getting updated users: %v\n", err)
		return
	}

	names := make([]string, 0, len(updatedUsers))
	for _, u := range updatedUsers {
		names = append(names, u.Name)
	}
	usernames := strings.Join(names, ", ")

	message, err := h.csvc.NewMessage(currentUserID, chatID, fmt.Sprintf("%s made %s %s(s)", currentUser.Name, usernames, req.NewRole), domain.SystemMessage)
	if err != nil {
		log.Printf("error sending system message: %v\n", err)
		return
	}

	participants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		log.Printf("error getting chat participants: %v\n", err)
		return
	}

	msg := outgoingMsg{
		Type:    "updatedUsersRoles",
		Payload: map[string]any{"message": message},
	}
	for _, p := range participants {
		h.hub.SendToUser(p.UserID, msg)
	}
}
