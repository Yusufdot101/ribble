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
	Name     string
	Email    string
	Tokens   []Token        `gorm:"constraint:OnDelete:CASCADE;"`
	Identity []UserIdentity `gorm:"constraint:OnDelete:CASCADE;"`
}

func (a *Adapter) InsertUser(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userModel := &User{
		Name:  user.Name,
		Email: user.Email,
	}
	res := a.DB.WithContext(ctx).Create(userModel)
	if res.Error == nil {
		user.ID = userModel.ID
	}

	return res.Error
}

func (a *Adapter) FindUserByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userModel := &User{}
	err := a.DB.WithContext(ctx).Where("email = ?", email).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	user := &domain.User{
		ID:        userModel.ID,
		Name:      userModel.Name,
		Email:     userModel.Email,
		CreatedAt: userModel.CreatedAt,
	}
	return user, nil
}

func (a *Adapter) FindUsersByID(ctx context.Context, userIDs []uint32) ([]*domain.User, error) {
	var userModels []User
	res := a.DB.WithContext(ctx).Where("id IN ?", userIDs).Find(&userModels)
	var users []*domain.User
	for _, userModel := range userModels {
		users = append(users, &domain.User{
			ID:        userModel.ID,
			Name:      userModel.Name,
			Email:     userModel.Email,
			CreatedAt: userModel.CreatedAt,
		})
	}
	return users, res.Error
}

func (a *Adapter) FindUsersByEmail(email string) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userModels []User
	var res *gorm.DB
	if email == "" {
		res = a.DB.WithContext(ctx).Find(&userModels)
	} else {
		// to_tsvector('simple', email) @@ plainto_tsquery('simple', ?)
		res = a.DB.WithContext(ctx).Where(`
			email ILIKE ?
			`, "%"+email+"%").Find(&userModels)
	}

	var users []*domain.User
	for _, userModel := range userModels {
		users = append(users, &domain.User{
			ID:        userModel.ID,
			Name:      userModel.Name,
			Email:     userModel.Email,
			CreatedAt: userModel.CreatedAt,
		})
	}
	return users, res.Error
}

func (a *Adapter) SearchUsers(ctx context.Context, query string, ids []uint32) ([]*domain.User, error) {
	if len(ids) == 0 {
		return []*domain.User{}, nil
	}

	var userModels []User
	var res *gorm.DB
	if query == "" {
		res = a.DB.WithContext(ctx).Where("id IN ?", ids).Find(&userModels)
	} else {
		res = a.DB.WithContext(ctx).Where("id IN ?", ids).Where(`
			email ILIKE ? OR name ILIKE ?
			`, "%"+query+"%", "%"+query+"%").Find(&userModels)
	}

	var users []*domain.User
	for _, userModel := range userModels {
		users = append(users, &domain.User{
			ID:        userModel.ID,
			Name:      userModel.Name,
			Email:     userModel.Email,
			CreatedAt: userModel.CreatedAt,
		})
	}
	return users, res.Error
}

func (a *Adapter) GetContacts(ctx context.Context, query string, excludeIds []uint32, currentUserID uint32) ([]*domain.User, error) {
	userModels := []*User{}
	tx := a.DB.WithContext(ctx).Model(&User{}).Where("id != ?", currentUserID)

	if query != "" {
		searchTerm := "%" + query + "%"
		tx = tx.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	if len(excludeIds) > 0 {
		tx = tx.Where("id NOT IN ?", excludeIds)
	}
	if err := tx.Find(&userModels).Error; err != nil {
		return nil, err
	}

	users := []*domain.User{}
	for _, userModel := range userModels {
		users = append(users, &domain.User{
			ID:        userModel.ID,
			Name:      userModel.Name,
			Email:     userModel.Email,
			CreatedAt: userModel.CreatedAt,
		})
	}
	return users, nil
}
