package domain

import "time"

type Message struct {
	ID        uint
	ChatID    uint
	SenderID  uint
	Content   string
	CreatedAt time.Time
}

func NewMessage(chatID, senderID uint, content string) *Message {
	return &Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
	}
}
