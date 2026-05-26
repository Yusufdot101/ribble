package api

import (
	"net/http"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/shared/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (h *handler) RegisterRoutes() *gin.Engine {
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{config.GetFrontendURL()},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
	}), middleware.RecoverPanic(), middleware.IPRateLimiter())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	group := r.Group("/auth")
	group.GET("/google", h.googleBegin)
	group.GET("/google/callback", h.googleCallback)
	group.GET("/refreshtoken", h.RefreshToken)

	group.GET("/verify", h.verify)
	loginGroup := group.Group("")
	loginGroup.Use(middleware.CredentialRateLimiter())
	loginGroup.POST("/signup", h.register)
	loginGroup.POST("/login", h.login)

	group.Match([]string{http.MethodPost, http.MethodOptions}, "/logout", h.logout)

	userGroup := r.Group("/users")
	userGroup.GET("", h.getUsersByEmail)
	return r
}
