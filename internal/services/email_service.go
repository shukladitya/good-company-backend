// internal/services/email_service.go
package services

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	SMTPHost     string
	SMTPPort     int
	SenderEmail  string
	SenderName   string
	SenderPass   string
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    587, // Default port, can be changed
		SenderEmail: os.Getenv("SMTP_EMAIL"),
		SenderName:  "Streaming Service",
		SenderPass:  os.Getenv("SMTP_PASSWORD"),
	}
}

// SendVerificationEmail sends an email with verification link
func (s *EmailService) SendVerificationEmail(toEmail, verificationCode string) error {

	// Construct verification link
	verificationLink := fmt.Sprintf(
		"%s/verify?code=%s", 
		os.Getenv("APP_FRONTEND_URL"), 
		verificationCode,
	)

	// Create message
	m := gomail.NewMessage()

	m.SetHeader("From", s.SenderEmail)
	
	m.SetHeader("To", toEmail)
	
	m.SetHeader("Subject", "Verify Your Account")
	
	m.SetBody("text/html", fmt.Sprintf(`
		<h1>Verify Your Account</h1>
		<p>Click the link below to verify your account:</p>
		<a href="%s">Verify Account</a>
		<p>If you didn't create an account, please ignore this email.</p>
	`, verificationLink))


	// Send email
	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.SenderEmail, s.SenderPass)

	return d.DialAndSend(m)
}

func (s *EmailService) SendPasswordResetEmail(toEmail, resetCode string) error {

	// Construct verification link
	resetLink := fmt.Sprintf(
		"%s/reset-password?code=%s", 
		os.Getenv("APP_FRONTEND_URL"), 
		resetCode,
	)

	// Create message
	m := gomail.NewMessage()

	m.SetHeader("From", s.SenderEmail)
	
	m.SetHeader("To", toEmail)
	
	m.SetHeader("Subject", "Reset Your Password")
	
	m.SetBody("text/html", fmt.Sprintf(`
		<h1>Reset Your Password</h1>
		<p>Click the link below to reset your password:</p>
		<a href="%s">Reset Password</a>
		<p>If you didn't request a password reset, please ignore this email.</p>
	`, resetLink))


	// Send email
	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.SenderEmail, s.SenderPass)

	return d.DialAndSend(m)
}