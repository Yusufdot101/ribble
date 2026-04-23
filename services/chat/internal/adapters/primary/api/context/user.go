package context

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserIDFromContext(ctx *gin.Context) uint {
	currentUserID, ok := ctx.MustGet("userID").(string)
	if !ok {
		panic("user id missing")
	}
	currentUserIDint, err := strconv.Atoi(currentUserID)
	if err != nil {
		panic("invalid user id type")
	}
	return uint(currentUserIDint)
}
