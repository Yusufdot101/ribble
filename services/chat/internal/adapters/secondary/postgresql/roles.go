package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name     domain.RoleType
	ChatRole []ChatRole `gorm:"constraint:OnDelete:CASCADE;"`
}

type ChatRole struct {
	gorm.Model
	ChatRolePermissions []ChatRolePermission `gorm:"constraint:OnDelete:CASCADE;"`
	ChatID              uint
	RoleID              uint
}

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

func (a *Adapter) NewChatRole(chatRole *domain.ChatRole, roleName domain.RoleType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roleModel := &Role{}
	err := a.db.WithContext(ctx).
		Where("name = ?", roleName).
		First(roleModel).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrInvalidRole
		}
		return err
	}

	chatRoleModel := &ChatRole{
		ChatID: chatRole.ChatID,
		RoleID: roleModel.ID,
	}

	err = a.db.WithContext(ctx).Save(chatRoleModel).Error
	if err != nil {
		if isForeignKeyViolation(err) {
			return domain.ErrInvalidChatRole
		}
		return err
	}

	chatRole.ID = chatRoleModel.ID
	return nil
}

func (a *Adapter) GrantChatRolePermission(chatRoleID uint, permission domain.PermissionType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the permission id
	permissionModel := &Permission{}
	err := a.db.WithContext(ctx).Where("name = ?", permission).First(permissionModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrInvalidPermission
		}
		return err
	}

	chatRolePermissionModel := &ChatRolePermission{
		ChatRoleID:   chatRoleID,
		PermissionID: permissionModel.ID,
	}

	err = a.db.WithContext(ctx).Save(chatRolePermissionModel).Error
	return err
}

func (a *Adapter) GrantUserChatRole(userID, chatID uint, roleName domain.RoleType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatRoleModel := &ChatRole{}
	err := a.db.WithContext(ctx).
		Joins("JOIN roles ON roles.id = chat_roles.role_id").
		Where("roles.name = ? AND chat_roles.chat_id = ?", roleName, chatID).
		First(chatRoleModel).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrInvalidRole
		}
		return err
	}

	res := a.db.WithContext(ctx).
		Table("chat_participants AS cp").
		Where("cp.user_id = ? AND cp.chat_id = ?", userID, chatID).
		Updates(map[string]any{
			"chat_role_id": chatRoleModel.ID,
		})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("no participant updated (invalid user/chat)")
	}
	return nil
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}
