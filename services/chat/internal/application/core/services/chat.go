package services

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	userpb "github.com/Yusufdot101/ripple-proto/golang/user/v4"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/chat/internal/ports"
)

type ChatService struct {
	repo        ports.Repository
	userService ports.UserService
}

func NewChatService(repo ports.Repository, userService ports.UserService) *ChatService {
	return &ChatService{
		repo:        repo,
		userService: userService,
	}
}

func (csvc *ChatService) NewChatWithParticipants(createChatRequest domain.CreateChatWithParticipantsRequestType) (*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userIDs := slices.Collect(maps.Keys(createChatRequest.UserRoles))

	// verify the users actually exist
	valid, err := csvc.userService.VerifyUsers(ctx, userIDs)
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
		chat = domain.NewChat(createChatRequest.Name, createChatRequest.IsGroup)
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
				case "send message":
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

func (csvc *ChatService) GetChatsByUserID(userID uint, query string) ([]*domain.Chat, error) {
	return csvc.repo.GetChatsByUserID(userID, query)
}

func (csvc *ChatService) GetChatByID(chatID, currentUserID uint) (*domain.Chat, error) {
	return csvc.repo.GetChatByID(chatID, currentUserID)
}

func (csvc *ChatService) GetChatParticipants(chatID, currentUserID uint) ([]*domain.ChatParticipant, error) {
	return csvc.repo.GetChatUsers(chatID, currentUserID)
}

func (csvc *ChatService) GetParticipantsByChatIDs(chatIDs []uint) (map[uint][]domain.ChatParticipant, error) {
	return csvc.repo.GetParticipantsByChatIDs(chatIDs)
}

func (csvc *ChatService) GetChatUsers(chatID, currentUserID uint) ([]*userpb.User, error) {
	chatParticipants, err := csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		return nil, err
	}
	userIDs := []uint{}
	for _, user := range chatParticipants {
		userIDs = append(userIDs, user.UserID)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcUsers, err := csvc.userService.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return grpcUsers, nil
}

func (csvc *ChatService) SearchUsers(query string, ids []uint) ([]*userpb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcUsers, err := csvc.userService.SearchUsers(ctx, query, ids)
	if err != nil {
		return nil, err
	}

	return grpcUsers, nil
}

func (csvc *ChatService) GetContacts(currentUserID uint, excludeIDs []uint, query string) ([]*userpb.User, error) {
	ctx := context.Background()
	grpcUsers, err := csvc.userService.GetContacts(ctx, query, excludeIDs, currentUserID)
	if err != nil {
		return nil, err
	}

	return grpcUsers, nil
}

func (csvc *ChatService) NewMessage(userID, chatID uint, content string) (*domain.Message, error) {
	message := domain.NewMessage(chatID, userID, content)
	return message, csvc.repo.InsertMessage(message)
}

func (csvc *ChatService) GetMessages(chatID uint, messageFilter domain.GetMessageFilter) ([]*domain.Message, error) {
	return csvc.repo.GetMessages(chatID, messageFilter)
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
