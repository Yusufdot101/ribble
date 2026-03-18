package main

import (
	"context"
	"log"
	"time"

	"github.com/Yusufdot101/ribble/services/user/config"
	"github.com/Yusufdot101/ribble/services/user/internal/adapters/primary/api"
	"github.com/Yusufdot101/ribble/services/user/internal/adapters/secondary/google"
	"github.com/Yusufdot101/ribble/services/user/internal/adapters/secondary/postgresql"
	"github.com/Yusufdot101/ribble/services/user/internal/application/core/services"
)

func main() {
	// get repo
	repo, err := postgresql.NewAdapter(config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("error : %v", err)
	}

	// get provider
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	googleOIDC, err := google.NewGoogleOIDC(ctx, config.GetGoogleClientID(), config.GetGoogleClientSecret(), google.CallbackURL)
	if err != nil {
		log.Fatalf("error : %v", err)
	}

	// get service
	tsvc := services.NewTokenService(repo)
	asvc := services.NewAuthService(repo, googleOIDC, tsvc)

	// make server listen
	server := api.NewServer(asvc)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}
