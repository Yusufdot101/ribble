package api

import (
	"net/http"

	"github.com/Yusufdot101/ripple/services/user/config"
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
	}))

	group := r.Group("/auth")
	group.GET("/google", h.googleBegin)
	group.GET("/google/callback", h.googleCallback)
	group.GET("/refreshtoken", h.RefreshToken)

	group.POST("/signup", h.register)
	group.GET("/verify", h.verify)
	group.POST("/login", h.login)

	group.Match([]string{http.MethodPost, http.MethodOptions}, "/logout", h.logout)

	userGroup := r.Group("/users")
	userGroup.GET("", h.getUsersByEmail)
	return r
}
