package main

import (
	"context"
	"log"
	"time"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/primary/grpc"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/secondary/mailer"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/secondary/postgresql"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/secondary/provider"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/secondary/provider/google"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/secondary/provider/local"
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/services"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
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

	providers := map[string]ports.Provider{}
	oauthProviders := map[string]ports.OAuthProvider{}

	googleOIDC, err := google.NewGoogleOIDC(ctx, config.GetGoogleClientID(), config.GetGoogleClientSecret(), google.CallbackURL, repo)
	if err != nil {
		log.Printf("google oidc disabled: %v", err)
	} else {
		providers[google.GoogleProviderName] = googleOIDC
		oauthProviders[google.GoogleProviderName] = googleOIDC
	}

	mailer := mailer.NewMailer()
	localProvider := local.NewLocalProvider(mailer, repo)
	providers[local.LocalProviderName] = localProvider

	registry := &provider.ProviderRegistry{
		Providers:      providers,
		OauthProviders: oauthProviders,
	}

	// get service
	tsvc := services.NewTokenService(repo)
	asvc := services.NewAuthService(repo, tsvc, registry)
	usvc := services.NewUserService(repo)

	// make grpc server and listen
	grpcAdapter := grpc.NewAdapter(9001, asvc, usvc)

	if err := grpcAdapter.Run(); err != nil {
		log.Fatalf("error starting grpc server: %v\n", err)
	}
}
