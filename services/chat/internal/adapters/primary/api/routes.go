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
	}), middleware.RecoverPanic(), middleware.IPRateLimiter())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	convoGroup := r.Group("/conversations").Use(middleware.RequireAuthentication())
	convoGroup.GET("", h.getConversations)

	group := r.Group("/chats")
	group.Use(middleware.RequireAuthentication())
	group.POST("", h.GetOrCreateChat)
	group.GET("/:chatId", h.getChatByID)
	group.GET("/:chatId/users", h.getChatUsers)

	group.GET("/:chatId/addable-users", h.getAddableChatUsers)
	group.POST("/:chatId/addToGroup", h.addToGroup)
	group.DELETE("/:chatId/users/:userId", h.removeFromGroup)
	group.GET("/:chatId/permissions", h.getUserPermissions)
	group.PATCH("/:chatId/users/:userId", h.updateUserRole)
	group.PATCH("/:chatId/roles/:role/permissions", h.updateRolePermission)
	group.POST("/:chatId/roles/permissions", h.getRolePermissions)

	group.GET("/:chatId/bans", h.getBannedUsers)
	group.POST("/:chatId/bans", h.banFromGroup)
	group.DELETE("/:chatId/bans/:userId", h.unbanFromGroup)

	messageGroup := group.Group("/:chatId/messages")
	messageGroup.GET("", h.getMessages)
	messageGroup.GET("/sync", h.syncMessages)
	messageGroup.DELETE(":messageId", h.deleteMessage)
	messageGroup.PATCH(":messageId", h.editMessage)

	r.GET("/ws", h.newWebsocket)
	return r
}
