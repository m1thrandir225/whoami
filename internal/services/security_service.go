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
	CheckAccountLockout(ctx context.Context, userID int64, ipAddress string) error
	RecordFailedLogin(ctx context.Context, userID int64, ipAddress string, userAgent string) error
	RecordSuccessfulLogin(ctx context.Context, userID int64, ipAddress string, userAgent string) error
	RecordSuspiciousActivity(ctx context.Context, req domain.CreateSuspiciousActivityAction) error
	GetSuspiciousActivities(ctx context.Context, userID int64) ([]domain.SuspiciousActivity, error)
	ResolveSuspiciousActivity(ctx context.Context, activityID int64) error
	CleanupExpiredLockouts(ctx context.Context) error
}

type securityService struct {
	lockoutRepo    repositories.AccountLockoutRepository
	suspiciousRepo repositories.SuspiciousActivityRepository
	userRepo       repositories.UserRepository
}

func NewSecurityService(
	lockoutRepo repositories.AccountLockoutRepository,
	suspiciousRepo repositories.SuspiciousActivityRepository,
	userRepo repositories.UserRepository,
) SecurityService {
	return &securityService{
		lockoutRepo:    lockoutRepo,
		suspiciousRepo: suspiciousRepo,
		userRepo:       userRepo,
	}
}

// CheckAccountLockout checks if the account is locked out for the given user and IP address
func (s *securityService) CheckAccountLockout(ctx context.Context, userID int64, ipAddress string) error {
	// Check for user-specific lockout
	userLockout, err := s.lockoutRepo.GetLockoutByUserID(ctx, userID)
	if err == nil && userLockout != nil {
		return fmt.Errorf("account is locked until %s", userLockout.ExpiresAt.Format(time.RFC3339))
	}

	// Check for IP-specific lockout
	ipLockout, err := s.lockoutRepo.GetLockoutByIP(ctx, ipAddress)
	if err == nil && ipLockout != nil {
		return fmt.Errorf("IP address is locked until %s", ipLockout.ExpiresAt.Format(time.RFC3339))
	}

	// Check for user+IP combination lockout
	userIPLockout, err := s.lockoutRepo.GetLockoutByUserAndIP(ctx, userID, ipAddress)
	if err == nil && userIPLockout != nil {
		return fmt.Errorf("account is locked for this IP until %s", userIPLockout.ExpiresAt.Format(time.RFC3339))
	}

	return nil
}

// RecordFailedLogin records a failed login attempt and checks if the account should be locked
func (s *securityService) RecordFailedLogin(ctx context.Context, userID int64, ipAddress string, userAgent string) error {
	// Record suspicious activity
	metadata, _ := json.Marshal(map[string]interface{}{
		"action":    "failed_login",
		"timestamp": time.Now().Unix(),
	})

	severity := domain.MediumActivity
	_, err := s.suspiciousRepo.CreateActivity(ctx, domain.CreateSuspiciousActivityAction{
		UserID:       userID,
		ActivityType: "failed_login",
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Description:  "Failed login attempt",
		Metadata:     metadata,
		Severity:     &severity,
	})
	if err != nil {
		return err
	}

	// Check if we should lock the account
	activityCount, err := s.suspiciousRepo.GetActivityCountByUser(ctx, userID)
	if err != nil {
		return err
	}

	// Lock account after 5 failed attempts in 24 hours
	if activityCount >= 5 {
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

		// Record the lockout as suspicious activity
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

func (s *securityService) RecordSuccessfulLogin(ctx context.Context, userID int64, ipAddress string, userAgent string) error {
	// Get user to check if this is a new device/location
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if this is a suspicious login (new IP, new user agent, etc.)
	metadata, _ := json.Marshal(map[string]interface{}{
		"action":     "successful_login",
		"timestamp":  time.Now().Unix(),
		"last_login": user.LastLoginAt,
	})

	severity := domain.LowActivity
	_, err = s.suspiciousRepo.CreateActivity(ctx, domain.CreateSuspiciousActivityAction{
		UserID:       userID,
		ActivityType: "successful_login",
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Description:  "Successful login",
		Metadata:     metadata,
		Severity:     &severity,
	})

	return err
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
