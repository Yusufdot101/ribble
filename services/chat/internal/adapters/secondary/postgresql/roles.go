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
	chatRole.RoleID = chatRoleModel.RoleID
	return nil
}

func (a *Adapter) GrantChatRolePermission(roleName domain.RoleType, chatID uint, permission domain.PermissionType) error {
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

	chatRoleModel := &ChatRole{}
	err = a.db.WithContext(ctx).
		Table("chat_roles AS cr").
		Joins("JOIN roles AS r ON cr.role_id = r.id").
		Where("cr.chat_id = ? AND r.name = ?", chatID, roleName).
		First(chatRoleModel).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrInvalidPermission
		}
		return err
	}

	chatRolePermissionModel := &ChatRolePermission{
		ChatRoleID:   chatRoleModel.ID,
		PermissionID: permissionModel.ID,
	}

	err = a.db.WithContext(ctx).Save(chatRolePermissionModel).Error
	return err
}

func (a *Adapter) RevokeChatRolePermission(
	roleName domain.RoleType, chatID uint, permission domain.PermissionType,
) error {
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

	chatRoleModel := &ChatRole{}
	err = a.db.WithContext(ctx).
		Table("chat_roles AS cr").
		Joins("JOIN roles AS r ON cr.role_id = r.id").
		Where("cr.chat_id = ? AND r.name = ?", chatID, roleName).
		First(chatRoleModel).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrInvalidPermission
		}
		return err
	}

	res := a.db.WithContext(ctx).
		Where("chat_role_id = ? AND permission_id = ?", chatRoleModel.ID, permissionModel.ID).
		Delete(&ChatRolePermission{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return domain.ErrRecordNotFound
	}
	return nil
}

func (a *Adapter) GrantUsersChatRoles(userIDs []uint, chatID uint, roleName domain.RoleType) error {
	if len(userIDs) == 0 {
		return nil
	}
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
		Where("cp.user_id IN (?) AND cp.chat_id = ?", userIDs, chatID).
		Updates(map[string]any{
			"chat_role_id": chatRoleModel.ID,
		})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("no participant updated (invalid user/chat)")
	}
	if int(res.RowsAffected) != len(userIDs) {
		return errors.New("partial role grant: some users are not chat participants")
	}

	return nil
}

func (a *Adapter) GetUserRole(userID, chatID uint) (*domain.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roleModel := &Role{}

	err := a.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN chat_roles ON chat_roles.role_id = roles.id").
		Joins("JOIN chat_participants ON chat_participants.chat_role_id = chat_roles.id").
		Where("chat_participants.user_id = ? AND chat_participants.chat_id = ?", userID, chatID).
		First(roleModel).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrInvalidRole
		}
		return nil, err
	}

	role := &domain.Role{
		ID:   roleModel.ID,
		Name: roleModel.Name,
	}

	return role, nil
}

func (a *Adapter) GetRoleByChatRoleID(chatRoleID uint) (*domain.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roleModel := &Role{}

	err := a.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN chat_roles ON chat_roles.role_id = roles.id").
		Where("chat_roles.id = ?", chatRoleID).
		First(roleModel).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrInvalidRole
		}
		return nil, err
	}

	role := &domain.Role{
		ID:   roleModel.ID,
		Name: roleModel.Name,
	}

	return role, nil
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}
