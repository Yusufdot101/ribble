package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	home, _ := os.UserHomeDir()
	_ = godotenv.Load(fmt.Sprintf("%s/Documents/projects/ribble/shared/middleware/config/.env", home))
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

func getEnvVariable(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env var %s missing\n", key)
	}
	return val
}
