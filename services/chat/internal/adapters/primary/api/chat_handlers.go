package api

import (
	"errors"
	"maps"
	"net/http"
	"slices"
	"strconv"

	userpb "github.com/Yusufdot101/ripple-proto/golang/user/v4"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/parameter"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/response"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

type Chat struct {
	Name    string `json:"name"`
	IsGroup bool   `json:"isGroup"`
	ID      uint   `json:"id"`
}

func (h *handler) GetOrCreateChat(ctx *gin.Context) {
	var createChatRequest domain.CreateChatWithParticipantsRequestType
	if err := ctx.ShouldBind(&createChatRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: err.Error(),
		})
		return
	}
	currentUserID := context.UserIDFromContext(ctx)
	if createChatRequest.UserRoles == nil {
		createChatRequest.UserRoles = make(map[uint]string)
	}

	createChatRequest.UserRoles[currentUserID] = "creator"
	if len(createChatRequest.UserRoles) < 2 {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "at least 2 participants required",
		})
		return
	}

	if createChatRequest.RolePermissions == nil {
		createChatRequest.RolePermissions = make(map[string][]string)
	}
	adminPerms, ok := createChatRequest.RolePermissions["admin"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "admin role permissions are required",
		})
		return
	}
	createChatRequest.RolePermissions["creator"] = append([]string(nil), adminPerms...)
	createChatRequest.RolePermissions["creator"] = append(
		createChatRequest.RolePermissions["creator"],
		string(domain.PromoteMembers), string(domain.DemoteAdmins),
		string(domain.UpdatePermissions),
	)

	userIDs := slices.Collect(maps.Keys(createChatRequest.UserRoles))
	var chat *domain.Chat
	var err error
	if !createChatRequest.IsGroup {
		chat, err = h.csvc.GetChatByUserIDs(userIDs, createChatRequest.IsGroup)
		if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, response.Response[any]{
				Error: err.Error(),
			})
			return
		}
	}

	// create chat if not exists
	if errors.Is(err, domain.ErrRecordNotFound) || createChatRequest.IsGroup {
		chat, err = h.csvc.NewChatWithParticipants(createChatRequest)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.Response[any]{
				Error: err.Error(),
			})
			return
		}
		payload := map[string]any{
			"chat":    chat,
			"userIds": userIDs,
		}
		msg := outgoingMsg{
			Type:    "message",
			SubType: "chatCreated",
			Payload: payload,
		}

		for _, user := range userIDs {
			h.hub.SendToUser(user, msg)
		}
	}

	ctx.JSON(http.StatusOK, response.Response[Chat]{
		Data: Chat{
			Name:    chat.Name,
			ID:      chat.ID,
			IsGroup: chat.IsGroup,
		},
	})
}

func (h *handler) getChatByID(ctx *gin.Context) {
	currentUserID := context.UserIDFromContext(ctx)

	chatID, err := strconv.ParseUint(ctx.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}
	if chatID > uint64(^uint(0)) {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}
	chatIDUint := uint(chatID)

	chat, err := h.csvc.GetChatByID(chatIDUint, currentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response[Chat]{
		Data: Chat{
			Name:    chat.Name,
			ID:      chat.ID,
			IsGroup: chat.IsGroup,
		},
	})
}

type ChatUser struct {
	*userpb.User
	Role string `json:"role"`
}

type ChatMemberResponse struct {
	Users []*ChatUser `json:"users"`
}

func (h *handler) getChatUsers(ctx *gin.Context) {
	currentUserID := context.UserIDFromContext(ctx)

	chatID, err := parameter.GetParameterValueUint(ctx, "chatId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}

	// get participants
	chatParticipants, err := h.csvc.GetChatParticipants(chatID, currentUserID)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			ctx.JSON(http.StatusForbidden, response.Response[any]{
				Error: "not a participant of this chat",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	// get users
	chatUsers, err := h.csvc.GetChatUsers(chatID, currentUserID, chatParticipants)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			ctx.JSON(http.StatusForbidden, response.Response[any]{
				Error: "not a participant of this chat",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	// get roles
	chatUsersRoles, err := h.csvc.GetChatUsersRoles(chatID, currentUserID, chatParticipants)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			ctx.JSON(http.StatusForbidden, response.Response[any]{
				Error: "not a participant of this chat",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}
	result := []*ChatUser{}
	for _, user := range chatUsers {
		role, exists := chatUsersRoles[uint(user.Id)]
		if !exists {
			ctx.JSON(http.StatusInternalServerError, response.Response[any]{
				Error: "an error occurred, please try again later",
			})
			return
		}
		result = append(result, &ChatUser{
			User: user,
			Role: role,
		})
	}

	ctx.JSON(http.StatusOK, response.Response[ChatMemberResponse]{
		Data: ChatMemberResponse{
			Users: result,
		},
	})
}

type AddableChatUsersResponse struct {
	Users []*userpb.User `json:"users"`
}

func (h *handler) getAddableChatUsers(c *gin.Context) {
	currentUserID := context.UserIDFromContext(c)
	q := c.Query("q")

	chatID, err := strconv.ParseUint(c.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}
	if chatID > uint64(^uint(0)) {
		c.JSON(http.StatusBadRequest, response.Response[any]{
			Error: "invalid chat id",
		})
		return
	}
	chatIDUint := uint(chatID)

	addableUsers, err := h.csvc.GetAddableChatUsers(chatIDUint, currentUserID, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response[any]{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response[AddableChatUsersResponse]{
		Data: AddableChatUsersResponse{
			Users: addableUsers,
		},
	})
}
