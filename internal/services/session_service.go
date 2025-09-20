package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/redis/go-redis/v9"
)

// SessionService is responsible for managing the user active sessions
type SessionService interface {
	CreateSession(ctx context.Context, userID int64, accessToken, refreshToken string, deviceInfo map[string]string) error
	GetSession(ctx context.Context, token string) (*domain.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	GetSessionByID(ctx context.Context, sessionID string) (*domain.Session, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID int64, reason string) error
	GetUserSessions(ctx context.Context, userID int64) ([]domain.Session, error)
	UpdateSessionActivity(ctx context.Context, token string) error
	UpdateSessionTokens(ctx context.Context, sessionID, newAccessToken, newRefreshToken string) error
	RevokeSessionByToken(ctx context.Context, token string) error
	CleanupExpiredSessions(ctx context.Context) error
}

type sessionService struct {
	redisClient    *redis.Client
	tokenBlacklist security.TokenBlacklist
}

func NewSessionService(redisClient *redis.Client, tokenBlacklist security.TokenBlacklist) SessionService {
	return &sessionService{
		redisClient:    redisClient,
		tokenBlacklist: tokenBlacklist,
	}
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *sessionService) CreateSession(ctx context.Context, userID int64, accessToken, refreshToken string, deviceInfo map[string]string) error {
	sessionID := generateSessionID()

	session := &domain.Session{
		ID:           sessionID,
		UserID:       userID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		DeviceInfo:   deviceInfo,
		IPAddress:    deviceInfo["ip_address"],
		UserAgent:    deviceInfo["user_agent"],
		CreatedAt:    time.Now(),
		LastActive:   time.Now(),
		IsActive:     true,
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	expiration := 7 * 24 * time.Hour // 7 days

	// Store session by ID (primary key)
	sessionIDKey := fmt.Sprintf("session:id:%s", sessionID)
	if err := s.redisClient.Set(ctx, sessionIDKey, sessionData, expiration).Err(); err != nil {
		return fmt.Errorf("failed to store session by ID: %w", err)
	}

	// Store session by access token (for quick lookup)
	accessTokenKey := fmt.Sprintf("session:token:%s", accessToken)
	if err := s.redisClient.Set(ctx, accessTokenKey, sessionID, expiration).Err(); err != nil {
		return fmt.Errorf("failed to store session by access token: %w", err)
	}

	// Store session by refresh token (for refresh operations)
	refreshTokenKey := fmt.Sprintf("session:refresh:%s", refreshToken)
	if err := s.redisClient.Set(ctx, refreshTokenKey, sessionID, expiration).Err(); err != nil {
		return fmt.Errorf("failed to store session by refresh token: %w", err)
	}

	// Add to user's active sessions
	userSessionsKey := fmt.Sprintf("user_sessions:%d", userID)
	if err := s.redisClient.SAdd(ctx, userSessionsKey, sessionID).Err(); err != nil {
		return fmt.Errorf("failed to add session to user sessions: %w", err)
	}

	// Set expiration for user sessions set
	s.redisClient.Expire(ctx, userSessionsKey, expiration)

	return nil
}

func (s *sessionService) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	// First try to get session ID from access token
	accessTokenKey := fmt.Sprintf("session:token:%s", token)
	sessionID, err := s.redisClient.Get(ctx, accessTokenKey).Result()
	if err != nil {
		return nil, fmt.Errorf("session not found for token: %w", err)
	}

	return s.GetSessionByID(ctx, sessionID)
}

func (s *sessionService) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	// Get session ID from refresh token
	refreshTokenKey := fmt.Sprintf("session:refresh:%s", refreshToken)
	sessionID, err := s.redisClient.Get(ctx, refreshTokenKey).Result()
	if err != nil {
		return nil, fmt.Errorf("session not found for refresh token: %w", err)
	}

	return s.GetSessionByID(ctx, sessionID)
}

func (s *sessionService) GetSessionByID(ctx context.Context, sessionID string) (*domain.Session, error) {
	sessionIDKey := fmt.Sprintf("session:id:%s", sessionID)
	sessionData, err := s.redisClient.Get(ctx, sessionIDKey).Result()
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	var session domain.Session
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	return &session, nil
}

func (s *sessionService) UpdateSessionTokens(ctx context.Context, sessionID, newAccessToken, newRefreshToken string) error {
	// Get current session
	session, err := s.GetSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Blacklist old tokens
	if err := s.tokenBlacklist.BlacklistToken(ctx, session.Token, 24*time.Hour); err != nil {
		fmt.Printf("Warning: failed to blacklist old access token: %v\n", err)
	}
	if err := s.tokenBlacklist.BlacklistToken(ctx, session.RefreshToken, 24*time.Hour); err != nil {
		fmt.Printf("Warning: failed to blacklist old refresh token: %v\n", err)
	}

	// Remove old token mappings
	oldAccessTokenKey := fmt.Sprintf("session:token:%s", session.Token)
	oldRefreshTokenKey := fmt.Sprintf("session:refresh:%s", session.RefreshToken)
	s.redisClient.Del(ctx, oldAccessTokenKey, oldRefreshTokenKey)

	// Update session with new tokens
	session.Token = newAccessToken
	session.RefreshToken = newRefreshToken
	session.LastActive = time.Now()

	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal updated session: %w", err)
	}

	expiration := 7 * 24 * time.Hour

	// Update session by ID
	sessionIDKey := fmt.Sprintf("session:id:%s", sessionID)
	if err := s.redisClient.Set(ctx, sessionIDKey, sessionData, expiration).Err(); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	// Create new token mappings
	newAccessTokenKey := fmt.Sprintf("session:token:%s", newAccessToken)
	newRefreshTokenKey := fmt.Sprintf("session:refresh:%s", newRefreshToken)

	if err := s.redisClient.Set(ctx, newAccessTokenKey, sessionID, expiration).Err(); err != nil {
		return fmt.Errorf("failed to create new access token mapping: %w", err)
	}

	if err := s.redisClient.Set(ctx, newRefreshTokenKey, sessionID, expiration).Err(); err != nil {
		return fmt.Errorf("failed to create new refresh token mapping: %w", err)
	}

	return nil
}

func (s *sessionService) RevokeSession(ctx context.Context, sessionID string) error {
	// Get session first to get token info
	session, err := s.GetSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session for revocation: %w", err)
	}

	// Remove from user's active sessions
	userSessionsKey := fmt.Sprintf("user_sessions:%d", session.UserID)
	if err := s.redisClient.SRem(ctx, userSessionsKey, sessionID).Err(); err != nil {
		fmt.Printf("Warning: failed to remove session from user sessions: %v\n", err)
	}

	// Delete all session keys
	sessionIDKey := fmt.Sprintf("session:id:%s", sessionID)
	accessTokenKey := fmt.Sprintf("session:token:%s", session.Token)
	refreshTokenKey := fmt.Sprintf("session:refresh:%s", session.RefreshToken)

	deletedCount, err := s.redisClient.Del(ctx, sessionIDKey, accessTokenKey, refreshTokenKey).Result()
	if err != nil {
		return fmt.Errorf("failed to delete session keys: %w", err)
	}

	fmt.Printf("Deleted %d session keys for session %s\n", deletedCount, sessionID)

	// Blacklist the tokens to prevent any remaining usage
	if err := s.tokenBlacklist.BlacklistToken(ctx, session.Token, 24*time.Hour); err != nil {
		fmt.Printf("Warning: failed to blacklist access token: %v\n", err)
	}
	if err := s.tokenBlacklist.BlacklistToken(ctx, session.RefreshToken, 24*time.Hour); err != nil {
		fmt.Printf("Warning: failed to blacklist refresh token: %v\n", err)
	}

	return nil
}

func (s *sessionService) RevokeAllUserSessions(ctx context.Context, userID int64, reason string) error {
	// Get all user sessions
	sessions, err := s.GetUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}

	// Revoke each session
	var revokeErrors []error
	for _, session := range sessions {
		if err := s.RevokeSession(ctx, session.ID); err != nil {
			revokeErrors = append(revokeErrors, fmt.Errorf("failed to revoke session %s: %w", session.ID, err))
			fmt.Printf("Failed to revoke session %s: %v\n", session.ID, err)
		}
	}

	// Clear the user sessions set
	userSessionsKey := fmt.Sprintf("user_sessions:%d", userID)
	if err := s.redisClient.Del(ctx, userSessionsKey).Err(); err != nil {
		fmt.Printf("Warning: failed to clear user sessions set: %v\n", err)
	}

	// Blacklist all user tokens (additional security measure)
	if err := s.tokenBlacklist.BlacklistUserTokens(ctx, userID, reason); err != nil {
		fmt.Printf("Warning: failed to blacklist user tokens: %v\n", err)
	}

	// Return the first error if any occurred, but continue processing all sessions
	if len(revokeErrors) > 0 {
		return revokeErrors[0]
	}

	return nil
}

func (s *sessionService) GetUserSessions(ctx context.Context, userID int64) ([]domain.Session, error) {
	userSessionsKey := fmt.Sprintf("user_sessions:%d", userID)
	sessionIDs, err := s.redisClient.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user session IDs: %w", err)
	}

	sessions := make([]domain.Session, 0, len(sessionIDs))
	for _, sessionID := range sessionIDs {
		session, err := s.GetSessionByID(ctx, sessionID)
		if err != nil {
			// Skip invalid/expired sessions but log the issue
			fmt.Printf("Warning: failed to get session %s, skipping: %v\n", sessionID, err)
			// Remove invalid session ID from user sessions
			s.redisClient.SRem(ctx, userSessionsKey, sessionID)
			continue
		}
		sessions = append(sessions, *session)
	}

	return sessions, nil
}

func (s *sessionService) UpdateSessionActivity(ctx context.Context, token string) error {
	session, err := s.GetSession(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to get session for activity update: %w", err)
	}

	session.LastActive = time.Now()
	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session for activity update: %w", err)
	}

	sessionIDKey := fmt.Sprintf("session:id:%s", session.ID)
	if err := s.redisClient.Set(ctx, sessionIDKey, sessionData, 7*24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to update session activity: %w", err)
	}

	return nil
}

func (s *sessionService) CleanupExpiredSessions(ctx context.Context) error {
	// This method cleans up orphaned session mappings and expired sessions
	// Note: Redis should handle most expiration automatically, but this catches edge cases

	// Get all user session keys
	userSessionKeys, err := s.redisClient.Keys(ctx, "user_sessions:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get user session keys: %w", err)
	}

	cleanedCount := 0
	for _, userSessionKey := range userSessionKeys {
		// Get all session IDs for this user
		sessionIDs, err := s.redisClient.SMembers(ctx, userSessionKey).Result()
		if err != nil {
			fmt.Printf("Warning: failed to get session IDs for key %s: %v\n", userSessionKey, err)
			continue
		}

		// Check each session ID
		for _, sessionID := range sessionIDs {
			sessionIDKey := fmt.Sprintf("session:id:%s", sessionID)
			exists, err := s.redisClient.Exists(ctx, sessionIDKey).Result()
			if err != nil {
				fmt.Printf("Warning: failed to check session existence %s: %v\n", sessionID, err)
				continue
			}

			// If session doesn't exist, remove it from user sessions
			if exists == 0 {
				s.redisClient.SRem(ctx, userSessionKey, sessionID)
				cleanedCount++
				fmt.Printf("Cleaned up orphaned session ID: %s\n", sessionID)
			}
		}
	}

	fmt.Printf("Cleaned up %d orphaned session references\n", cleanedCount)
	return nil
}

func (s *sessionService) RevokeSessionByToken(ctx context.Context, token string) error {
	session, err := s.GetSession(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	return s.RevokeSession(ctx, session.ID)
}
