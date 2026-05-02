package postgresql

import (
	"context"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"gorm.io/gorm"
)

type ChatParticipant struct {
	gorm.Model
	UserID     uint `gorm:"uniqueIndex:user_chat_idx"`
	ChatID     uint `gorm:"uniqueIndex:user_chat_idx"`
	ChatRoleID uint
}

func (a *Adapter) InsertChatParticipants(chatParticipants []*domain.ChatParticipant) error {
	chatParticipantModels := []*ChatParticipant{}
	for _, chatParticipant := range chatParticipants {
		chatParticipantModels = append(chatParticipantModels, &ChatParticipant{
			UserID:     chatParticipant.UserID,
			ChatID:     chatParticipant.ChatID,
			ChatRoleID: chatParticipant.ChatRoleID,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := a.db.WithContext(ctx).Create(chatParticipantModels)

	return res.Error
}

func (a *Adapter) DeleteChatParticipant(chatID, userID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := a.db.WithContext(ctx).Where("chat_id = ? AND user_id = ?", chatID, userID).Unscoped().Delete(&ChatParticipant{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrRecordNotFound
	}

	return nil
}

func (a *Adapter) GetChatUsers(chatID, currentUserID uint) ([]*domain.ChatParticipant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check user is in the chat
	var count int64
	if err := a.db.WithContext(ctx).
		Model(&ChatParticipant{}).
		Where("chat_id = ? AND user_id = ?", chatID, currentUserID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, domain.ErrRecordNotFound
	}

	// get the chat users
	chatParticipantModels := []*ChatParticipant{}
	if err := a.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Find(&chatParticipantModels).Error; err != nil {
		return nil, err
	}

	chatParticipants := []*domain.ChatParticipant{}
	for _, chatParticipantModel := range chatParticipantModels {
		chatParticipant := &domain.ChatParticipant{
			ID:     chatParticipantModel.ID,
			UserID: chatParticipantModel.UserID,
			ChatID: chatParticipantModel.ChatID,
		}
		chatParticipants = append(chatParticipants, chatParticipant)
	}

	return chatParticipants, nil
}

func (a *Adapter) GetParticipantsByChatIDs(chatIDs []uint) (map[uint][]domain.ChatParticipant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatUsersModels := []*ChatParticipant{}
	err := a.db.WithContext(ctx).
		Joins("JOIN chats ON chats.id = chat_participants.chat_id").
		Where("chats.id IN (?)", chatIDs).
		Find(&chatUsersModels).Error
	if err != nil {
		return nil, err
	}

	chatUsers := make(map[uint][]domain.ChatParticipant)
	for _, chatUserModel := range chatUsersModels {
		chatUsers[chatUserModel.ChatID] = append(chatUsers[chatUserModel.ChatID], domain.ChatParticipant{
			ChatID:     chatUserModel.ChatID,
			ChatRoleID: chatUserModel.ChatRoleID,
			ID:         chatUserModel.ID,
			UserID:     chatUserModel.UserID,
		})
	}

	return chatUsers, nil
}
