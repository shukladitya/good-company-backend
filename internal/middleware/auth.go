// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"serveMovies/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	AuthService *services.AuthService
}

// RequireAuth is a middleware to protect routes
func (m *AuthMiddleware) RequireAuth(c *gin.Context) {
	// Get Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token required",
		})
		c.Abort()
		return
	}

	// Extract token (expecting "Bearer <token>")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token format",
		})
		c.Abort()
		return
	}

	token := parts[1]

	// Validate token
	jwtToken, err := m.AuthService.ValidateJWTToken(token)
	if err != nil || !jwtToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired token",
		})
		c.Abort()
		return
	}

	// Optional: Extract and set user claims in context
	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if ok {
		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Set("username", claims["username"])
	}

	// Continue to the next middleware/handler
	c.Next()
}