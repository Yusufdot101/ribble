package ports

import (
	userpb "github.com/Yusufdot101/ripple-proto/golang/user/v4"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
)

type ChatService interface {
	NewChatWithParticipants(createChatRequest domain.CreateChatWithParticipantsRequestType) (*domain.Chat, error)
	GetChatParticipants(chatID, currentUserID uint) ([]*domain.ChatParticipant, error)
	GetParticipantsByChatIDs(chatIDs []uint) (map[uint][]domain.ChatParticipant, error)
	GetChatUsers(chatID, currentUserID uint) ([]*userpb.User, error)
	GetChatsByUserID(userID uint, query string) ([]*domain.Chat, error)
	SearchUsers(query string, ids []uint) ([]*userpb.User, error)
	GetContacts(uint, []uint, string) ([]*userpb.User, error)

	GetChatByUserIDs(userIDs []uint) (*domain.Chat, error)
	GetChatByID(chatID, currentUserID uint) (*domain.Chat, error)
	NewMessage(userID, chatID uint, content string) (*domain.Message, error)
	GetMessages(chatID uint, messageFilter domain.GetMessageFilter) ([]*domain.Message, error)
	DeleteMessage(userID, messageID uint) (uint, error)
	EditMessage(userID, messageID uint, newContent string) (*domain.Message, error)

	UserHasPermission(userID, chatID uint, permissionName domain.PermissionType) (bool, error)
}
