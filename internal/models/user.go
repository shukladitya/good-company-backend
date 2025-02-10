package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username     string    `gorm:"not null" json:"username"`
    Email        string    `gorm:"unique;not null" json:"email"`
    PasswordHash string    `gorm:"not null" json:"-"` // Stored in DB but excluded from JSON
    Password     string    `gorm:"-" json:"password,omitempty"` // Ignored by GORM, used only for input
    LastLogin    time.Time `json:"last_login"`

	// Verification fields
	IsVerified       bool      `json:"is_verified"`
	VerificationCode string    `json:"verification_code"`
	VerifiedAt       time.Time `json:"verified_at"`

	// Reset password fields
	ResetPasswordCode string `json:"reset_password_code"`
	ResetPasswordAt   time.Time `json:"reset_password_at"`
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	if u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	// Generate a bcrypt hash
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	
	u.PasswordHash = string(bytes)

	return nil
}

func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// GenerateVerificationCode creates a unique verification code
func (u *User) GenerateVerificationCode() string {
	return uuid.New().String()
}