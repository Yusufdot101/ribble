package ports

import "github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"

type Repository interface {
	InsertChat(*domain.Chat) error
	InsertChatParticipants([]*domain.ChatParticipant) error
	DeleteChatParticipant(chatID, userID uint) error
	GetChatByUserIDs(userIDs []uint, isGroup bool) (*domain.Chat, error)
	GetChatUsers(chatID, currentUser uint) ([]*domain.ChatParticipant, error)
	GetParticipantsByChatIDs(chatIDs []uint) (map[uint][]domain.ChatParticipant, error)
	GetChatByID(chatID, currentUser uint) (*domain.Chat, error)
	WithTx(fn func(repo Repository) error) error

	GetChatsByUserID(userID uint, query string) ([]*domain.Chat, error)

	InsertMessage(message *domain.Message) error
	GetMessages(chatID uint, messageFilter domain.GetMessageFilter) ([]*domain.Message, error)
	DeleteMessage(chatID, userID, messageID uint) (*domain.Message, error)
	DeleteAnyMessage(chatID, userID uint, username string, messageID uint) (*domain.Message, error)

	EditMessage(userID, messageID uint, newContent string) (*domain.Message, error)

	NewRole(role *domain.Role) error
	NewChatRole(chatRole *domain.ChatRole, roleName domain.RoleType) error
	GrantChatRolePermission(chatRoleID uint, permission domain.PermissionType) error
	GrantUsersChatRoles(userIDs []uint, chatID uint, roleName domain.RoleType) error

	NewPermission(permission *domain.Permission) error
	GetUserPermissions(userID, chatID uint) ([]*domain.Permission, error)

	InsertChatBan(chatBan *domain.ChatBan) error
	GetChatBans(chatID uint) ([]*domain.ChatBan, error)
	DeleteChatBan(chatID, userID uint) error
}
