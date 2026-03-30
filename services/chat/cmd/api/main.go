package main

import (
	"log"

	"github.com/Yusufdot101/ripple/services/chat/config"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/primary/api"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/secondary/grpc"
	"github.com/Yusufdot101/ripple/services/chat/internal/adapters/secondary/postgresql"
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/services"
)

func main() {
	repo, err := postgresql.NewAdapter(config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("error : %v", err)
	}

	grpcClient, err := grpc.NewAdapter(9001)
	if err != nil {
		log.Fatalf("error initiating user grpc client: %v", err)
	}
	defer grpcClient.Close()

	csvc := services.NewChatService(repo, grpcClient)
	srv := api.NewServer(csvc)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}
