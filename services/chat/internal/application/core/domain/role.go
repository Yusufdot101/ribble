package domain

type RoleType string

var (
	Admin  RoleType = "admin"
	Member RoleType = "member"
)

type Role struct {
	ID   uint
	Name RoleType
}
type ChatRole struct {
	ID     uint
	ChatID uint
	RoleID uint
}

func NewRole(name RoleType) *Role {
	return &Role{
		Name: name,
	}
}

func NewChatRole(chatID uint) *ChatRole {
	return &ChatRole{
		ChatID: chatID,
	}
}
