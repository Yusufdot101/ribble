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

func (csvc *ChatService) NewChatWithParticipants(userIDs []uint) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// verify the users actually exist
	valid, err := csvc.userVerifier.VerifyUsers(ctx, userIDs)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, domain.ErrInvalidUserIDs
	}

	chatID := uint(0)
	err = csvc.repo.WithTx(func(repo ports.Repository) error {
		chat := domain.NewChat()
		err := repo.InsertChat(chat)
		if err != nil {
			return err
		}
		chatID = chat.ID

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
		return 0, err
	}

	return chatID, nil
}
