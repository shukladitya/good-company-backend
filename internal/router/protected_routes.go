// internal/router/protected_routes.go
package router

import (
	"github.com/gin-gonic/gin"
)

func (r *Router) setupProtectedRoutes(router *gin.Engine) {
	// Protected routes group
	protected := router.Group("/api")
	protected.Use(r.authMiddleware.RequireAuth)
	{
		protected.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message":  "Access to protected resource",
				"user_id":  c.GetString("user_id"),
				"email":    c.GetString("email"),
				"username": c.GetString("username"),
			})
		})
	}
}