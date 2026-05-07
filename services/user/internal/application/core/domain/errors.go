package domain

import "errors"

var (
	ErrRecordNotFound          = errors.New("record not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrInvalidToken            = errors.New("invalid token")
	ErrNoIDToken               = errors.New("no id_token field in oauth2 token")
	ErrInvalidNonce            = errors.New("invalid nonce in id_token")
	ErrInvalidTokenUse         = errors.New("invalid token use")
	ErrInvalidTokeType         = errors.New("invalid token type")
	ErrInvalidJWT              = errors.New("unexpected signing method")
	ErrInvalidJWTSigningMethod = errors.New("unexpected signing method")

	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidEmail = errors.New("invalid email")

	ErrInvalidGoogleOIDCInputs = errors.New("invalid google oidc inputs")
)
