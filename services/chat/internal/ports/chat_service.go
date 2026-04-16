package ports

import "github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"

type ChatService interface {
	NewChatWithParticipants(userIDs []uint) (*domain.Chat, error)
	GetChatParticipants(chatID uint) ([]*domain.ChatParticipant, error)
	NewMessage(userID, chatID uint, content string) error
}
