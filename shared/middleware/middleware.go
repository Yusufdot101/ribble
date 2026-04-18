package middleware

import (
	"net/http"
	"strings"

	"github.com/Yusufdot101/ripple/shared/middleware/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
