// internal/router/public_routes.go
package router

import (
	"github.com/gin-gonic/gin"
)

func (r *Router) setupPublicRoutes(router *gin.Engine) {
	// auth routes
	router.POST("/register", r.authHandler.Register)
	router.POST("/login", r.authHandler.Login)
	router.GET("/verify", r.authHandler.VerifyEmail)
	router.POST("/reset-password", r.passwordResetHandler.GeneratePasswordResetLink)
	router.GET("/reset-password", r.passwordResetHandler.ResetPasswordLinkValidation)
	router.POST("/reset-password/:code", r.passwordResetHandler.ResetPassword)
}