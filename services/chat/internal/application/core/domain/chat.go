package domain

type CreateChatWithParticipantsRequestType struct {
	Name            string              `json:"name"`
	IsGroup         bool                `json:"isGroup"`
	RolePermissions map[string][]string `json:"rolePermissions"`
	UserRoles       map[uint]string     `json:"userRoles"`
}

type Chat struct {
	Name    string `json:"name"`
	IsGroup bool   `json:"isGroup"`
	ID      uint   `json:"id"`
}

func NewChat(name string, isGroup bool) *Chat {
	return &Chat{
		Name:    name,
		IsGroup: isGroup,
	}
}
