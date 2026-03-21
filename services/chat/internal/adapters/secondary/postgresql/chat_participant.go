package postgresql

import (
	"context"
	"time"

	"github.com/Yusufdot101/ribble/services/chat/internal/application/core/domain"
	"gorm.io/gorm"
)

type ChatParticipant struct {
	gorm.Model
	UserID uint
	ChatID uint
}

func (a *Adapter) InsertChatParticipant(chatParticipant *domain.ChatParticipant) error {
	chatParticipantModel := &ChatParticipant{
		UserID: chatParticipant.UserID,
		ChatID: chatParticipant.ChatID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := a.db.WithContext(ctx).Save(chatParticipantModel)
	if res.Error == nil {
		chatParticipant.ID = chatParticipantModel.ID
	}

	return res.Error
}
