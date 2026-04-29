package domain

import "time"

type GetMessageFilter struct {
	LastMessageID uint
}

type MessageStatus string

var (
	MessageDelivered MessageStatus = "delivered"
	MessageFailed    MessageStatus = "failed"
)

type Message struct {
	ID        uint
	ChatID    uint
	SenderID  uint
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Deleted   bool
	Status    MessageStatus
}

func NewMessage(chatID, senderID uint, content string) *Message {
	return &Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
	}
}
