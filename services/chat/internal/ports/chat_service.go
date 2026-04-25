package ports

import "github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"

type ChatService interface {
	NewChatWithParticipants(createChatRequest domain.CreateChatWithParticipantsRequestType) (*domain.Chat, error)
	GetChatParticipants(chatID uint) ([]*domain.ChatParticipant, error)
	GetChatByUserIDs(userIDs []uint) (*domain.Chat, error)
	NewMessage(userID, chatID uint, content string) (*domain.Message, error)
	GetMessages(chatID uint) ([]*domain.Message, error)
	DeleteMessage(userID, messageID uint) (uint, error)
	EditMessage(userID, messageID uint, newContent string) (*domain.Message, error)

	UserHasPermission(userID, chatID uint, permissionName domain.PermissionType) (bool, error)
}
