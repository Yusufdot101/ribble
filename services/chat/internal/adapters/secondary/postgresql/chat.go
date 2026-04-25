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
	Participants []ChatParticipant `gorm:"constraint:OnDelete:CASCADE;"`
	Messages     []Message         `gorm:"constraint:OnDelete:CASCADE;"`
	ChatRoles    []ChatRole        `gorm:"constraint:OnDelete:CASCADE;"`
}

func (a *Adapter) InsertChat(chat *domain.Chat) error {
	chatModel := &Chat{}

	res := a.db.Save(chatModel)
	if res.Error == nil {
		chat.ID = chatModel.ID
	}
	return res.Error
}

func (a *Adapter) GetChatByParticipantIDs(participantIDs []uint) (*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var chatModel Chat
	db := a.db.WithContext(ctx)
	subQuery := db.
		Table("chat_participants").
		Select("chat_id").
		Where("user_id IN ?", participantIDs).
		Group("chat_id").
		Having("Count(DISTINCT user_id) = ?", len(participantIDs))

	res := db.WithContext(ctx).Where("id IN (?)", subQuery).First(&chatModel)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, res.Error
	}

	chat := &domain.Chat{}
	chat.ID = chatModel.ID

	return chat, nil
}

func (a *Adapter) WithTx(fn func(repo ports.Repository) error) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &Adapter{db: tx} // same struct, but with tx as the db
		return fn(txRepo)
	})
}
