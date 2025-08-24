package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/redis/go-redis/v9"
)

type SessionService interface {
	CreateSession(ctx context.Context, userID int64, token string, deviceInfo map[string]string) error
	GetSession(ctx context.Context, token string) (*domain.Session, error)
	RevokeSession(ctx context.Context, token string) error
	RevokeAllUserSessions(ctx context.Context, userID int64, reason string) error
	GetUserSessions(ctx context.Context, userID int64) ([]domain.Session, error)
	UpdateSessionActivity(ctx context.Context, token string) error
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

func (s *sessionService) CreateSession(ctx context.Context, userID int64, token string, deviceInfo map[string]string) error {
	session := &domain.Session{
		UserID:     userID,
		Token:      token,
		DeviceInfo: deviceInfo,
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
		IsActive:   true,
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}
	// Store session with token expiration
	key := fmt.Sprintf("session:%s", token)
	expiration := 7 * 24 * time.Hour // 7 days
	err = s.redisClient.Set(ctx, key, sessionData, expiration).Err()
	if err != nil {
		return err
	}

	// Add to user's active sessions
	userSessionsKey := fmt.Sprintf("user_sessions:%d", userID)
	err = s.redisClient.SAdd(ctx, userSessionsKey, token).Err()
	if err != nil {
		return err
	}

	// Set expiration for user sessions set
	return s.redisClient.Expire(ctx, userSessionsKey, expiration).Err()
}

func (s *sessionService) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	key := fmt.Sprintf("session:%s", token)
	sessionData, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var session domain.Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *sessionService) RevokeSession(ctx context.Context, token string) error {
	// Get session to find user ID
	session, err := s.GetSession(ctx, token)
	if err != nil {
		return err
	}

	// Remove from user's active sessions
	userSessionsKey := fmt.Sprintf("user_sessions:%d", session.UserID)
	err = s.redisClient.SRem(ctx, userSessionsKey, token).Err()
	if err != nil {
		return err
	}

	// Delete session
	key := fmt.Sprintf("session:%s", token)
	err = s.redisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	// Blacklist the token
	return s.tokenBlacklist.BlacklistToken(ctx, token, 24*time.Hour)
}

func (s *sessionService) RevokeAllUserSessions(ctx context.Context, userID int64, reason string) error {
	// Get all user sessions
	sessions, err := s.GetUserSessions(ctx, userID)
	if err != nil {
		return err
	}

	// Revoke each session
	for _, session := range sessions {
		err = s.RevokeSession(ctx, session.Token)
		if err != nil {
			// Log error but continue with other sessions
			fmt.Printf("Failed to revoke session %s: %v\n", session.Token, err)
		}
	}

	// Blacklist all user tokens
	return s.tokenBlacklist.BlacklistUserTokens(ctx, userID, reason)
}

func (s *sessionService) GetUserSessions(ctx context.Context, userID int64) ([]domain.Session, error) {
	userSessionsKey := fmt.Sprintf("user_sessions:%d", userID)
	tokens, err := s.redisClient.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]domain.Session, 0, len(tokens))
	for _, token := range tokens {
		session, err := s.GetSession(ctx, token)
		if err != nil {
			// Skip invalid sessions
			continue
		}
		sessions = append(sessions, *session)
	}

	return sessions, nil
}

func (s *sessionService) UpdateSessionActivity(ctx context.Context, token string) error {
	session, err := s.GetSession(ctx, token)
	if err != nil {
		return err
	}

	session.LastActive = time.Now()
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("session:%s", token)
	return s.redisClient.Set(ctx, key, sessionData, 7*24*time.Hour).Err()
}
