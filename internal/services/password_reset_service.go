package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/repositories"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/util"
)

type PasswordResetService interface {
	RequestPasswordReset(ctx context.Context, email string) error
	VerifyResetToken(ctx context.Context, token string) (*domain.PasswordReset, error)
	ResetPassword(ctx context.Context, token string, newPassword string) error
}

type passwordResetService struct {
	passwordResetRepo       repositories.PasswordResetRepository
	userRepo                repositories.UserRepository
	passwordSecurityService PasswordSecurityService
	config                  *util.Config
}

func NewPasswordResetService(
	passwordResetRepo repositories.PasswordResetRepository,
	userRepo repositories.UserRepository,
	passwordSecurityService PasswordSecurityService,
	config *util.Config,
) PasswordResetService {
	return &passwordResetService{
		passwordResetRepo:       passwordResetRepo,
		userRepo:                userRepo,
		passwordSecurityService: passwordSecurityService,
		config:                  config,
	}
}
func (s *passwordResetService) RequestPasswordReset(ctx context.Context, email string) error {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists or not
		return nil
	}

	// Generate HOTP secret and counter
	hotpSecret, err := security.GenerateHOTPSecret()
	if err != nil {
		return fmt.Errorf("failed to generate HOTP secret: %v", err)
	}

	// Generate token
	token, err := s.generateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %v", err)
	}
	// Delete any existing unused resets for this user
	err = s.passwordResetRepo.DeleteUnusedPasswordResets(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to clean up existing resets: %v", err)
	}

	// Create password reset record
	err = s.passwordResetRepo.CreatePasswordReset(ctx, domain.CreatePasswordResetAction{
		UserID:     user.ID,
		TokenHash:  token,
		HotpSecret: hotpSecret,
		Counter:    0,
	})
	if err != nil {
		return fmt.Errorf("failed to create password reset: %v", err)
	}

	// Send password reset email
	return s.sendPasswordResetEmail(user.Email, token)
}

func (s *passwordResetService) VerifyResetToken(ctx context.Context, token string) (*domain.PasswordReset, error) {
	// Get reset record
	reset, err := s.passwordResetRepo.GetPasswordResetByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired reset token")
	}

	// Check if token is expired
	if time.Now().After(reset.ExpiresAt) {
		return nil, fmt.Errorf("reset token has expired")
	}

	// Check if already used
	if reset.UsedAt != nil {
		return nil, fmt.Errorf("reset token has already been used")
	}

	return reset, nil
}

func (s *passwordResetService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// Verify the token
	reset, err := s.VerifyResetToken(ctx, token)
	if err != nil {
		return err
	}

	// Validate the new password
	if err := s.passwordSecurityService.ValidateNewUserPassword(ctx, newPassword); err != nil {
		return fmt.Errorf("invalid password: %v", err)
	}

	// Update the user's password
	if err := s.passwordSecurityService.UpdatePassword(ctx, reset.UserID, newPassword); err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	// Mark the reset token as used
	if err := s.passwordResetRepo.MarkPasswordResetAsUsed(ctx, reset.ID); err != nil {
		return fmt.Errorf("failed to mark reset as used: %v", err)
	}

	return nil
}

func (s *passwordResetService) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *passwordResetService) sendPasswordResetEmail(email, token string) error {
	if s.config.SMTPHost == "" || s.config.SMTPPort == 0 {
		// In development, just log the email
		fmt.Printf("=== PASSWORD RESET EMAIL ===\nTo: %s\nSubject: Password Reset\nBody:\n%s\n=== END EMAIL ===\n",
			email, s.buildPasswordResetEmailContent(token))
		return nil
	}

	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.config.SMTPUsername, email, "Password Reset", s.buildPasswordResetEmailContent(token))

	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	return smtp.SendMail(addr, auth, s.config.SMTPUsername, []string{email}, []byte(msg))
}

func (s *passwordResetService) buildPasswordResetEmailContent(token string) string {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", "", token)

	return fmt.Sprintf(`
Password Reset Request

You requested a password reset. Click the link below to reset your password:

%s

This link will expire in 1 hour.

If you didn't request a password reset, you can safely ignore this email.

Best regards,
The Whoami Team
`, resetURL)
}
