package handlers

import (
	"net/http"
	"theGoodCompany/internal/services"

	"time"

	"github.com/gin-gonic/gin"
)

type PasswordResetHandler struct {
	PasswordResetService *services.PasswordResetService
	EmailService         *services.EmailService
	// DB 					 *gorm.DB   // this is incorrect, DB is not injected in handler in main, its only inside services
}

type Credentials struct {
	Email    string `json:"email"`
}

type NewPassword struct {
	NewPassword string `json:"new_password"`
}

func (h *PasswordResetHandler) GeneratePasswordResetLink(c *gin.Context) {

	var Credentials Credentials

	err := c.BindJSON(&Credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	user,err := h.PasswordResetService.FindUserByEmail(Credentials.Email)
	if(err != nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}
	
	//create hash
	user.ResetPasswordCode = user.GenerateVerificationCode()
	user.ResetPasswordAt = time.Now()

	err = h.EmailService.SendPasswordResetEmail(user.Email, user.ResetPasswordCode)

	if(err != nil){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to send email",
			"message": err.Error(),
		})
		return
	}

	//update user information with passwordCode and resetPasswordAt
	if err := h.PasswordResetService.SaveUserWithResetPasswordCode(user, user.ResetPasswordCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save user information",
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset link sent successfully",
	})

}

func (h *PasswordResetHandler) ResetPasswordLinkValidation(c *gin.Context) {
	// extract reset code from url
	passwordResetCode := c.Query("code")

	// validate reset code
	 user, err := h.PasswordResetService.FindUserByResetPasswordCode(passwordResetCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid reset code",
			"message": err.Error(),
		})
		return
	}

	// check if link was send more than 1 day before
	if time.Since(user.ResetPasswordAt).Hours() > 24 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Link expired",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Link is valid",
		"user_id": user.ID,
	})
}

func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	// get new password from request body
	passwordResetCode := c.Param("code")

	var newPassword NewPassword

	if err := c.BindJSON(&newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	user, err := h.PasswordResetService.FindUserByResetPasswordCode(passwordResetCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid reset code",
		})
		return
	}

	user.Password = newPassword.NewPassword
	user.ResetPasswordCode = ""
	user.ResetPasswordAt = time.Time{}

	user.HashPassword()

	if err := h.PasswordResetService.SaveUserWithNewPasswordHash(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save user information",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}	