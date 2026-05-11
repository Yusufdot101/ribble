package config

import (
	"log"
	"os"
	"time"
)

func GetDatabaseURL() string {
	return getEnvVariable("DATABASE_URL")
}

func GetGoogleClientID() string {
	return getEnvVariable("GOOGLE_CLIENT_ID")
}

func GetGoogleClientSecret() string {
	return getEnvVariable("GOOGLE_CLIENT_SECRET")
}

func GetRefreshTokenTTL() time.Duration {
	duration, err := time.ParseDuration(getEnvVariable("REFRESH_TOKEN_TTL"))
	if err != nil {
		log.Fatalf("invalid refresh token ttl")
	}

	return duration
}

func RefreshTokenIsSecure() bool {
	return getEnvVariable("REFRESH_TOKEN_COOKIE_SECURE") != "false" // default true
}

func GetFrontendURL() string {
	return getEnvVariable("FRONTEND_URL")
}

func GetJWTIssuer() string {
	return getEnvVariable("JWT_ISSUER")
}

func GetAccessTokenTTL() time.Duration {
	duration, err := time.ParseDuration(getEnvVariable("ACCESS_TOKEN_TTL"))
	if err != nil {
		log.Fatalf("invalid access token ttl")
	}

	return duration
}

func GetEmailVerificationTTL() time.Duration {
	duration, err := time.ParseDuration(getEnvVariable("EMAIL_VERIFICATION_TTL"))
	if err != nil {
		log.Fatalf("invalid email verification ttl")
	}

	return duration
}

func GetJWTSecret() []byte {
	jwtSecret := getEnvVariable("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET variable must be at least 32 bytes for HS256")
	}

	return []byte(jwtSecret)
}

func GetEnv() string {
	return getEnvVariable("ENV")
}

func getEnvVariable(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env var %s missing\n", key)
	}
	return val
}
