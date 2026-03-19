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

func NewToken(use TokenUse, tokenType TokenType, userID uint, tokenString string, expires time.Time) *Token {
	return &Token{
		Use:         use,
		TokenType:   tokenType,
		TokenString: tokenString,
		Expires:     expires,
		UserID:      userID,
	}
}
