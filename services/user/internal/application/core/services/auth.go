package services

import (
	"context"
	"errors"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
)

type AuthService struct {
	repo     ports.Repository
	provider ports.OAuthProvider
	tsvc     ports.TokenService
}

func NewAuthService(repo ports.Repository, provider ports.OAuthProvider, tsvc ports.TokenService) *AuthService {
	return &AuthService{
		repo:     repo,
		provider: provider,
		tsvc:     tsvc,
	}
}

func (asvc *AuthService) NewUser(user *domain.User) error {
	return asvc.repo.InsertUser(user)
}

func (asvc *AuthService) HandleCallback(ctx context.Context, inputs map[string]string) (string, string, error) {
	user, err := asvc.provider.GetUserInfo(ctx, inputs)
	if err != nil {
		return "", "", err
	}

	gotUser, err := asvc.repo.FindUserByProviderAndSub(user.Provider, user.Sub)
	if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
		return "", "", err
	} else if errors.Is(err, domain.ErrRecordNotFound) {
		// new user
		err = asvc.repo.InsertUser(user)
		if err != nil {
			return "", "", err
		}
	} else {
		user.ID = gotUser.ID
	}

	refreshToken, err := asvc.tsvc.New(domain.UUID, domain.REFRESH, user.ID)
	if err != nil {
		return "", "", err
	}

	err = asvc.tsvc.Save(refreshToken)
	if err != nil {
		return "", "", err
	}

	accessToken, err := asvc.tsvc.New(domain.JWT, domain.ACCESS, user.ID)
	if err != nil {
		return "", "", err
	}

	return refreshToken.TokenString, accessToken.TokenString, nil
}

func (asvc *AuthService) BeginAuth() (string, string, string) {
	state := generateUUID()
	nonce := generateUUID()

	url := asvc.provider.GetAuthURL(state, nonce)
	return url, state, nonce
}

func (asvc *AuthService) VerifyUsers(ctx context.Context, userIDs []uint32) (bool, error) {
	users, err := asvc.repo.FindUsersByID(ctx, userIDs)
	if err != nil {
		return false, err
	}
	if len(userIDs) != len(users) || len(userIDs) == 0 {
		return false, domain.ErrInvalidID
	}

	return true, nil
}
