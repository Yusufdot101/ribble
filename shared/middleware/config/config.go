package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	home, _ := os.UserHomeDir()
	_ = godotenv.Load(fmt.Sprintf("%s/Documents/projects/ripple/shared/middleware/config/.env", home))
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

func GetRateLimit() float64 {
	limit, err := strconv.ParseFloat(getEnvVariable("RATE_LIMIT"), 64)
	if err != nil {
		log.Fatalf("RATE_LIMIT variable not found: %v", err)
	}

	if limit <= 0 {
		log.Fatal("RATE_LIMIT must be > 0")
	}
	return limit
}

func GetRateLimitBucketSize() int {
	bucketSize, err := strconv.Atoi(getEnvVariable("RATE_LIMIT_BUCKET_SIZE"))
	if err != nil {
		panic(fmt.Sprintf("error: bucket size not found: %v", err))
	}
	if bucketSize <= 0 {
		log.Fatal("RATE_LIMIT_BUCKET_SIZE must be > 0")
	}
	return bucketSize
}

// GetRateLimitResetTime is the duration a client's rate limit is kept
func GetRateLimitResetTime() time.Duration {
	duration, err := time.ParseDuration(getEnvVariable("RATE_LIMIT_RESET_TIME"))
	if err != nil {
		log.Fatalf("invalid RATE_LIMIT_RESET_TIME: %v", err)
	}
	if duration <= 0 {
		log.Fatal("RATE_LIMIT_RESET_TIME must be > 0")
	}

	return duration
}

func GetRateLimitCredentials() float64 {
	limit, err := strconv.ParseFloat(getEnvVariable("RATE_LIMIT_CREDENTIALS"), 64)
	if err != nil {
		log.Fatalf("RATE_LIMIT variable not found: %v", err)
	}

	if limit <= 0 {
		log.Fatal("RATE_LIMIT must be > 0")
	}
	return limit
}

func GetRateLimitBucketSizeCredentials() int {
	bucketSize, err := strconv.Atoi(getEnvVariable("RATE_LIMIT_BUCKET_SIZE_CREDENTIALS"))
	if err != nil {
		panic(fmt.Sprintf("error: bucket size not found: %v", err))
	}
	if bucketSize <= 0 {
		log.Fatal("RATE_LIMIT_BUCKET_SIZE must be > 0")
	}
	return bucketSize
}

func GetRateLimitResetTimeCredentials() time.Duration {
	duration, err := time.ParseDuration(getEnvVariable("RATE_LIMIT_RESET_TIME_CREDENTIALS"))
	if err != nil {
		log.Fatalf("invalid RATE_LIMIT_RESET_TIME: %v", err)
	}
	if duration <= 0 {
		log.Fatal("RATE_LIMIT_RESET_TIME must be > 0")
	}

	return duration
}

func getEnvVariable(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env var %s missing\n", key)
	}
	return val
}
