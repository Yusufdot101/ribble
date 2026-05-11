package main

import (
	"context"
	"log"
	"time"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/adapters/primary/api"
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	googleOIDC, err := google.NewGoogleOIDC(ctx, config.GetGoogleClientID(), config.GetGoogleClientSecret(), google.CallbackURL, repo)
	if err != nil {
		log.Printf("error : %v", err)
	}

	mailer := mailer.NewMailer()
	localProvider := local.NewLocalProvider(mailer, repo)

	registry := &provider.ProviderRegistry{
		Providers: map[string]ports.Provider{
			google.GoogleProviderName: googleOIDC,
			local.LocalProviderName:   localProvider,
		},
		OauthProviders: map[string]ports.OAuthProvider{
			google.GoogleProviderName: googleOIDC,
		},
	}
	tsvc := services.NewTokenService(repo)
	asvc := services.NewAuthService(repo, tsvc, registry)
	usvc := services.NewUserService(repo)

	// make server listen
	server := api.NewServer(asvc, tsvc, usvc)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}
