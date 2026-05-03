package services

import (
	"context"
	"errors"
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
				case "add users to group":
					permission = domain.AddToGroup
				case "remove users from group":
					permission = domain.RemoveUserFromGroup
				case "delete messages":
					permission = domain.DeleteMessages
				case "ban users":
					permission = domain.BanUsers
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
		participants := []*domain.ChatParticipant{}
		for _, userID := range userIDs {
			participant := domain.NewChatParticipant(uint(userID), chat.ID)
			participants = append(participants, participant)
		}
		err = repo.InsertChatParticipants(participants)
		if err != nil {
			return err
		}

		memberParticipants := []uint{}
		adminParticipants := []uint{}
		for _, userID := range userIDs {
			// 5. grant role to users
			roleName := createChatRequest.UserRoles[userID]
			switch roleName {
			case "member":
				memberParticipants = append(memberParticipants, userID)
			case "admin":
				adminParticipants = append(adminParticipants, userID)
			default:
				return fmt.Errorf("%w: %s", domain.ErrInvalidRole, roleName)
			}
		}
		err = repo.GrantUsersChatRoles(memberParticipants, chat.ID, domain.Member)
		if err != nil {
			return err
		}

		err = repo.GrantUsersChatRoles(adminParticipants, chat.ID, domain.Admin)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (csvc *ChatService) GetChatByUserIDs(userIDs []uint, isGroup bool) (*domain.Chat, error) {
	return csvc.repo.GetChatByUserIDs(userIDs, isGroup)
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

func (csvc *ChatService) NewMessage(userID, chatID uint, content string, messageType domain.MessageType) (*domain.Message, error) {
	message := domain.NewMessage(chatID, userID, content, messageType)
	return message, csvc.repo.InsertMessage(message)
}

func (csvc *ChatService) GetMessages(chatID uint, messageFilter domain.GetMessageFilter) ([]*domain.Message, error) {
	return csvc.repo.GetMessages(chatID, messageFilter)
}

func (csvc *ChatService) DeleteMessage(chatID, userID, messageID uint) (*domain.Message, error) {
	hasPermission, err := csvc.UserHasPermission(userID, chatID, domain.DeleteMessages)
	if err != nil {
		return nil, err
	}

	var message *domain.Message
	if hasPermission {
		users, err := csvc.SearchUsers("", []uint{userID})
		if err != nil {
			return nil, err
		}
		if len(users) != 1 {
			return nil, errors.New("error getting the user")
		}
		user := users[0]
		message, err = csvc.repo.DeleteAnyMessage(chatID, userID, user.Name, messageID)
		if err != nil {
			return nil, err
		}
	} else {
		message, err = csvc.repo.DeleteMessage(chatID, userID, messageID)
	}
	if err != nil {
		return nil, err
	}
	return message, err
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

func (csvc *ChatService) GetUserPermissions(chatID, userID uint) ([]*domain.Permission, error) {
	return csvc.repo.GetUserPermissions(userID, chatID)
}

func (csvc *ChatService) AddUsersToGroup(chatID, currentUserID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return domain.ErrInvalidUserIDs
	}

	userHasPermission, err := csvc.UserHasPermission(currentUserID, chatID, domain.AddToGroup)
	if err != nil {
		return domain.ErrNotPermitted
	}

	if !userHasPermission {
		return domain.ErrNotPermitted
	}

	chatBans, err := csvc.repo.GetChatBans(chatID)
	if err != nil {
		return err
	}

	excludedUsersID := []uint{}
	for _, ban := range chatBans {
		excludedUsersID = append(excludedUsersID, ban.UserID)
	}

	chatParticipants := []*domain.ChatParticipant{}
	for _, userID := range userIDs {
		if slices.Contains(excludedUsersID, userID) {
			return domain.ErrBannedUser
		}
		chatParticipants = append(chatParticipants, domain.NewChatParticipant(userID, chatID))
	}
	return csvc.repo.WithTx(func(repo ports.Repository) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		isValid, err := csvc.userService.VerifyUsers(ctx, userIDs)
		if err != nil {
			return err
		}

		if !isValid {
			return domain.ErrInvalidUserIDs
		}

		err = repo.InsertChatParticipants(chatParticipants)
		if err != nil {
			return err
		}
		return repo.GrantUsersChatRoles(userIDs, chatID, domain.Member)
	})
}

func (csvc *ChatService) RemoveUserFromGroup(chatID, currentUserID, userID uint) error {
	if currentUserID != userID {
		userHasPermission, err := csvc.UserHasPermission(currentUserID, chatID, domain.RemoveUserFromGroup)
		if err != nil {
			return domain.ErrNotPermitted
		}

		if !userHasPermission {
			return domain.ErrNotPermitted
		}
	}

	return csvc.repo.DeleteChatParticipant(chatID, userID)
}

func (csvc *ChatService) BanUser(
	chatID, currentUserID, userID uint, reason string, expiresAt *time.Time,
) error {
	userHasPermission, err := csvc.UserHasPermission(currentUserID, chatID, domain.BanUsers)
	if err != nil {
		return domain.ErrNotPermitted
	}

	if !userHasPermission {
		return domain.ErrNotPermitted
	}

	err = csvc.repo.WithTx(func(repo ports.Repository) error {
		err := repo.DeleteChatParticipant(chatID, userID)
		if err != nil {
			return err
		}

		chatBan := domain.NewChatBan(chatID, userID, currentUserID, reason, expiresAt)
		err = repo.InsertChatBan(chatBan)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (csvc *ChatService) GetAddableChatUsers(chatID, currentUserID uint, query string) ([]*userpb.User, error) {
	ctx := context.Background()
	chatBans, err := csvc.repo.GetChatBans(chatID)
	if err != nil {
		return nil, err
	}

	chatUsers, err := csvc.repo.GetChatUsers(chatID, currentUserID)
	if err != nil {
		return nil, err
	}

	excludedUsersID := []uint{}
	for _, ban := range chatBans {
		excludedUsersID = append(excludedUsersID, ban.UserID)
	}
	for _, user := range chatUsers {
		excludedUsersID = append(excludedUsersID, user.UserID)
	}

	grpcUsers, err := csvc.userService.GetContacts(ctx, query, excludedUsersID, currentUserID)
	if err != nil {
		return nil, err
	}

	return grpcUsers, nil
}

func (csvc *ChatService) UnbanUser(chatID, currentUserID, userID uint) error {
	userHasPermission, err := csvc.UserHasPermission(currentUserID, chatID, domain.BanUsers)
	if err != nil {
		return domain.ErrNotPermitted
	}

	if !userHasPermission {
		return domain.ErrNotPermitted
	}

	return csvc.repo.DeleteChatBan(chatID, userID)
}
