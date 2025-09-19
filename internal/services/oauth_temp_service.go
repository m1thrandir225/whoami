package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/redis/go-redis/v9"
)

type OAuthTempService interface {
	StoreTemporaryAuthData(ctx context.Context, authData *TempOAuthData) (string, error)
	GetTemporaryAuthData(ctx context.Context, tempToken string) (*TempOAuthData, error)
	DeleteTemporaryAuthData(ctx context.Context, tempToken string) error
}

type TempOAuthData struct {
	User                  domain.User `json:"user"`
	AccessToken           string      `json:"access_token"`
	RefreshToken          string      `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time   `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time   `json:"refresh_token_expires_at"`
	Device                interface{} `json:"device,omitempty"`
}

type oauthTempService struct {
	redis *redis.Client
}

func NewOAuthTempService(redisClient *redis.Client) OAuthTempService {
	return &oauthTempService{
		redis: redisClient,
	}
}

func (s *oauthTempService) StoreTemporaryAuthData(ctx context.Context, authData *TempOAuthData) (string, error) {
	// Generate secure temporary token
	tempToken, err := s.generateSecureToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate temporary token: %w", err)
	}

	// Serialize auth data
	authDataBytes, err := json.Marshal(authData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth data: %w", err)
	}

	// Store in Redis with 5 minute expiration
	key := fmt.Sprintf("oauth_temp:%s", tempToken)
	err = s.redis.SetEx(ctx, key, authDataBytes, 5*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store temporary auth data: %w", err)
	}

	return tempToken, nil
}

func (s *oauthTempService) GetTemporaryAuthData(ctx context.Context, tempToken string) (*TempOAuthData, error) {
	key := fmt.Sprintf("oauth_temp:%s", tempToken)

	authDataBytes, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("temporary token not found or expired")
		}
		return nil, fmt.Errorf("failed to get temporary auth data: %w", err)
	}

	var authData TempOAuthData
	err = json.Unmarshal(authDataBytes, &authData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal auth data: %w", err)
	}

	return &authData, nil
}

func (s *oauthTempService) DeleteTemporaryAuthData(ctx context.Context, tempToken string) error {
	key := fmt.Sprintf("oauth_temp:%s", tempToken)
	return s.redis.Del(ctx, key).Err()
}

func (s *oauthTempService) generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
