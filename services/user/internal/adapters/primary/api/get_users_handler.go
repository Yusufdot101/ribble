package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) getUsersByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	users, err := h.usvc.GetUsersByEmail(email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
