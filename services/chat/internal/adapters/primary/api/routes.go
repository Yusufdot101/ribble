package api

import (
	"github.com/Yusufdot101/ripple/shared/middleware"
	"github.com/gin-gonic/gin"
)

func (h *handler) RegisterRoutes() *gin.Engine {
	r := gin.New()
	group := r.Group("/chat")
	group.POST("", middleware.RequireAuthentication(h.NewChatWithParticipants))
	group.POST("/find", middleware.RequireAuthentication(h.GetByUserIDs))
	return r
}
