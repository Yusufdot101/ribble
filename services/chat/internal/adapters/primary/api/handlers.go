package api

import (
	"errors"
	"maps"
	"net/http"
	"slices"

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
	chat, err := h.csvc.GetChatByUserIDs(userIDs)
	if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// create chat if not exists
	if errors.Is(err, domain.ErrRecordNotFound) {
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
