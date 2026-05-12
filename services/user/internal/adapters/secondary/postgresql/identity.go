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
	Provider      string
	Sub           string
	EmailVerified bool
	UserID        uint `gorm:"index"`
	PasswordHash  *[]byte
}

func (a *Adapter) InsertIdentity(identity *domain.UserIdentity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	identityModel := &UserIdentity{
		Provider:      identity.Provider,
		Sub:           identity.Sub,
		UserID:        identity.UserID,
		PasswordHash:  identity.PasswordHash,
		EmailVerified: identity.EmailVerified,
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
	res := a.DB.WithContext(ctx).First(identityModel, "provider = ? AND sub = ? AND email_verified = true", provider, sub)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, res.Error
	}

	identity := &domain.UserIdentity{
		ID:            identityModel.ID,
		CreatedAt:     identityModel.CreatedAt,
		Provider:      identityModel.Provider,
		Sub:           identityModel.Sub,
		UserID:        identityModel.UserID,
		EmailVerified: identityModel.EmailVerified,
		PasswordHash:  identityModel.PasswordHash,
	}
	return identity, nil
}

func (a *Adapter) FindIdentityByUserIDAndID(userID, identityID uint) (*domain.UserIdentity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	identityModel := &UserIdentity{}
	res := a.DB.WithContext(ctx).
		Joins("JOIN users ON users.id = user_identities.user_id").
		Where("users.id = ? AND user_identities.id = ?", userID, identityID).
		First(identityModel)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, res.Error
	}

	identity := &domain.UserIdentity{
		ID:            identityModel.ID,
		CreatedAt:     identityModel.CreatedAt,
		Provider:      identityModel.Provider,
		Sub:           identityModel.Sub,
		UserID:        identityModel.UserID,
		EmailVerified: identityModel.EmailVerified,
		PasswordHash:  identityModel.PasswordHash,
	}
	return identity, nil
}

func (a *Adapter) ActivateIdentity(identityID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := a.DB.WithContext(ctx).
		Table("user_identities AS ui").
		Where("ui.id = ? AND email_verified = false", identityID).
		Updates(map[string]any{
			"email_verified": true,
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrRecordNotFound
	}

	return nil
}
