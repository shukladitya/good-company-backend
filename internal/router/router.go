// internal/router/router.go
package router

import (
	"serveMovies/internal/handlers"
	"serveMovies/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler         *handlers.AuthHandler
	passwordResetHandler *handlers.PasswordResetHandler
	authMiddleware      *middleware.AuthMiddleware
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	passwordResetHandler *handlers.PasswordResetHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		authHandler:         authHandler,
		passwordResetHandler: passwordResetHandler,
		authMiddleware:      authMiddleware,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Healthy",
		})
	})

	// Setup routes
	r.setupPublicRoutes(router)
	r.setupProtectedRoutes(router)

	return router
}