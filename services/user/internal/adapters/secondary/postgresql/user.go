package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string         `gorm:"index:idx_email,unique"`
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
	if isDuplicateViolation(res.Error) {
		return domain.ErrDuplicateEmail
	}
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

	tx := a.DB.WithContext(ctx).
		Joins("JOIN user_identities ON users.id = user_identities.user_id").
		Distinct("users.id, users.name, users.email, users.created_at, users.updated_at, users.deleted_at").
		Where("users.id IN ? AND user_identities.email_verified = true", ids)

	var userModels []User

	if query != "" {
		searchTerm := "%" + query + "%"
		tx = tx.Where(`
			users.email ILIKE ? OR users.name ILIKE ?
			`, searchTerm, searchTerm)
	}
	if err := tx.Find(&userModels).Error; err != nil {
		return nil, err
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
	return users, nil
}

func (a *Adapter) GetContacts(ctx context.Context, query string, excludeIds []uint32, currentUserID uint32) ([]*domain.User, error) {
	userModels := []*User{}
	tx := a.DB.WithContext(ctx).
		Joins("JOIN user_identities ON users.id = user_identities.user_id").
		Model(&User{}).
		Distinct("users.id, users.name, users.email, users.created_at, users.updated_at, users.deleted_at").
		Where("users.id != ? AND user_identities.email_verified = true", currentUserID)

	if query != "" {
		searchTerm := "%" + query + "%"
		tx = tx.Where("users.name ILIKE ? OR users.email ILIKE ?", searchTerm, searchTerm)
	}

	if len(excludeIds) > 0 {
		tx = tx.Where("users.id NOT IN ?", excludeIds)
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

func isDuplicateViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
