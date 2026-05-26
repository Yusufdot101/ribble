package google

import (
	"context"
	"errors"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var CallbackURL = config.GetRedirectURL()

const (
	issuerURL          = "https://accounts.google.com"
	GoogleProviderName = "google"
)

type GoogleOIDC struct {
	config   *oauth2.Config
	provider *oidc.Provider

	ProviderName string
	repo         ports.IdentityRepository
}

func NewGoogleOIDC(ctx context.Context, clientID, clientSecret, redirectURL string, repo ports.IdentityRepository) (*GoogleOIDC, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}

	cfg := &oauth2.Config{
		ClientID:     config.GetGoogleClientID(),
		ClientSecret: config.GetGoogleClientSecret(),
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &GoogleOIDC{
		config:       cfg,
		provider:     provider,
		repo:         repo,
		ProviderName: GoogleProviderName,
	}, nil
}

func (g *GoogleOIDC) Authenticate(ctx context.Context, credentials map[string]string) (*domain.UserIdentity, error) {
	code, hasCode := credentials["code"]
	nonce, hasNonce := credentials["nonce"]
	if !hasCode || !hasNonce || code == "" || nonce == "" {
		return nil, domain.ErrInvalidProviderInputs
	}

	userInfo, err := g.getUserInfo(ctx, code, nonce)
	if err != nil {
		return nil, err
	}

	identity, err := g.repo.FindIdentityByProviderAndSub(userInfo.Provider, userInfo.Sub)
	if err == nil {
		identity.EmailVerified = true
		return identity, nil
	}
	if !errors.Is(err, domain.ErrRecordNotFound) {
		return nil, err
	}

	return &domain.UserIdentity{
		Provider:      userInfo.Provider,
		Sub:           userInfo.Sub,
		EmailVerified: userInfo.EmailVerified,
		Email:         userInfo.Email,
		Name:          userInfo.Name,
	}, err
}

func (g *GoogleOIDC) GetAuthURL(state, nonce string) string {
	url := g.config.AuthCodeURL(state, oidc.Nonce(nonce))
	return url
}

func (g *GoogleOIDC) getUserInfo(ctx context.Context, code, nonce string) (*domain.UserInfo, error) {
	rawIDToken, err := g.exchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	idToken, err := g.verifyIDToken(ctx, rawIDToken, nonce)
	if err != nil {
		return nil, err
	}

	var claims struct {
		Sub   string `json:"sub"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	// all the OIDC token exchange + verification code lives here
	return &domain.UserInfo{
		Provider:      g.ProviderName,
		Sub:           claims.Sub,
		Email:         claims.Email,
		Name:          claims.Name,
		EmailVerified: true,
	}, nil
}

func (g *GoogleOIDC) verifyIDToken(ctx context.Context, rawIDToken string, expectedNonce string) (*oidc.IDToken, error) {
	verifier := g.provider.Verifier(&oidc.Config{ClientID: g.config.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	if expectedNonce == "" || idToken.Nonce != expectedNonce {
		return nil, domain.ErrInvalidNonce
	}

	return idToken, nil
}

func (g *GoogleOIDC) exchangeCode(ctx context.Context, code string) (string, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", domain.ErrNoIDToken
	}

	return rawIDToken, nil
}
