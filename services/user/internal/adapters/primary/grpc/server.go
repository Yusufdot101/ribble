package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	userpb "github.com/Yusufdot101/ripple-proto/golang/user/v4"
	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	userpb.UnimplementedUserServiceServer
	port int
	asvc ports.AuthService
	usvc ports.UserService
}

func NewAdapter(port int, asvc ports.AuthService, usvc ports.UserService) *Adapter {
	return &Adapter{
		port: port,
		asvc: asvc,
		usvc: usvc,
	}
}

func (a *Adapter) Run() error {
	address := fmt.Sprintf("0.0.0.0:%d", a.port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		sig := <-quit

		log.Println("grpc server shutting down", map[string]string{
			"signal": sig.String(),
		})
		grpcServer.GracefulStop()
	}()

	log.Printf("grpc server listening on :%s\n", address)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to server grpc server: %v", err)
	}

	log.Println("grpc server stopped")

	return nil
}
