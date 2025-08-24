package security

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlacklist interface {
	BlacklistToken(ctx context.Context, token string, expiration time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	RemoveFromBlacklist(ctx context.Context, token string) error
	BlacklistUserTokens(ctx context.Context, userID int64, reason string) error
	GetBlacklistedTokensForUser(ctx context.Context, userID int64) ([]string, error)
}

type tokenBlacklist struct {
	redisClient *redis.Client
}

func NewTokenBlacklist(redisClient *redis.Client) TokenBlacklist {
	return &tokenBlacklist{
		redisClient: redisClient,
	}
}

func (tb *tokenBlacklist) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:token:%s", token)
	return tb.redisClient.Set(ctx, key, "blacklisted", expiration).Err()
}

func (tb *tokenBlacklist) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:token:%s", token)
	exists, err := tb.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (tb *tokenBlacklist) RemoveFromBlacklist(ctx context.Context, token string) error {
	key := fmt.Sprintf("blacklist:token:%s", token)
	return tb.redisClient.Del(ctx, key).Err()
}

func (tb *tokenBlacklist) BlacklistUserTokens(ctx context.Context, userID int64, reason string) error {
	// This would typically involve:
	// 1. Getting all active tokens for the user
	// 2. Blacklisting each token
	// 3. Logging the security event

	// For now, we'll use a pattern-based approach
	pattern := fmt.Sprintf("blacklist:user:%d:*", userID)
	_, err := tb.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	// Add a marker for the user's blacklisted status
	userBlacklistKey := fmt.Sprintf("blacklist:user:%d:status", userID)
	err = tb.redisClient.Set(ctx, userBlacklistKey, reason, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	// Log the security event
	eventKey := fmt.Sprintf("blacklist:events:%d:%d", userID, time.Now().Unix())
	err = tb.redisClient.Set(ctx, eventKey, reason, 7*24*time.Hour).Err()

	return err
}

func (tb *tokenBlacklist) GetBlacklistedTokensForUser(ctx context.Context, userID int64) ([]string, error) {
	pattern := fmt.Sprintf("blacklist:user:%d:*", userID)
	keys, err := tb.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	tokens := make([]string, 0, len(keys))
	for _, key := range keys {
		// Extract token from key
		if len(key) > 0 {
			tokens = append(tokens, key)
		}
	}

	return tokens, nil
}
