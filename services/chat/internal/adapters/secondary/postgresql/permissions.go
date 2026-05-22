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

	err := a.db.WithContext(ctx).
		Table("permissions p").
		Joins("JOIN chat_role_permissions crp ON crp.permission_id = p.id").
		Joins("JOIN chat_participants cp ON cp.chat_role_id = crp.chat_role_id").
		Where("cp.user_id = ? AND cp.chat_id = ? AND crp.deleted_at IS NULL", userID, chatID).
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
