package config

import (
	"log"
	"os"
)

func GetDatabaseURL() string {
	return getEnvVariable("DATABASE_URL")
}

func GetJWTIssuer() string {
	return getEnvVariable("JWT_ISSUER")
}

func GetJWTSecret() []byte {
	jwtSecret := getEnvVariable("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET variable must be at least 32 bytes for HS256")
	}

	return []byte(jwtSecret)
}

func GetFrontendURL() string {
	return getEnvVariable("FRONTEND_URL")
}

func GetGRPCUrl() string {
	return getEnvVariable("GRPC_URL")
}

func getEnvVariable(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env var %s missing\n", key)
	}
	return val
}
