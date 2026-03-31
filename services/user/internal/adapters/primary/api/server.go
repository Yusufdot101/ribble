package api

import (
	"log"
	"net/http"

	"github.com/Yusufdot101/ripple/services/user/internal/ports"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer(svc ports.AuthService, tsvc ports.TokenService, usvc ports.UserService) *Server {
	h := NewHandler(svc, tsvc, usvc)
	r := h.RegisterRoutes()
	return &Server{
		router: r,
	}
}

const PORT = ":8080"

func (s *Server) ListenAndServe() error {
	log.Printf("server listening on port %v\n", PORT)
	return http.ListenAndServe(PORT, s.router)
}
