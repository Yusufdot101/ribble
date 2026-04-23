package domain

type (
	RoleType       string
	PermissionType string
)

var (
	Admin  RoleType = "admin"
	Member RoleType = "member"

	ReadMessage  PermissionType = "read message"
	WriteMessage PermissionType = "write message"
)

type Permission struct {
	ID   uint
	Name PermissionType
}

type Role struct {
	ID   uint
	Name RoleType
}

type RolePermission struct {
	RoleID       uint
	PermissionID uint
}

func NewRole(name RoleType) *Role {
	return &Role{
		Name: name,
	}
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
