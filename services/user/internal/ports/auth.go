package ports

import (
	"context"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
)

type AuthService interface {
	NewUser(user *domain.User) error
	BeginAuth(providerName string) (authURL, state, nonce string, err error)
	HandleAuth(ctx context.Context, credentials map[string]string, provider string) (refreshToken, accessToken string, err error)
	VerifyUsers(ctx context.Context, userIDs []uint32) (bool, error)
	ActivateAccount(tokenString string) (string, string, error)
}
