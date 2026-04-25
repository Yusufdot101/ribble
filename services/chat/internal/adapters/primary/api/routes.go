package api

import (
	"net/http"

	"github.com/Yusufdot101/ripple/services/chat/config"
	"github.com/Yusufdot101/ripple/shared/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (h *handler) RegisterRoutes() *gin.Engine {
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{config.GetFrontendURL()},
		AllowMethods:     []string{http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodPatch},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
	}))
	group := r.Group("/chats")
	// group.POST("", middleware.RequireAuthentication(h.NewChatWithParticipants))
	group.POST("", middleware.RequireAuthentication(h.GetOrCreateChat))

	messageGroup := r.Group("/messages")
	messageGroup.POST("", middleware.RequireAuthentication(h.getMessages))
	messageGroup.DELETE(":id", middleware.RequireAuthentication(h.deleteMessage))
	messageGroup.PATCH(":id", middleware.RequireAuthentication(h.editMessage))
	messageGroup.GET("/new", h.newMessage)
	return r
}
