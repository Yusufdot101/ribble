package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	message, err := h.csvc.NewMessage(currentUserID, chatIDUint, fmt.Sprintf("%s added %s", currentUser.Name, usernames), domain.SystemMessage)
	if err != nil {
		log.Printf("error sending system message: %v\n", err)
		return
	}

	participants, err := h.csvc.GetChatParticipants(chatIDUint, currentUserID)
	if err != nil {
		log.Printf("error getting chat participants: %v\n", err)
		return
	}

	for _, p := range participants {
		h.hub.SendToUser(p.UserID, message)
	}
}
