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
	ID        uint `gorm:"primarykey"`
	ChatID    uint
	SenderID  uint
	Content   string
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Status    domain.MessageStatus
}

func (a *Adapter) InsertMessage(message *domain.Message) error {
	messageModel := &Message{
		ChatID:   message.ChatID,
		SenderID: message.SenderID,
		Content:  message.Content,
		Status:   domain.MessageDelivered,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := a.db.WithContext(ctx).Save(messageModel)
	if res.Error == nil {
		message.ID = messageModel.ID
		message.CreatedAt = messageModel.CreatedAt
		message.UpdatedAt = messageModel.UpdatedAt
		message.Status = domain.MessageDelivered
	}
	return res.Error
}

func (a *Adapter) GetMessages(chatID uint, messageFilter domain.GetMessageFilter) ([]*domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var messageModels []*Message
	res := a.db.WithContext(ctx).
		Where("chat_id = ? AND id > ?", chatID, messageFilter.LastMessageID).
		Order("created_at ASC").
		Find(&messageModels)

	if res.Error != nil {
		return nil, res.Error
	}

	messages := []*domain.Message{}
	for _, messageModel := range messageModels {
		status := messageModel.Status
		if status == "" {
			status = domain.MessageDelivered
		}
		message := &domain.Message{
			ID:        messageModel.ID,
			ChatID:    messageModel.ChatID,
			CreatedAt: messageModel.CreatedAt,
			UpdatedAt: messageModel.UpdatedAt,
			SenderID:  messageModel.SenderID,
			Content:   messageModel.Content,
			DeletedAt: messageModel.DeletedAt,
			Deleted:   messageModel.Deleted,
			Status:    status,
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
		Model(&Message{}).
		Where("id = ? AND sender_id = ?", messageID, userID).
		Updates(map[string]any{
			"content":    "",
			"deleted_at": gorm.Expr("Now()"),
			"deleted":    true,
		})
	if res.Error != nil {
		return 0, res.Error
	}

	if res.RowsAffected == 0 {
		return 0, domain.ErrRecordNotFound
	}

	return messageModel.ChatID, nil
}

func (a *Adapter) EditMessage(userID, messageID uint, newContent string) (*domain.Message, error) {
	if strings.TrimSpace(newContent) == "" {
		return nil, domain.ErrInvalidMessageContent
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageModel := &Message{}
	// fetch the message
	err := a.db.WithContext(ctx).Where("id = ? AND sender_id = ?", messageID, userID).First(messageModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	// verify updateablity
	updateWindow := time.Hour
	if time.Since(messageModel.CreatedAt) > updateWindow {
		return nil, domain.ErrUpdateWindowOver
	}

	// update
	err = a.db.WithContext(ctx).
		Model(&Message{}).
		Where("id = ? AND sender_id = ?", messageID, userID).
		Update("content", newContent).Error
	if err != nil {
		return nil, err
	}

	// fetch the updated message
	err = a.db.WithContext(ctx).Where("id = ? AND sender_id = ?", messageID, userID).First(messageModel).Error
	if err != nil {
		return nil, err
	}

	message := &domain.Message{
		ID:        messageModel.ID,
		Content:   messageModel.Content,
		CreatedAt: messageModel.CreatedAt,
		UpdatedAt: messageModel.UpdatedAt,
		ChatID:    messageModel.ChatID,
		SenderID:  messageModel.SenderID,
		DeletedAt: messageModel.DeletedAt,
		Deleted:   messageModel.Deleted,
		Status:    messageModel.Status,
	}

	return message, err
}
