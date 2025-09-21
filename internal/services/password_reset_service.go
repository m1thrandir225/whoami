package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/mail"
	"github.com/m1thrandir225/whoami/internal/repositories"
	"github.com/m1thrandir225/whoami/internal/security"
)

type PasswordResetService interface {
	RequestPasswordReset(ctx context.Context, email string) error
	VerifyResetToken(ctx context.Context, token string) (*domain.PasswordReset, error)
	VerifyResetOTP(ctx context.Context, token, otp string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
}

type passwordResetService struct {
	passwordResetRepo       repositories.PasswordResetRepository
	userRepo                repositories.UserRepository
	passwordSecurityService PasswordSecurityService
	mailService             mail.MailService
	frontendURL             string
}

func NewPasswordResetService(
	passwordResetRepo repositories.PasswordResetRepository,
	userRepo repositories.UserRepository,
	passwordSecurityService PasswordSecurityService,
	mailService mail.MailService,
	frontendURL string,
) PasswordResetService {
	return &passwordResetService{
		passwordResetRepo:       passwordResetRepo,
		userRepo:                userRepo,
		passwordSecurityService: passwordSecurityService,
		mailService:             mailService,
		frontendURL:             frontendURL,
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

	hotp, err := security.GenerateHOTP(reset.HotpSecret, uint64(reset.Counter))
	if err != nil {
		return nil, fmt.Errorf("failed to generate HOTP: %v", err)
	}

	user, err := s.userRepo.GetUserByID(ctx, reset.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	if err := s.sendOTPEmail(user.Email, hotp); err != nil {
		return nil, fmt.Errorf("failed to send OTP mail: %v", err)
	}

	return reset, nil
}

func (s *passwordResetService) VerifyResetOTP(ctx context.Context, token, otp string) error {
	// Get reset record
	reset, err := s.passwordResetRepo.GetPasswordResetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Check if token is expired
	if time.Now().After(reset.ExpiresAt) {
		return fmt.Errorf("reset token has expired")
	}

	// Check if already used
	if reset.UsedAt != nil {
		return fmt.Errorf("reset token has already been used")
	}

	// Verify OTP using HOTP
	hotp, err := security.ValidateHOTP(reset.HotpSecret, otp, uint64(reset.Counter))
	if err != nil {
		return fmt.Errorf("failed to validate HOTP: %v", err)
	}

	if !hotp {
		return fmt.Errorf("invalid OTP")
	}

	if err := s.passwordResetRepo.IncrementPasswordResetCounter(ctx, reset.ID); err != nil {
		return fmt.Errorf("failed to increment counter: %v", err)
	}

	return nil
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
	return s.mailService.SendMail("whoami@sebastijanzindl.me", email, "Password Reset", s.buildPasswordResetEmailContent(token))
}

func (s *passwordResetService) sendOTPEmail(email, otp string) error {
	return s.mailService.SendMail("whoami@sebastijanzindl.me", email, "Password Reset Verification Code", s.buildOTPEmailContent(otp))
}

func (s *passwordResetService) buildPasswordResetEmailContent(token string) string {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, token)

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

func (s *passwordResetService) buildOTPEmailContent(otp string) string {
	return fmt.Sprintf(`
Password Reset Verification Code

Your verification code is: %s

Enter this code to continue with your password reset.

This code will expire in 10 minutes.

If you didn't request a password reset, you can safely ignore this email.

Best regards,
The Whoami Team
`, otp)
}
