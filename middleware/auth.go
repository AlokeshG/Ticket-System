package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ticket-system/utils"
)

// AuthMiddleware protects routes that require authentication.
func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Get Authorization header.
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			c.Abort()
			return
		}

		// Expected format:
		//
		// Authorization: Bearer <token>
		//
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token is required",
			})
			c.Abort()
			return
		}

		// Validate JWT.
		userID, err := utils.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store user ID in Gin context.
		c.Set("user_id", userID)

		// Continue to the protected handler.
		c.Next()
	}
}
