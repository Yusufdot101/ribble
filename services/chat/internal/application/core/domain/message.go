package domain

import "time"

type GetMessageFilter struct {
	LastMessageID uint
}

type (
	MessageStatus string
	MessageType   string // for groups, to show when a user is added, etc
)

const (
	MessageDelivered MessageStatus = "delivered"
	MessageFailed    MessageStatus = "failed"
	StandardMessage  MessageType   = "standard"
	SystemMessage    MessageType   = "information message"
)

type Message struct {
	ID          uint
	ChatID      uint
	SenderID    uint
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Deleted     bool
	Status      MessageStatus
	MessageType MessageType
	ClientID    string
}

func NewMessage(chatID, senderID uint, content string, messageType MessageType) *Message {
	return &Message{
		ChatID:      chatID,
		SenderID:    senderID,
		Content:     content,
		MessageType: messageType,
	}
}
