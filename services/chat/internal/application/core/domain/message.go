package domain

type Message struct {
	ID       uint
	ChatID   uint
	SenderID uint
	Content  string
}

func NewMessage(chatID, senderID uint, content string) *Message {
	return &Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
	}
}
