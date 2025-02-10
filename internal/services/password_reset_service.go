package services

import (
	"serveMovies/internal/models"

	"gorm.io/gorm"
)

type PasswordResetService struct {
	DB *gorm.DB
}

func (s *PasswordResetService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *PasswordResetService) FindUserByResetPasswordCode(resetPasswordCode string) (*models.User, error) {
	var user models.User
	result := s.DB.Where("reset_password_code = ?", resetPasswordCode).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *PasswordResetService) SaveUserWithResetPasswordCode(user *models.User, resetPasswordCode string) error {
	user.ResetPasswordCode = resetPasswordCode
	return s.DB.Save(user).Error
}

func (s *PasswordResetService) SaveUserWithNewPasswordHash(user *models.User) error {
	return s.DB.Save(user).Error
}