package services

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
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

func (csvc *ChatService) NewChatWithParticipants(createChatRequest domain.CreateChatWithParticipantsRequestType) (*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userIDs := slices.Collect(maps.Keys(createChatRequest.UserRoles))

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
		// steps:
		// 1. create chat
		chat = domain.NewChat()
		err := repo.InsertChat(chat)
		if err != nil {
			return err
		}

		// 2. create chat roles
		for roleName, permissions := range createChatRequest.RolePermissions {
			var chatRole *domain.ChatRole
			switch roleName {
			case "member":
				chatRole = domain.NewChatRole(chat.ID)
				err = repo.NewChatRole(chatRole, domain.Member)
			case "admin":
				chatRole = domain.NewChatRole(chat.ID)
				err = repo.NewChatRole(chatRole, domain.Admin)
			default:
				return fmt.Errorf("%w: %s", domain.ErrInvalidRole, roleName)
			}

			if err != nil {
				return err
			}

			// 3. grant permissions to roles
			var permission domain.PermissionType
			for _, permissionName := range permissions {
				switch permissionName {
				case "sendMessages":
					permission = domain.SendMessage
				default:
					return fmt.Errorf("%w: %s", domain.ErrInvalidPermission, permissionName)
				}

				err = repo.GrantChatRolePermission(chatRole.ID, permission)
				if err != nil {
					return err
				}
			}
		}

		// 4. create chat_users
		for _, userID := range userIDs {
			participant := domain.NewChatParticipant(uint(userID), chat.ID)
			err = repo.InsertChatParticipant(participant)
			if err != nil {
				return err
			}

			// 5. grant role to users
			var userRole domain.RoleType
			roleName := createChatRequest.UserRoles[userID]
			switch roleName {
			case "member":
				userRole = domain.Member
			case "admin":
				userRole = domain.Admin
			default:
				return fmt.Errorf("%w: %s", domain.ErrInvalidRole, roleName)
			}
			err = repo.GrantUserChatRole(participant.UserID, participant.ChatID, userRole)
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

func (csvc *ChatService) EditMessage(userID, messageID uint, newContent string) (*domain.Message, error) {
	if strings.TrimSpace(newContent) == "" {
		return nil, domain.ErrInvalidMessageContent
	}
	return csvc.repo.EditMessage(userID, messageID, newContent)
}

func (csvc *ChatService) UserHasPermission(userID, chatID uint, permissionName domain.PermissionType) (bool, error) {
	userPermissions, err := csvc.repo.GetUserPermissions(userID, chatID)
	if err != nil {
		return false, err
	}

	p := domain.NewPermission(permissionName)
	return p.IncludedIn(userPermissions), nil
}
