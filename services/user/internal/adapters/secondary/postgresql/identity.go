package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"gorm.io/gorm"
)

type UserIdentity struct {
	gorm.Model
	Provider string `gorm:"index:idx_provider_sub,unique"`
	Sub      string `gorm:"index:idx_provider_sub,unique"`
	UserID   uint   `gorm:"index"`
}

func (a *Adapter) InsertIdentity(identity *domain.UserIdentity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	identityModel := &UserIdentity{
		Provider: identity.Provider,
		Sub:      identity.Sub,
		UserID:   identity.UserID,
	}
	res := a.DB.WithContext(ctx).Create(identityModel)
	if res.Error == nil {
		identity.ID = identityModel.ID
	}

	return res.Error
}

func (a *Adapter) FindIdentityByProviderAndSub(provider, sub string) (*domain.UserIdentity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	identityModel := &UserIdentity{}
	res := a.DB.WithContext(ctx).First(identityModel, "provider = ? AND sub = ?", provider, sub)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, res.Error
	}

	identity := &domain.UserIdentity{
		ID:        identityModel.ID,
		CreatedAt: identityModel.CreatedAt,
		Provider:  identityModel.Provider,
		Sub:       identityModel.Sub,
		UserID:    identityModel.UserID,
	}
	return identity, nil
}
