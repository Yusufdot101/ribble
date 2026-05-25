package middleware

import (
	"net/http"
	"strings"

	"github.com/Yusufdot101/ripple/shared/middleware/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var CtxUserIDKey = "userID"

func RequireAuthentication() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// read the token(jwt) from the request headers
		header := c.Request.Header.Get("Authorization")
		if len(header) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrMissingInvalidToken.Error(),
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrMissingInvalidToken.Error(),
			})
			return
		}
		// validate it
		token, err := ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidJWT.Error(),
			})
			return
		}

		// exctract the fields from it
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidJWT.Error(),
			})
			return
		}
		issuer, ok := claims["iss"].(string)
		if !ok || issuer != config.GetJWTIssuer() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidJWT.Error(),
			})
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidJWT.Error(),
			})
			return
		}

		// add the userID to the request context
		c.Set(CtxUserIDKey, userID)
		c.Next()
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
