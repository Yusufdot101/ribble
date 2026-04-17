package ports

import "github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"

type Repository interface {
	InsertChat(*domain.Chat) error
	InsertChatParticipant(*domain.ChatParticipant) error
	GetChatByParticipantIDs(participantIDs []uint) (*domain.Chat, error)
	GetChatUsers(chatID uint) ([]*domain.ChatParticipant, error)
	WithTx(fn func(repo Repository) error) error
	InsertMessage(message *domain.Message) error
}
