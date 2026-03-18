package domain

import (
	"time"
)

type (
	TokenType string
	TokenUse  string
)

const (
	UUID TokenType = "uuid"
	JWT  TokenType = "jwt"

	REFRESH TokenUse = "refresh"
	ACCESS  TokenUse = "access"
)

type Token struct {
	ID          uint
	TokenString string
	CreatedAt   time.Time
	UserID      uint
	Expires     time.Time
	Use         TokenUse
	TokenType   TokenType
}
