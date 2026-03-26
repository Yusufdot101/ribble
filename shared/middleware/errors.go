package middleware

import "errors"

var (
	ErrMissingInvalidToken     = errors.New("authorization token missing or error")
	ErrInvalidJWT              = errors.New("invalid jwt")
	ErrInvalidJWTSigningMethod = errors.New("unexpected signing method")
)
