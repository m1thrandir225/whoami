package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/repositories"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/util"
)

type PasswordSecurityService interface {
	ValidatePassword(ctx context.Context, userID int64, newPassword string) error
	UpdatePassword(ctx context.Context, userID int64, newPassword string) error
	ValidateNewUserPassword(ctx context.Context, newPassword string) error
	CheckPasswordStrength(password string) error
	AddInitialPasswordToHistory(ctx context.Context, userID int64, passwordHash string) error
}

type passwordSecurityService struct {
	passwordHistoryRepo repositories.PasswordHistoryRepository
	userRepo            repositories.UserRepository
	hibpClient          *security.HaveIBeenPwnedClient
}

func NewPasswordSecurityService(
	passwordHistoryRepo repositories.PasswordHistoryRepository,
	userRepo repositories.UserRepository,
) PasswordSecurityService {
	return &passwordSecurityService{
		passwordHistoryRepo: passwordHistoryRepo,
		userRepo:            userRepo,
		hibpClient:          security.NewHaveIBeenPwnedClient(),
	}
}

func (s *passwordSecurityService) ValidatePassword(ctx context.Context, userID int64, newPassword string) error {
	// Check password strength
	if err := s.CheckPasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash the password to check against history
	passwordHash, err := util.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Check if password is in history (prevent reuse of last 5 passwords)
	inHistory, err := s.passwordHistoryRepo.CheckPasswordInHistory(ctx, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to check password history: %v", err)
	}

	if inHistory {
		return errors.New("password has been used recently and cannot be reused")
	}

	// Check if password has been compromised
	pwnedPassword, err := s.hibpClient.CheckPasswordWithRetry(newPassword, 3)
	if err != nil {
		// If we can't check, log the error but don't block the password change
		// In production, you might want to be more strict about this
		fmt.Printf("Warning: Could not check password against HaveIBeenPwned: %v\n", err)
		return nil
	}

	if pwnedPassword.IsPwned {
		return fmt.Errorf("password has been compromised %d times and should not be used", pwnedPassword.Count)
	}

	return nil
}

func (s *passwordSecurityService) ValidateNewUserPassword(ctx context.Context, newPassword string) error {
	// Check password strength
	if err := s.CheckPasswordStrength(newPassword); err != nil {
		return err
	}

	// Check if password has been compromised (for new users, we don't check history)
	pwnedPassword, err := s.hibpClient.CheckPasswordWithRetry(newPassword, 3)
	if err != nil {
		// If we can't check, log the error but don't block registration
		fmt.Printf("Warning: Could not check password against HaveIBeenPwned: %v\n", err)
		return nil
	}

	if pwnedPassword.IsPwned {
		return fmt.Errorf("password has been compromised %d times and should not be used", pwnedPassword.Count)
	}

	return nil
}

func (s *passwordSecurityService) UpdatePassword(ctx context.Context, userID int64, newPassword string) error {
	// Validate the password first
	if err := s.ValidatePassword(ctx, userID, newPassword); err != nil {
		return err
	}

	// Hash the password
	passwordHash, err := util.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Get current user to update
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	// Update user password
	user.Password = passwordHash
	now := time.Now()
	user.PasswordChangedAt = &now

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user password: %v", err)
	}

	// Add to password history
	req := domain.CreatePasswordHistory{
		UserID:       userID,
		PasswordHash: passwordHash,
	}
	if err := s.passwordHistoryRepo.CreatePasswordHistory(ctx, req); err != nil {
		return fmt.Errorf("failed to add password to history: %v", err)
	}

	// Clean up old password history (keep only last 5 passwords)
	if err := s.passwordHistoryRepo.DeleteOldPasswordHistory(ctx, userID); err != nil {
		// Log but don't fail the password update
		fmt.Printf("Warning: Failed to clean up old password history: %v\n", err)
	}

	return nil
}

func (s *passwordSecurityService) AddInitialPasswordToHistory(ctx context.Context, userID int64, passwordHash string) error {
	req := domain.CreatePasswordHistory{
		UserID:       userID,
		PasswordHash: passwordHash,
	}
	return s.passwordHistoryRepo.CreatePasswordHistory(ctx, req)
}

func (s *passwordSecurityService) CheckPasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return errors.New("password must be less than 128 characters")
	}

	commonPasswords := []string{
		"password", "123456", "123456789", "qwerty", "abc123",
		"password123", "admin", "letmein", "welcome", "monkey",
	}

	lowerPassword := strings.ToLower(password)
	for _, common := range commonPasswords {
		if lowerPassword == common {
			return errors.New("password is too common and easily guessable")
		}
	}

	// Check for character variety
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("password must contain uppercase, lowercase, and numeric characters")
	}

	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
