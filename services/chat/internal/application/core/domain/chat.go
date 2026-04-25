package domain

type CreateChatWithParticipantsRequestType struct {
	RolePermissions map[string][]string `json:"rolePermissions"`
	UserRoles       map[uint]string     `json:"userRoles"`
}

type Chat struct {
	ID uint
}

func NewChat() *Chat {
	return &Chat{}
}
