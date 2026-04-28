package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/chat/internal/ports"
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Name         string
	IsGroup      bool
	Participants []ChatParticipant `gorm:"constraint:OnDelete:CASCADE;"`
	Messages     []Message         `gorm:"constraint:OnDelete:CASCADE;"`
	ChatRoles    []ChatRole        `gorm:"constraint:OnDelete:CASCADE;"`
}

func (a *Adapter) InsertChat(chat *domain.Chat) error {
	chatModel := &Chat{
		Name:    chat.Name,
		IsGroup: chat.IsGroup,
	}

	res := a.db.Save(chatModel)
	if res.Error == nil {
		chat.ID = chatModel.ID
	}
	return res.Error
}

func (a *Adapter) GetChatByID(chatID, currentUserID uint) (*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatModel := &Chat{}
	err := a.db.WithContext(ctx).
		Model(&Chat{}).
		Where("id = ?", chatID).
		Where("EXISTS (?)",
			a.db.Table("chat_participants").
				Select("1").
				Where("chat_participants.chat_id = chats.id").
				Where("chat_participants.user_id = ?", currentUserID),
		).
		First(chatModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	chat := &domain.Chat{
		ID:      chatModel.ID,
		Name:    chatModel.Name,
		IsGroup: chatModel.IsGroup,
	}

	return chat, nil
}

func (a *Adapter) GetChatsByUserID(userID uint, query string) ([]*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatModels := []*Chat{}
	err := a.db.WithContext(ctx).
		Joins("JOIN chat_participants ON chat_participants.chat_id = chats.id").
		Where("chat_participants.user_id = ? AND chats.name ILIKE ?", userID, "%"+query+"%").
		Find(&chatModels).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	chats := []*domain.Chat{}
	for _, chatModel := range chatModels {
		chats = append(chats, &domain.Chat{
			ID:      chatModel.ID,
			Name:    chatModel.Name,
			IsGroup: chatModel.IsGroup,
		})
	}

	return chats, nil
}

func (a *Adapter) GetChatByParticipantIDs(participantIDs []uint) (*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatModel := Chat{}
	participantCount := len(participantIDs)

	// Subquery to find chat_ids that have the exact participants
	subQuery := a.db.Table("chat_participants").
		Select("chat_id").
		Group("chat_id").
		// 1. Ensure the total number of members in the chat matches your list size
		Having("COUNT(*) = ?", participantCount).
		// 2. Ensure the number of members matching your list also matches your list size
		Having("SUM(CASE WHEN user_id IN (?) THEN 1 ELSE 0 END) = ?", participantIDs, participantCount)

	res := a.db.WithContext(ctx).
		Where("id IN (?)", subQuery).
		First(&chatModel)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, res.Error
	}

	chat := &domain.Chat{
		ID:      chatModel.ID,
		Name:    chatModel.Name,
		IsGroup: chatModel.IsGroup,
	}

	return chat, nil
}

func (a *Adapter) WithTx(fn func(repo ports.Repository) error) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &Adapter{db: tx} // same struct, but with tx as the db
		return fn(txRepo)
	})
}
