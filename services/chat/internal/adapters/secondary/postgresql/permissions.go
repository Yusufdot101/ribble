package postgresql

import (
	"context"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name               domain.PermissionType
	ChatRolePermission []ChatRolePermission `gorm:"constraint:OnDelete:CASCADE;"`
}

type ChatRolePermission struct {
	gorm.Model
	ChatRoleID   uint
	PermissionID uint
}

func (a *Adapter) NewPermission(permission *domain.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	permissionModel := &Permission{
		Name: permission.Name,
	}

	err := a.db.WithContext(ctx).Save(permissionModel).Error
	if err == nil {
		permission.ID = permissionModel.ID
	}
	return err
}

func (a *Adapter) GetUserPermissions(userID, chatID uint) ([]*domain.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	permissionModels := []*Permission{}

	// chat_participants: role_id -> role_permissions: permission_id -> permissions
	err := a.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN chat_role_permissions ON chat_role_permissions.permission_id = permissions.id").
		Joins("JOIN chat_participants ON chat_participants.chat_role_id = chat_role_permissions.chat_role_id").
		Where("chat_participants.user_id = ? AND chat_participants.chat_id = ?", userID, chatID).
		Find(&permissionModels).
		Error
	if err != nil {
		return nil, err
	}

	permissions := []*domain.Permission{}
	for _, permission := range permissionModels {
		permission := &domain.Permission{
			ID:   permission.ID,
			Name: permission.Name,
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}
