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
	"github.com/m1thrandir225/whoami/internal/util"
)

type EmailService interface {
	SendVerificationEmail(ctx context.Context, userID int64, email string) error
	VerifyEmailToken(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, userID int64, email string) error
}

type emailService struct {
	emailVerificationRepo repositories.EmailVerificationRepository
	userRepo              repositories.UserRepository
	config                *util.Config
}

func NewEmailService(
	emailVerificationRepo repositories.EmailVerificationRepository,
	userRepo repositories.UserRepository,
	config *util.Config,
) EmailService {
	return &emailService{
		emailVerificationRepo: emailVerificationRepo,
		userRepo:              userRepo,
		config:                config,
	}
}

func (s *emailService) SendVerificationEmail(ctx context.Context, userID int64, email string) error {
	// Generate verification token
	token, err := s.generateVerificationToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %v", err)
	}

	// Store verification token
	err = s.emailVerificationRepo.CreateEmailVerification(ctx, domain.CreateEmailVerificationAction{
		UserID: userID,
		Token:  token,
	})
	if err != nil {
		return fmt.Errorf("failed to store verification token: %v", err)
	}

	// Send email
	return s.sendEmail(email, "Email Verification", s.buildVerificationEmailContent(token))
}

func (s *emailService) VerifyEmailToken(ctx context.Context, token string) error {
	// Get verification record
	verification, err := s.emailVerificationRepo.GetEmailVerificationByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	// Check if token is expired
	if time.Now().After(verification.ExpiresAt) {
		return fmt.Errorf("verification token has expired")
	}

	// Check if already verified
	if verification.UsedAt != nil {
		return fmt.Errorf("email is already verified")
	}

	// Mark as verified
	err = s.emailVerificationRepo.MarkEmailVerified(ctx, verification.ID)
	if err != nil {
		return fmt.Errorf("failed to mark email as verified: %v", err)
	}

	// Update user email verification status
	user, err := s.userRepo.GetUserByID(ctx, verification.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	err = s.userRepo.MarkEmailVerified(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user verification status: %v", err)
	}

	return nil
}

func (s *emailService) ResendVerificationEmail(ctx context.Context, userID int64, email string) error {
	// Delete any existing unverified tokens for this user
	err := s.emailVerificationRepo.DeleteUnverifiedTokens(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to clean up existing tokens: %v", err)
	}

	// Send new verification email
	return s.SendVerificationEmail(ctx, userID, email)
}

func (s *emailService) generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *emailService) sendEmail(to, subject, body string) error {
	if s.config.SMTPHost == "" || s.config.SMTPPort == 0 {
		// In development, just log the email
		fmt.Printf("=== EMAIL SENT ===\nTo: %s\nSubject: %s\nBody:\n%s\n=== END EMAIL ===\n", to, subject, body)
		return nil
	}

	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.config.SMTPUsername, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	return smtp.SendMail(addr, auth, s.config.SMTPUsername, []string{to}, []byte(msg))
}

func (s *emailService) buildVerificationEmailContent(token string) string {
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", "", token)

	return fmt.Sprintf(`
Welcome to Whoami!

Please verify your email address by clicking the link below:

%s

This link will expire in 24 hours.

If you didn't create an account, you can safely ignore this email.

Best regards,
The Whoami Team
`, verificationURL)
}
