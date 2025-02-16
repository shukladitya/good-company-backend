// internal/services/auth_service.go
package services

import (
	"errors"
	"fmt"
	"os"
	"theGoodCompany/internal/models"

	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

// Credentials struct for login
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateUser registers a new user
func (s *AuthService) CreateUser(user *models.User) error {
	// Validate input
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return errors.New("username, email, and password are required")
	}

	// Hash the password
	if err := user.HashPassword(); err != nil {
		return err
	}

	// Create user in database
	result := s.DB.Create(user)
	return result.Error
}

// Authenticate validates user credentials and returns a JWT token
func (s *AuthService) Authenticate(creds Credentials) (string, error) {
	var user models.User
	
	// Find user by username
	result := s.DB.Where("username = ?", creds.Username).First(&user)
	
	if result.Error != nil {
		return "", errors.New("invalid username or password")
	}

	// Check password
	if !user.CheckPasswordHash(creds.Password) {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT token
	return s.GenerateJWTToken(user)
}

// generateJWTToken creates a new JWT token
func (s *AuthService) GenerateJWTToken(user models.User) (string, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT secret not configured")
	}

	
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   strconv.Itoa(int(user.ID)),
		"username":  user.Username,
		"email":     user.Email,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // Token valid for 7 days
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWTToken checks if a token is valid
func (s *AuthService) ValidateJWTToken(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT secret not configured")
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
}


// VerifyEmail verifies the email using the provided verification code
func (s *AuthService) VerifyEmail(verificationCode string) (*models.User, error) {
	if verificationCode == "" {
		return nil, errors.New("verification code is required")
	}

	// Find the user by the verification code
	var user models.User
	result := s.DB.Where("email_verification_code = ?", verificationCode).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid verification code")
		}
		return nil, result.Error
	}

	// Update user's email verification status
	user.IsVerified = true
	user.VerificationCode = "" // Clear the code

	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
