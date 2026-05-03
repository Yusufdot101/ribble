package postgresql

import (
	"context"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"gorm.io/gorm"
)

type ChatBan struct {
	gorm.Model
	UserID         uint `gorm:"uniqueIndex:user_chat_idx"`
	ChatID         uint `gorm:"uniqueIndex:user_chat_idx"`
	BannedByUserID uint
	Reason         string
	ExpiresAt      *time.Time
}

func (a *Adapter) InsertChatBan(chatBan *domain.ChatBan) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatBanModel := &ChatBan{
		ChatID:         chatBan.ChatID,
		UserID:         chatBan.UserID,
		BannedByUserID: chatBan.BannedByUserID,
		Reason:         chatBan.Reason,
		ExpiresAt:      chatBan.ExpiresAt,
	}
	res := a.db.WithContext(ctx).Save(chatBanModel)
	if res.Error == nil {
		chatBan.ID = chatBanModel.ID
	}

	return res.Error
}

func (a *Adapter) GetChatBans(chatID uint) ([]*domain.ChatBan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatBanModels := []*domain.ChatBan{}
	err := a.db.WithContext(ctx).Where("expires_at > ?", time.Now()).Find(&chatBanModels).Error
	if err != nil {
		return nil, err
	}

	chatBans := []*domain.ChatBan{}
	for _, ban := range chatBanModels {
		chatBans = append(chatBans, &domain.ChatBan{
			ID:             ban.ID,
			ChatID:         ban.ChatID,
			UserID:         ban.UserID,
			BannedByUserID: ban.BannedByUserID,
			Reason:         ban.Reason,
			ExpiresAt:      ban.ExpiresAt,
		})
	}

	return chatBans, nil
}
