package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.router,
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		sig := <-quit

		log.Println("api server shutting down", map[string]string{
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

	log.Println("api server stopped")
	return nil
}
