package ports

import (
	"context"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
)

type IdentityRepository interface {
	InsertIdentity(identity *domain.UserIdentity) error
	FindIdentityByProviderAndSub(provider, sub string) (*domain.UserIdentity, error)
	FindIdentityByUserIDAndID(userID, identityID uint) (*domain.UserIdentity, error)
	FindUserByEmail(email string) (*domain.User, error)
	InsertUser(user *domain.User) error
	WithTx(fn func(repo Repository) error) error
	ActivateIdentity(identityID uint) error
}

type Repository interface {
	IdentityRepository
	FindUsersByID(context.Context, []uint32) ([]*domain.User, error)
	FindUsersByEmail(email string) ([]*domain.User, error)
	SearchUsers(ctx context.Context, query string, ids []uint32) ([]*domain.User, error)
	GetContacts(ctx context.Context, query string, excludeIds []uint32, currentUserID uint32) ([]*domain.User, error)

	InsertToken(token *domain.Token) error
	GetTokenByStringAndUse(tokenString string, tokenUse domain.TokenUse) (*domain.Token, error)
	DeleteTokenByStringAndUse(tokenString string, tokenUse domain.TokenUse) error
}
