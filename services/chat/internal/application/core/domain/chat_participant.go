package domain

import "time"

type ChatParticipant struct {
	ID         uint
	UserID     uint
	ChatID     uint
	ChatRoleID uint
}

func NewChatParticipant(userID, chatID uint) *ChatParticipant {
	return &ChatParticipant{
		UserID: userID,
		ChatID: chatID,
	}
}

type ChatBan struct {
	ID             uint
	ChatID         uint
	UserID         uint
	BannedByUserID uint
	Reason         string
	ExpiresAt      *time.Time
}

func NewChatBan(chatID, userID, bannedByUserID uint, reason string, expires *time.Time) *ChatBan {
	return &ChatBan{
		ChatID:         chatID,
		UserID:         userID,
		BannedByUserID: bannedByUserID,
		Reason:         reason,
		ExpiresAt:      expires,
	}
}
