package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var createChatWithParticipantsRequests struct {
	UserIDs []uint `json:"userIDs"`
}

func (h *handler) NewChatWithParticipants(ctx *gin.Context) {
	if err := ctx.ShouldBind(&createChatWithParticipantsRequests); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	currentUserID, ok := ctx.MustGet("userID").(string)
	if !ok {
		panic("user id missing")
	}
	currentUserIDint, err := strconv.Atoi(currentUserID)
	if err != nil {
		panic("invalid user id type")
	}
	createChatWithParticipantsRequests.UserIDs = append(createChatWithParticipantsRequests.UserIDs, uint(currentUserIDint))
	if len(createChatWithParticipantsRequests.UserIDs) < 2 {
		ctx.String(http.StatusBadRequest, "userIDs cannot be less than 2")
		return
	}
	log.Println("here: ", createChatWithParticipantsRequests.UserIDs)
	chatID, err := h.csvc.NewChatWithParticipants(createChatWithParticipantsRequests.UserIDs)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"chatID": chatID,
	})
}
