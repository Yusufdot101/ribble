package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/ports"
	"github.com/gin-gonic/gin"
)

type handler struct {
	csvc ports.ChatService
	hub  *hub
}

type Server struct {
	r *gin.Engine
}

func NewServer(csvc ports.ChatService) *Server {
	h := handler{
		hub:  newHub(),
		csvc: csvc,
	}
	r := h.RegisterRoutes()
	return &Server{
		r: r,
	}
}

const PORT = ":8081"

func (s *Server) ListenAndServe() error {
	srv := http.Server{
		Addr:         PORT,
		Handler:      s.r,
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		sig := <-quit

		log.Println("server shutting down", map[string]string{
			"signal": sig.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		shutdownErr <- err
	}()
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	if err := <-shutdownErr; err != nil {
		return err
	}

	log.Println("server stopped")
	return nil
}
