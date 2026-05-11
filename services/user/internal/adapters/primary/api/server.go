package api

import (
	"fmt"
	"net/http"

	"github.com/Yusufdot101/ripple/services/user/internal/ports"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	port   int
}

func NewServer(port int, svc ports.AuthService, tsvc ports.TokenService, usvc ports.UserService) *Server {
	h := NewHandler(svc, tsvc, usvc)
	r := h.RegisterRoutes()
	return &Server{
		router: r,
		port:   port,
	}
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}
