package google

import (
	"context"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const (
	issuerURL   = "https://accounts.google.com"
	CallbackURL = "http://localhost:8080/auth/google/callback"
)

type GoogleOIDC struct {
	config   *oauth2.Config
	provider *oidc.Provider
}

func NewGoogleOIDC(ctx context.Context, clientID, clientSecret, redirectURL string) (*GoogleOIDC, error) {
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
		config:   cfg,
		provider: provider,
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

func (g *GoogleOIDC) Authenticate(ctx context.Context, credentials map[string]string) (*domain.UserInfo, error) {
	code, hasCode := credentials["code"]
	nonce, hasNonce := credentials["nonce"]
	if !hasCode || !hasNonce || code == "" || nonce == "" {
		return nil, domain.ErrInvalidGoogleOIDCInputs
	}

	return g.getUserInfo(ctx, code, nonce)
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
		Provider: "google",
		Sub:      claims.Sub,
		Email:    claims.Email,
		Name:     claims.Name,
	}, nil
}

func (g *GoogleOIDC) GetAuthURL(state, nonce string) string {
	url := g.config.AuthCodeURL(state, oidc.Nonce(nonce))
	return url
}
