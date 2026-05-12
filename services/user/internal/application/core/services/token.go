package services

import (
	"log"
	"strconv"
	"time"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenService struct {
	repo ports.Repository
}

func NewTokenService(repo ports.Repository) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

func (tsvc *TokenService) New(tokenType domain.TokenType, tokenUse domain.TokenUse, userID uint) (*domain.Token, error) {
	var ttl time.Duration
	switch tokenUse {
	case domain.REFRESH:
		ttl = config.GetRefreshTokenTTL()
	case domain.ACCESS:
		ttl = config.GetAccessTokenTTL()
	case domain.ACTIVATE:
		ttl = config.GetActivationTokenTTL()
	default:
		return nil, domain.ErrInvalidTokenUse
	}

	token := &domain.Token{
		TokenType: tokenType,
		UserID:    userID,
		Expires:   time.Now().Add(ttl),
		Use:       tokenUse,
	}
	var tokenString string
	switch tokenType {
	case domain.UUID:
		tokenString = generateUUID()
	case domain.JWT:
		tokenString = genenerateJWT(userID)
	default:
		return nil, domain.ErrInvalidTokeType
	}
	token.TokenString = tokenString
	return token, nil
}

func generateUUID() string {
	return uuid.New().String()
}

func genenerateJWT(userID uint) string {
	claims := jwt.RegisteredClaims{
		Issuer:    config.GetJWTIssuer(),
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.GetAccessTokenTTL())),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		log.Fatalf("error signing jwt: %v\n", err)
	}
	return tokenString
}

func (tsvc *TokenService) Save(token *domain.Token) error {
	return tsvc.repo.InsertToken(token)
}

func (tsvc *TokenService) GetTokenByStringAndUse(tokenString string, tokenUse domain.TokenUse) (*domain.Token, error) {
	return tsvc.repo.GetTokenByStringAndUse(tokenString, tokenUse)
}

func (tsvc *TokenService) RefreshAccessToken(refreshTokenString string) (string, error) {
	refreshToken, err := tsvc.GetTokenByStringAndUse(refreshTokenString, domain.REFRESH)
	if err != nil {
		return "", err
	}
	accessToken := genenerateJWT(refreshToken.UserID)
	return accessToken, nil
}

func (tsvc *TokenService) DeleteTokenByStringAndUse(tokenString string, tokenUse domain.TokenUse) error {
	return tsvc.repo.DeleteTokenByStringAndUse(tokenString, tokenUse)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// ensure the token was signed with HMAC, not something else
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidJWTSigningMethod
		}
		return config.GetJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidJWT
	}

	return token, nil
}
