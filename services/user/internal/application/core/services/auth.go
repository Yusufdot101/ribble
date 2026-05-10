package services

import (
	"context"
	"errors"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
)

type AuthService struct {
	repo     ports.Repository
	registry ports.AuthProviderRegistry
	tsvc     ports.TokenService
}

func NewAuthService(repo ports.Repository, tsvc ports.TokenService, registry ports.AuthProviderRegistry) *AuthService {
	return &AuthService{
		repo:     repo,
		registry: registry,
		tsvc:     tsvc,
	}
}

func (asvc *AuthService) NewUser(user *domain.User) error {
	return asvc.repo.InsertUser(user)
}

func (asvc *AuthService) HandleAuth(ctx context.Context, credentials map[string]string, provider string) (string, string, error) {
	p, err := asvc.registry.GetProvider(provider)
	if err != nil {
		return "", "", err
	}

	identity, err := p.Authenticate(ctx, credentials)
	if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
		return "", "", err
	}
	if !identity.EmailVerified {
		return "", "", domain.ErrUnverifiedAccount
	}

	if errors.Is(err, domain.ErrRecordNotFound) {
		user, err := asvc.repo.FindUserByEmail(identity.Email)
		if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
			return "", "", err
		}

		// TODO: make use of transaction
		if errors.Is(err, domain.ErrRecordNotFound) {
			// create entry
			user = &domain.User{
				Email: identity.Email,
				Name:  identity.Name,
			}
			err = asvc.repo.InsertUser(user)
			if err != nil {
				return "", "", err
			}
		}
		identity.UserID = user.ID
		err = asvc.repo.InsertIdentity(identity)
		if err != nil {
			return "", "", err
		}
	}

	refreshToken, err := asvc.tsvc.New(domain.UUID, domain.REFRESH, identity.UserID)
	if err != nil {
		return "", "", err
	}

	err = asvc.tsvc.Save(refreshToken)
	if err != nil {
		return "", "", err
	}

	accessToken, err := asvc.tsvc.New(domain.JWT, domain.ACCESS, identity.UserID)
	if err != nil {
		return "", "", err
	}

	return refreshToken.TokenString, accessToken.TokenString, nil
}

func (asvc *AuthService) BeginAuth(provider string) (string, string, string, error) {
	state := generateUUID()
	nonce := generateUUID()
	p, err := asvc.registry.GetOauthProvider(provider)
	if err != nil {
		return "", "", "", err
	}

	url := p.GetAuthURL(state, nonce)
	return url, state, nonce, nil
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
