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

	r.GET("/conversations", middleware.RequireAuthentication(h.getConversations))

	group := r.Group("/chats")
	group.POST("", middleware.RequireAuthentication(h.GetOrCreateChat))
	group.GET("/:chatId", middleware.RequireAuthentication(h.getChatByID))
	group.GET("/:chatId/users", middleware.RequireAuthentication(h.getChatUsers))
	group.GET("/:chatId/addable-users", middleware.RequireAuthentication(h.getAddableChatUsers))
	group.POST("/:chatId/addToGroup", middleware.RequireAuthentication(h.addToGroup))
	group.DELETE("/:chatId/users/:userId", middleware.RequireAuthentication(h.removeFromGroup))
	group.GET("/:chatId/permissions", middleware.RequireAuthentication(h.getUserPermissions))
	group.POST("/:chatId/ban", middleware.RequireAuthentication(h.banFromGroup))
	group.DELETE("/:chatId/ban/:userId", middleware.RequireAuthentication(h.unbanFromGroup))

	messageGroup := group.Group("/:chatId/messages")
	messageGroup.GET("", middleware.RequireAuthentication(h.getMessages))
	messageGroup.GET("/sync", middleware.RequireAuthentication(h.syncMessages))
	messageGroup.DELETE(":messageId", middleware.RequireAuthentication(h.deleteMessage))
	messageGroup.PATCH(":messageId", middleware.RequireAuthentication(h.editMessage))

	r.GET("/ws", h.newMessage)
	return r
}
