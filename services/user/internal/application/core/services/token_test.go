package services

import (
	"errors"
	"strconv"
	"testing"

	"github.com/Yusufdot101/ribble/services/user/config"
	"github.com/Yusufdot101/ribble/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ribble/services/user/internal/ports"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func TestGenerateUUID(t *testing.T) {
	uuidString := generateUUID()
	err := uuid.Validate(uuidString)
	if err != nil {
		t.Fatalf("unexpected error generating uuid: %v", err)
	}
}

func TestGenerateJWT(t *testing.T) {
	sub := uint(1)
	tokenString := genenerateJWT(sub)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.GetJWTSecret(), nil
	})
	if err != nil {
		t.Fatalf("unexpected error generating uuid: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	issuer, ok := claims["iss"].(string)
	if !ok || issuer != config.GetJWTIssuer() {
		t.Errorf("expected issuer = %v, got issuer = %v", config.GetJWTIssuer(), issuer)
	}

	subject, ok := claims["sub"]
	if !ok || subject != strconv.Itoa(int(sub)) {
		t.Errorf("expected sub = %v, got sub = %v", sub, sub)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		tokenType domain.TokenType
		tokenUse  domain.TokenUse
		wantErr   bool
		err       error
		userID    uint
	}{
		{
			name:      "uuid refresh token",
			tokenType: domain.UUID,
			tokenUse:  domain.REFRESH,
			userID:    1,
		},
		{
			name:      "unknown token type",
			tokenType: domain.TokenType("randomType"),
			tokenUse:  domain.REFRESH,
			userID:    1,
			wantErr:   true,
			err:       domain.ErrInvalidTokeType,
		},
		{
			name:      "unknown token use",
			tokenType: domain.UUID,
			tokenUse:  domain.TokenUse("randomUse"),
			userID:    1,
			wantErr:   true,
			err:       domain.ErrInvalidTokenUse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := ports.NewMockRepository(t)
			tsvc := NewTokenService(repo)
			_, err := tsvc.New(tt.tokenType, tt.tokenUse, tt.userID)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("unexpected error %v\n", err)
				}
				if !errors.Is(err, tt.err) {
					t.Fatalf("expected error: %v, got error: %v\n", tt.err, err)
				}
				return
			}
		})
	}
}
