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
	"github.com/gin-gonic/gin"
)

func main() {
	env := config.GetEnv()
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// get repo
	repo, err := postgresql.NewAdapter(config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("error : %v", err)
	}

	// get provider
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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

	tsvc := services.NewTokenService(repo)
	mailer := mailer.NewMailer(
		config.GetSMTPHost(), config.GetSMTPPort(), config.GetSMTPUsername(),
		config.GetSMTPPassword(), config.GetSMTPSender(),
	)
	localProvider := local.NewLocalProvider(mailer, repo, tsvc)
	providers[local.LocalProviderName] = localProvider

	registry := &provider.ProviderRegistry{
		Providers:      providers,
		OauthProviders: oauthProviders,
	}
	asvc := services.NewAuthService(repo, tsvc, registry)
	usvc := services.NewUserService(repo)

	// make server listen
	server := api.NewServer(config.GetPort(8080), asvc, tsvc, usvc)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}
