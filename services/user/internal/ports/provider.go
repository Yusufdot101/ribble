package ports

import (
	"context"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
)

type AuthProviderRegistry interface {
	GetProvider(provider string) (Provider, error)
	GetOauthProvider(provider string) (OAuthProvider, error)
}

type Provider interface {
	Authenticate(ctx context.Context, inputs map[string]string) (*domain.UserInfo, error)
}

type OAuthProvider interface {
	Provider
	GetAuthURL(state, nonce string) string
}
