package domain

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
