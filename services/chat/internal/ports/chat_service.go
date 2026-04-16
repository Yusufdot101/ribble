package ports

import "github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"

type ChatService interface {
	NewChatWithParticipants(userIDs []uint) (*domain.Chat, error)
<<<<<<< Updated upstream
	GetChatParticipants(chatID uint) ([]*domain.ChatParticipant, error)
=======
	GetChatByUserIDs(userIDs []uint) (*domain.Chat, error)
>>>>>>> Stashed changes
	NewMessage(userID, chatID uint, content string) error
}
