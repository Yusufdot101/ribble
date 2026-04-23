package postgresql

import (
	"context"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name           domain.PermissionType
	RolePermission []RolePermission `gorm:"constraint:OnDelete:CASCADE;"`
}

type Role struct {
	gorm.Model
	Name           domain.RoleType
	RolePermission []RolePermission `gorm:"constraint:OnDelete:CASCADE;"`
}

type RolePermission struct {
	gorm.Model
	RoleID       uint
	PermissionID uint
}

// create role
func (a *Adapter) NewRole(role *domain.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roleModel := &Role{
		Name: role.Name,
	}

	err := a.db.WithContext(ctx).Save(roleModel).Error
	if err == nil {
		role.ID = roleModel.ID
	}
	return err
}

// create permission
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
