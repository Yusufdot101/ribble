package api

import (
	"errors"
	"net/http"

	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api/context"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/gin-gonic/gin"
)

func (h *handler) NewChatWithParticipants(ctx *gin.Context) {
	var createChatWithParticipantsRequests struct {
		UserIDs []uint `json:"userIDs"`
	}
	if err := ctx.ShouldBind(&createChatWithParticipantsRequests); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	currentUserID := context.UserIDFromContext(ctx)
	createChatWithParticipantsRequests.UserIDs = append(createChatWithParticipantsRequests.UserIDs, currentUserID)
	if len(createChatWithParticipantsRequests.UserIDs) < 2 {
		ctx.String(http.StatusBadRequest, "userIDs cannot be less than 2")
		return
	}
	chat, err := h.csvc.NewChatWithParticipants(createChatWithParticipantsRequests.UserIDs)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"chat": chat,
	})
}

func (h *handler) GetByUserIDs(ctx *gin.Context) {
	var GetChatRequest struct {
		UserIDs []uint `json:"userIDs"`
	}

	if err := ctx.ShouldBind(&GetChatRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	currentUserID := context.UserIDFromContext(ctx)

	GetChatRequest.UserIDs = append(GetChatRequest.UserIDs, currentUserID)
	if len(GetChatRequest.UserIDs) < 2 {
		ctx.String(http.StatusBadRequest, "userIDs cannot be less than 2")
		return
	}

	chat, err := h.csvc.GetChatByUserIDs(GetChatRequest.UserIDs)
	if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// create chat if not exists
	if errors.Is(err, domain.ErrRecordNotFound) {
		chat, err = h.csvc.NewChatWithParticipants(GetChatRequest.UserIDs)
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
