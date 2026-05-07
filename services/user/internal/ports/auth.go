package ports

import (
	"context"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
)

type AuthService interface {
	NewUser(user *domain.User) error
	BeginAuth() (authURL, state, nonce string)
	HandleCallback(ctx context.Context, inputs map[string]string) (refreshToken, accessToken string, err error)
	VerifyUsers(ctx context.Context, userIDs []uint32) (bool, error)
}
