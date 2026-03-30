package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Sub      string
	Name     string
	Email    string
	Provider string
	Tokens   []Token `gorm:"constraint:OnDelete:CASCADE;"`
}

func (a *Adapter) InsertUser(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userModel := &User{
		Name:     user.Name,
		Email:    user.Email,
		Sub:      user.Sub,
		Provider: user.Provider,
	}
	res := a.DB.WithContext(ctx).Create(userModel)
	if res.Error == nil {
		user.ID = userModel.ID
	}

	return res.Error
}

func (a *Adapter) FindUserByProviderAndSub(provider, sub string) (*domain.User, error) {
	userModel := &User{}
	res := a.DB.First(userModel, "provider = ? AND sub = ?", provider, sub)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, res.Error
	}

	user := &domain.User{
		ID:        userModel.ID,
		Name:      userModel.Name,
		Email:     userModel.Email,
		CreatedAt: userModel.CreatedAt,
		Provider:  userModel.Provider,
		Sub:       userModel.Sub,
	}
	return user, nil
}

func (a *Adapter) FindUsersByID(ctx context.Context, userIDs []uint32) ([]*domain.User, error) {
	var usersModel []User
	res := a.DB.WithContext(ctx).Where("id IN ?", userIDs).Find(&usersModel)
	var users []*domain.User
	for _, userModel := range usersModel {
		users = append(users, &domain.User{
			ID:        userModel.ID,
			Sub:       userModel.Sub,
			Provider:  userModel.Provider,
			Name:      userModel.Name,
			Email:     userModel.Email,
			CreatedAt: userModel.CreatedAt,
		})
	}
	return users, res.Error
}
