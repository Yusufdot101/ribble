package api

import (
	"errors"
	"maps"
	"net/http"
	"slices"
	"strconv"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

func (h *handler) GetOrCreateChat(ctx *gin.Context) {
	var createChatRequest domain.CreateChatWithParticipantsRequestType
	if err := ctx.ShouldBind(&createChatRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	currentUserID := context.UserIDFromContext(ctx)
	if createChatRequest.UserRoles == nil {
		createChatRequest.UserRoles = make(map[uint]string)
	}

	createChatRequest.UserRoles[currentUserID] = "admin"
	if len(createChatRequest.UserRoles) < 2 {
		ctx.String(http.StatusBadRequest, "at least 2 participants required")
		return
	}

	userIDs := slices.Collect(maps.Keys(createChatRequest.UserRoles))
	var chat *domain.Chat
	var err error
	if !createChatRequest.IsGroup {
		chat, err = h.csvc.GetChatByUserIDs(userIDs)
		if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// create chat if not exists
	if errors.Is(err, domain.ErrRecordNotFound) || createChatRequest.IsGroup {
		chat, err = h.csvc.NewChatWithParticipants(createChatRequest)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"chat": chat,
	})
}

func (h *handler) getChatByID(ctx *gin.Context) {
	currentUserID := context.UserIDFromContext(ctx)

	chatID, err := strconv.ParseUint(ctx.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}
	if chatID > uint64(^uint(0)) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}
	chatIDUint := uint(chatID)

	chat, err := h.csvc.GetChatByID(chatIDUint, currentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"chat": chat,
	})
}

func (h *handler) getChatUsers(ctx *gin.Context) {
	currentUserID := context.UserIDFromContext(ctx)

	chatID, err := strconv.ParseUint(ctx.Param("chatId"), 10, strconv.IntSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}
	if chatID > uint64(^uint(0)) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}
	chatIDUint := uint(chatID)

	chatUsers, err := h.csvc.GetChatUsers(chatIDUint, currentUserID)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "not a participant of this chat",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"chatUsers": chatUsers,
	})
}
