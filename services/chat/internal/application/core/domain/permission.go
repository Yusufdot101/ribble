package domain

type PermissionType string

var SendMessage PermissionType = "send message"

type Permission struct {
	ID   uint
	Name PermissionType
}

type RolePermission struct {
	ChatRoleID   uint
	PermissionID uint
}

func NewPermission(name PermissionType) *Permission {
	return &Permission{
		Name: name,
	}
}

func (p *Permission) IncludedIn(permissions []*Permission) bool {
	for _, permission := range permissions {
		if p.Name == permission.Name {
			return true
		}
	}
	return false
}
