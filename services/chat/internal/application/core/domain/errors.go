package domain

import "errors"

var (
	ErrInvalidUserIDs        = errors.New("invalid user ids")
	ErrRecordNotFound        = errors.New("record not found")
	ErrUpdateWindowOver      = errors.New("message update window over")
	ErrInvalidMessageContent = errors.New("invalid message content")

	ErrInvalidPermission = errors.New("invalid permission")
	ErrInvalidRole       = errors.New("invalid role")
	ErrInvalidChatRole   = errors.New("invalid chat role")
)
