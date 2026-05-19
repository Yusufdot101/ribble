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
	ID          uint          `json:"id"`
	ChatID      uint          `json:"chatId"`
	SenderID    uint          `json:"senderId"`
	Content     string        `json:"content"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	DeletedAt   *time.Time    `json:"deletedAt"`
	Deleted     bool          `json:"deleted"`
	Status      MessageStatus `json:"status"`
	MessageType MessageType   `json:"messageType"`
	ClientID    string        `json:"clientId"`
}

func NewMessage(chatID, senderID uint, content string, messageType MessageType) *Message {
	return &Message{
		ChatID:      chatID,
		SenderID:    senderID,
		Content:     content,
		MessageType: messageType,
	}
}
