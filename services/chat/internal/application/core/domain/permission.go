package domain

type PermissionType string

const (
	SendMessage         PermissionType = "send message"
	AddToGroup          PermissionType = "add users to group"
	RemoveUserFromGroup PermissionType = "remove users from group"
	DeleteMessages      PermissionType = "delete messages"
)

type Permission struct {
	ID   uint           `json:"id"`
	Name PermissionType `json:"name"`
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
