package middleware

import (
	"github.com/Yusufdot101/ribble/shared/middleware/config"
	"github.com/golang-jwt/jwt/v4"
)

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// ensure the token was signed with HMAC, not something else
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidJWTSigningMethod
		}
		return config.GetJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidJWT
	}

	return token, nil
}
