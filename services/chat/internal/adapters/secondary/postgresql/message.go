package postgresql

import (
	"context"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ChatID   uint
	SenderID uint
	Content  string
}

func (a *Adapter) InsertMessage(message *domain.Message) error {
	messageModel := &Message{
		ChatID:   message.ChatID,
		SenderID: message.SenderID,
		Content:  message.Content,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := a.db.WithContext(ctx).Save(messageModel)
	if res.Error == nil {
		message.ID = messageModel.ID
		message.CreatedAt = messageModel.CreatedAt
	}
	return res.Error
}
