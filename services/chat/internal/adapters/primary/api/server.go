package api

import (
	"log"
	"net/http"

	"github.com/Yusufdot101/ripple/services/chat/internal/ports"
	"github.com/gin-gonic/gin"
)

type handler struct {
	csvc ports.ChatService
}

type Server struct {
	r *gin.Engine
}

func NewServer(csvc ports.ChatService) *Server {
	h := handler{
		csvc: csvc,
	}
	r := h.RegisterRoutes()
	return &Server{
		r: r,
	}
}

const PORT = ":8081"

func (s *Server) ListenAndServe() error {
	log.Printf("server listening on port %v\n", PORT)
	return http.ListenAndServe(PORT, s.r)
}
