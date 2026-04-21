package postgresql

import (
	"context"
	"errors"
	"strings"
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

func (a *Adapter) DeleteMessage(userID, messageID uint) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var messageModel Message

	err := a.db.WithContext(ctx).
		Select("chat_id").
		Where("id = ? AND sender_id = ?", messageID, userID).
		First(&messageModel).Error
	if err != nil {
		return 0, err
	}

	res := a.db.WithContext(ctx).
		Where("id = ? AND sender_id = ?", messageID, userID).
		Delete(&Message{})
	if res.Error != nil {
		return 0, err
	}

	if res.RowsAffected == 0 {
		return 0, domain.ErrRecordNotFound
	}

	return messageModel.ChatID, nil
}

func (a *Adapter) EditMessage(userID, messageID uint, newContent string) (uint, error) {
	if strings.TrimSpace(newContent) == "" {
		return 0, domain.ErrInvalidMessageContent
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageModel := &Message{}
	// fetch the message
	err := a.db.WithContext(ctx).Where("id = ? AND sender_id = ?", messageID, userID).First(messageModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, domain.ErrRecordNotFound
		}
		return 0, err
	}

	// verify updateablity
	updateWindow := time.Hour
	if time.Since(messageModel.CreatedAt) > updateWindow {
		return 0, domain.ErrUpdateWindowOver
	}

	// update
	err = a.db.WithContext(ctx).
		Model(&Message{}).
		Where("id = ? AND sender_id = ?", messageID, userID).
		Update("content", newContent).Error
	if err != nil {
		return 0, err
	}

	return messageModel.ChatID, err
}
