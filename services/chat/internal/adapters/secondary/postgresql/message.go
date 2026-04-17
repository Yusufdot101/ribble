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

func (a *Adapter) GetMessages(chatID uint) ([]*domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var messageModels []*Message
	res := a.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Find(&messageModels)

	if res.Error != nil {
		return nil, res.Error
	}

	messages := []*domain.Message{}
	for _, messageModel := range messageModels {
		message := &domain.Message{
			ID:        messageModel.ID,
			ChatID:    messageModel.ChatID,
			CreatedAt: messageModel.CreatedAt,
			SenderID:  messageModel.SenderID,
			Content:   messageModel.Content,
		}
		messages = append(messages, message)
	}

	return messages, nil
}
