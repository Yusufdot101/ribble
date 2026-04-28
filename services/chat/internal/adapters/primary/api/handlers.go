package api

import (
	"net/http"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

type DirectChat struct {
	Chat        *domain.Chat
	OtherUserID uint
}

func (h *handler) getConversations(c *gin.Context) {
	currentUserID := context.UserIDFromContext(c)
	q := c.Query("q")

	// 1. get all chats
	chats, err := h.csvc.GetChatsByUserID(currentUserID, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	chatIDs := make([]uint, 0, len(chats))
	for _, ch := range chats {
		chatIDs = append(chatIDs, ch.ID)
	}

	// 2. batch load participants (NO N+1)
	participantsByChat, err := h.csvc.GetParticipantsByChatIDs(chatIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var groups []*domain.Chat
	var directChats []DirectChat
	var otherUserIDs []uint

	// 3. classify in memory
	for _, chat := range chats {
		members := participantsByChat[chat.ID]

		if chat.IsGroup {
			groups = append(groups, chat)
			continue
		}

		for _, m := range members {
			if m.UserID != currentUserID {
				directChats = append(directChats, DirectChat{
					Chat:        chat,
					OtherUserID: m.UserID,
				})
				otherUserIDs = append(otherUserIDs, m.UserID)
			}
		}
	}

	// 4. search matched users
	userMatch := make(map[uint32]bool)
	if len(otherUserIDs) > 0 {
		matchedUsers, err := h.csvc.SearchUsers(q, otherUserIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, u := range matchedUsers {
			userMatch[u.Id] = true
		}
	}

	// 5. filter chats safely
	var filteredChats []*domain.Chat
	for _, dc := range directChats {
		if userMatch[uint32(dc.OtherUserID)] {
			filteredChats = append(filteredChats, dc.Chat)
		}
	}

	// 6. contacts (unchanged logic)
	contacts, err := h.csvc.GetContacts(currentUserID, otherUserIDs, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Groups":   groups,
		"Chats":    filteredChats,
		"Contacts": contacts,
	})
}
