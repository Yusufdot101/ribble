package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Yusufdot101/ripple/shared/middleware/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/time/rate"
)

var CtxUserIDKey = "userID"

func RequireAuthentication(next gin.HandlerFunc) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		// read the token(jwt) from the request headers
		header := ctx.Request.Header.Get("Authorization")
		if len(header) == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrMissingInvalidToken.Error(),
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrMissingInvalidToken.Error(),
			})
			return
		}
		// validate it
		token, err := ValidateJWT(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidJWT.Error(),
			})
			return
		}

		// exctract the fields from it
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidJWT.Error(),
			})
			return
		}
		issuer, ok := claims["iss"].(string)
		if !ok || issuer != config.GetJWTIssuer() {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidJWT.Error(),
			})
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidJWT.Error(),
			})
			return
		}

		// add the userID to the request context
		ctx.Set(CtxUserIDKey, userID)
		next(ctx)
	}
	return fn
}

func RecoverPanic() gin.HandlerFunc {
	fn := func(c *gin.Context, err any) {
		message := "the server encountered and error and could not resolve your request"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": message,
		})
	}
	return gin.CustomRecovery(fn)
}

func IPRateLimiter() gin.HandlerFunc {
	limit := config.GetRateLimit()
	bucketSize := config.GetRateLimitBucketSize()
	type client struct {
		lastActive time.Time
		limiter    *rate.Limiter
		mu         *sync.Mutex
	}
	var (
		clients = make(map[string]*client)
		mu      sync.Mutex
	)
	// clean up
	go func() {
		for {
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastActive) > config.GetRateLimitResetTime() {
					delete(clients, ip)
				}
			}
			mu.Unlock()
			time.Sleep(3 * time.Minute)
		}
	}()

	fn := func(c *gin.Context) {
		ipAddr := c.ClientIP()
		mu.Lock()
		defer mu.Unlock()
		if _, exists := clients[ipAddr]; !exists {
			clients[ipAddr] = &client{
				limiter: rate.NewLimiter(rate.Limit(limit), bucketSize),
			}
		}
		clients[ipAddr].lastActive = time.Now()
		if !clients[ipAddr].limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": ErrTooManyRequests,
			})
			c.Abort()
			return
		}
		c.Next()
	}
	return fn
}
