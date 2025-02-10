package handlers

import (
	"log"
	"net/http"
	"serveMovies/internal/models"
	"serveMovies/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related operations
type AuthHandler struct {
	AuthService  *services.AuthService
	EmailService *services.EmailService
}


// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User

	// Bind JSON body to user struct
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Generate verification code
	user.VerificationCode = user.GenerateVerificationCode()

	// Create user
	if err := h.AuthService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send verification email
	go func() {
		if err := h.EmailService.SendVerificationEmail(user.Email, user.VerificationCode); err != nil {
			// Log error, but don't block response
			log.Printf("Failed to send verification email: %v", err)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please check your email to verify your account.",
		"user_id": user.ID,
	})
}

// VerifyEmail handles email verification
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	verificationCode := c.Query("code")

	// Verify email
	user, err := h.AuthService.VerifyEmail(verificationCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or expired verification code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
		"user_id": user.ID,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var creds services.Credentials

	// Bind JSON body to credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Authenticate user
	token, err := h.AuthService.Authenticate(creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Login successful",
	})
}
