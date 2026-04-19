package services

import (
	"context"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/chat/internal/ports"
)

type ChatService struct {
	repo         ports.Repository
	userVerifier ports.UserVerifier
}

func NewChatService(repo ports.Repository, userVerifier ports.UserVerifier) *ChatService {
	return &ChatService{
		repo:         repo,
		userVerifier: userVerifier,
	}
}

func (csvc *ChatService) NewChatWithParticipants(userIDs []uint) (*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// verify the users actually exist
	valid, err := csvc.userVerifier.VerifyUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, domain.ErrInvalidUserIDs
	}

	chat := &domain.Chat{}
	err = csvc.repo.WithTx(func(repo ports.Repository) error {
		chat = domain.NewChat()
		err := repo.InsertChat(chat)
		if err != nil {
			return err
		}

		for _, userID := range userIDs {
			participant := domain.NewChatParticipant(uint(userID), chat.ID)
			err = repo.InsertChatParticipant(participant)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (csvc *ChatService) GetChatByUserIDs(userIDs []uint) (*domain.Chat, error) {
	return csvc.repo.GetChatByParticipantIDs(userIDs)
}

func (csvc *ChatService) GetChatParticipants(chatID uint) ([]*domain.ChatParticipant, error) {
	return csvc.repo.GetChatUsers(chatID)
}

func (csvc *ChatService) NewMessage(userID, chatID uint, content string) (*domain.Message, error) {
	message := domain.NewMessage(chatID, userID, content)
	return message, csvc.repo.InsertMessage(message)
}

func (csvc *ChatService) GetMessages(chatID uint) ([]*domain.Message, error) {
	return csvc.repo.GetMessages(chatID)
}

func (csvc *ChatService) DeleteMessage(userID, messageID uint) (uint, error) {
	return csvc.repo.DeleteMessage(userID, messageID)
}
