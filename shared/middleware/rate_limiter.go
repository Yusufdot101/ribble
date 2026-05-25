package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Yusufdot101/ripple/shared/middleware/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func IPRateLimiter() gin.HandlerFunc {
	limit := config.GetRateLimit()
	bucketSize := config.GetRateLimitBucketSize()
	type client struct {
		lastActive time.Time
		limiter    *rate.Limiter
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
		if _, exists := clients[ipAddr]; !exists {
			clients[ipAddr] = &client{
				limiter: rate.NewLimiter(rate.Limit(limit), bucketSize),
			}
		}
		clients[ipAddr].lastActive = time.Now()
		allowed := clients[ipAddr].limiter.Allow()
		mu.Unlock()
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": ErrTooManyRequests.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
	return fn
}

func IdentityRateLimiter() gin.HandlerFunc {
	limit := 1
	bucketSize := 100
	type client struct {
		lastActive time.Time
		limiter    *rate.Limiter
	}
	var (
		clients = make(map[uint]*client)
		mu      sync.Mutex
	)
	// clean up
	go func() {
		for {
			mu.Lock()
			for id, client := range clients {
				if time.Since(client.lastActive) > time.Hour {
					delete(clients, id)
				}
			}
			mu.Unlock()
			time.Sleep(3 * time.Minute)
		}
	}()

	fn := func(c *gin.Context) {
		userID := userIDFromContext(c)
		mu.Lock()
		if _, exists := clients[userID]; !exists {
			clients[userID] = &client{
				limiter: rate.NewLimiter(rate.Limit(limit), bucketSize),
			}
		}
		clients[userID].lastActive = time.Now()
		allowed := clients[userID].limiter.Allow()
		mu.Unlock()
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
	return fn
}

func CredentialRateLimiter() gin.HandlerFunc {
	limit := config.GetRateLimit()
	bucketSize := config.GetRateLimitBucketSize()
	type client struct {
		lastActive time.Time
		limiter    *rate.Limiter
	}
	var (
		clients = make(map[string]*client)
		mu      sync.Mutex
	)
	// clean up
	go func() {
		for {
			mu.Lock()
			for email, client := range clients {
				if time.Since(client.lastActive) > config.GetRateLimitResetTime() {
					delete(clients, email)
				}
			}
			mu.Unlock()
			time.Sleep(3 * time.Minute)
		}
	}()

	fn := func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			c.Abort()
			return
		}
		mu.Lock()
		if _, exists := clients[req.Email]; !exists {
			clients[req.Email] = &client{
				limiter: rate.NewLimiter(rate.Limit(limit), bucketSize),
			}
		}
		clients[req.Email].lastActive = time.Now()
		allowed := clients[req.Email].limiter.Allow()
		mu.Unlock()
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": ErrTooManyRequests.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
	return fn
}

func userIDFromContext(ctx *gin.Context) uint {
	currentUserID, ok := ctx.MustGet("userID").(string)
	if !ok {
		panic("user id missing")
	}

	currentUserIDint, err := strconv.ParseUint(currentUserID, 10, 32)
	if err != nil {
		panic("invalid user id type")
	}
	return uint(currentUserIDint)
}
