package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/repositories"
)

type SecurityService interface {
	RecordFailedLogin(ctx context.Context, userID int64, email, ipAddress, userAgent string) error
	RecordSuccessfulLogin(ctx context.Context, userID int64, email, ipAddress, userAgent string) error
	CheckAccountLockout(ctx context.Context, userID int64, ipAddress string) error
	RecordSuspiciousActivity(ctx context.Context, req domain.CreateSuspiciousActivityAction) error
	GetSuspiciousActivities(ctx context.Context, userID int64) ([]domain.SuspiciousActivity, error)
	ResolveSuspiciousActivity(ctx context.Context, activityID int64) error
	CleanupExpiredLockouts(ctx context.Context) error
}

type securityService struct {
	loginAttemptsRepo repositories.LoginAttemptsRepository
	suspiciousRepo    repositories.SuspiciousActivityRepository
	lockoutRepo       repositories.AccountLockoutRepository
	userRepo          repositories.UserRepository
}

func NewSecurityService(
	loginAttemptsRepo repositories.LoginAttemptsRepository,
	suspiciousRepo repositories.SuspiciousActivityRepository,
	lockoutRepo repositories.AccountLockoutRepository,
	userRepo repositories.UserRepository,
) SecurityService {
	return &securityService{
		loginAttemptsRepo: loginAttemptsRepo,
		suspiciousRepo:    suspiciousRepo,
		lockoutRepo:       lockoutRepo,
		userRepo:          userRepo,
	}
}

// RecordFailedLogin records a failed login attempt and checks if the account should be locked
func (s *securityService) RecordFailedLogin(ctx context.Context, userID int64, email, ipAddress, userAgent string) error {
	// 1. Record in login_attempts table (for detailed tracking)
	failureReason := "invalid_credentials"
	_, err := s.loginAttemptsRepo.CreateLoginAttempt(ctx, domain.CreateLoginAttemptAction{
		UserID:        &userID,
		Email:         email,
		IPAddress:     ipAddress,
		UserAgent:     &userAgent,
		Success:       false,
		FailureReason: &failureReason,
	})
	if err != nil {
		return err
	}

	// 2. Check if we should create suspicious activity
	var recentAttempts []domain.LoginAttempt
	if userID > 0 {
		recentAttempts, err = s.loginAttemptsRepo.GetRecentFailedAttemptsByUserID(ctx, userID)
	} else {
		recentAttempts, err = s.loginAttemptsRepo.GetRecentFailedAttemptsByEmail(ctx, email)
	}
	if err != nil {
		return err
	}

	// 3. Record suspicious activity if multiple failed attempts
	if len(recentAttempts) >= 3 {
		metadata, _ := json.Marshal(map[string]interface{}{
			"action":          "multiple_failed_logins",
			"failed_attempts": len(recentAttempts),
			"timestamp":       time.Now().Unix(),
		})

		mediumSeverity := domain.MediumActivity
		_, err = s.suspiciousRepo.CreateActivity(ctx, domain.CreateSuspiciousActivityAction{
			UserID:       userID,
			ActivityType: "multiple_failed_logins",
			IPAddress:    ipAddress,
			UserAgent:    userAgent,
			Description:  fmt.Sprintf("Multiple failed login attempts (%d)", len(recentAttempts)),
			Metadata:     metadata,
			Severity:     &mediumSeverity,
		})
	}

	// 4. Lock account if too many failed attempts
	if len(recentAttempts) >= 5 {
		lockoutDuration := 30 * time.Minute
		expiresAt := time.Now().Add(lockoutDuration)

		_, err = s.lockoutRepo.CreateLockout(ctx, domain.CreateAccountLockoutAction{
			UserID:      userID,
			IPAddress:   ipAddress,
			LockoutType: domain.LockoutTypeAccount,
			ExpiresAt:   expiresAt.Format(time.RFC3339),
		})
		if err != nil {
			return err
		}

		// Record lockout as suspicious activity
		lockoutMetadata, _ := json.Marshal(map[string]interface{}{
			"action":           "account_locked",
			"reason":           "too_many_failed_logins",
			"lockout_duration": lockoutDuration.String(),
		})

		highSeverity := domain.HighActivity
		_, err = s.suspiciousRepo.CreateActivity(ctx, domain.CreateSuspiciousActivityAction{
			UserID:       userID,
			ActivityType: "account_locked",
			IPAddress:    ipAddress,
			UserAgent:    userAgent,
			Description:  "Account locked due to multiple failed login attempts",
			Metadata:     lockoutMetadata,
			Severity:     &highSeverity,
		})
	}

	return err
}

func (s *securityService) RecordSuccessfulLogin(ctx context.Context, userID int64, email, ipAddress, userAgent string) error {
	// 1. Record the successful login attempt
	_, err := s.loginAttemptsRepo.CreateLoginAttempt(ctx, domain.CreateLoginAttemptAction{
		UserID:    &userID,
		Email:     email,
		IPAddress: ipAddress,
		UserAgent: &userAgent,
		Success:   true,
	})
	if err != nil {
		return err
	}

	// 2. Get user to check if this is a new device/location
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// 3. Record successful login as suspicious activity (for monitoring)
	metadata, _ := json.Marshal(map[string]interface{}{
		"action":     "successful_login",
		"timestamp":  time.Now().Unix(),
		"last_login": user.LastLoginAt,
	})

	lowSeverity := domain.LowActivity
	_, err = s.suspiciousRepo.CreateActivity(ctx, domain.CreateSuspiciousActivityAction{
		UserID:       userID,
		ActivityType: "successful_login",
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Description:  "Successful login",
		Metadata:     metadata,
		Severity:     &lowSeverity,
	})

	return err
}

func (s *securityService) CheckAccountLockout(ctx context.Context, userID int64, ipAddress string) error {
	// Check if user is locked out
	lockout, err := s.lockoutRepo.GetLockoutByUserID(ctx, userID)
	if err != nil {
		// No lockout found, user can proceed
		return nil
	}

	// Check if lockout has expired
	if time.Now().After(lockout.ExpiresAt) {
		// Lockout expired, remove it
		err = s.lockoutRepo.DeleteLockoutByID(ctx, lockout.ID)
		if err != nil {
			return err
		}
		return nil
	}

	// User is still locked out
	return fmt.Errorf("account is locked until %s", lockout.ExpiresAt.Format(time.RFC3339))
}

func (s *securityService) RecordSuspiciousActivity(ctx context.Context, req domain.CreateSuspiciousActivityAction) error {
	_, err := s.suspiciousRepo.CreateActivity(ctx, req)
	return err
}

func (s *securityService) GetSuspiciousActivities(ctx context.Context, userID int64) ([]domain.SuspiciousActivity, error) {
	return s.suspiciousRepo.GetActivitiesByUserID(ctx, userID, 50)
}

func (s *securityService) ResolveSuspiciousActivity(ctx context.Context, activityID int64) error {
	return s.suspiciousRepo.ResolveActivity(ctx, activityID)
}

func (s *securityService) CleanupExpiredLockouts(ctx context.Context) error {
	return s.lockoutRepo.DeleteExpiredLockouts(ctx)
}
