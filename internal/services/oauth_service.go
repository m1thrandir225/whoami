package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/repositories"
)

type OAuthService interface {
	LinkOAuthAccount(ctx context.Context, userID int64, provider, providerUserID string, userInfo *OAuthUserInfo) (*domain.OAuthAccount, error)
	UnlinkOAuthAccount(ctx context.Context, userID int64, provider string) error
	GetOAuthAccounts(ctx context.Context, userID int64) ([]domain.OAuthAccount, error)
	GetOAuthAccountByProvider(ctx context.Context, provider, providerUserID string) (*domain.OAuthAccount, error)
	UpdateOAuthTokens(ctx context.Context, accountID, userID int64, accessToken, refreshToken *string, expiresAt *time.Time) (*domain.OAuthAccount, error)
	AuthenticateWithOAuth(ctx context.Context, provider, providerUserID string, userInfo *OAuthUserInfo) (*domain.User, *domain.OAuthAccount, error)
	GenerateOAuthState() (string, error)
	ValidateOAuthState(state string) bool
}

type OAuthUserInfo struct {
	ProviderUserID string
	Email          *string
	Name           *string
	AvatarURL      *string
	AccessToken    *string
	RefreshToken   *string
	TokenExpiresAt *time.Time
}

type oauthService struct {
	oauthRepo  repositories.OAuthAccountsRepository
	userRepo   repositories.UserRepository
	stateStore map[string]time.Time // In production, use Redis
}

func NewOAuthService(
	oauthRepo repositories.OAuthAccountsRepository,
	userRepo repositories.UserRepository,
) OAuthService {
	return &oauthService{
		oauthRepo:  oauthRepo,
		userRepo:   userRepo,
		stateStore: make(map[string]time.Time),
	}
}

func (s *oauthService) LinkOAuthAccount(ctx context.Context, userID int64, provider, providerUserID string, userInfo *OAuthUserInfo) (*domain.OAuthAccount, error) {
	// Check if OAuth account already exists
	existingAccount, err := s.oauthRepo.GetOAuthAccountByProvider(ctx, provider, providerUserID)
	if err == nil {
		// Account exists, check if it's already linked to this user
		if existingAccount.UserID == userID {
			return existingAccount, nil // Already linked
		}
		return nil, fmt.Errorf("oauth account already linked to another user")
	}

	// Create new OAuth account
	account, err := s.oauthRepo.CreateOAuthAccount(ctx, domain.CreateOAuthAccountAction{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		Email:          userInfo.Email,
		Name:           userInfo.Name,
		AvatarURL:      userInfo.AvatarURL,
		AccessToken:    userInfo.AccessToken,
		RefreshToken:   userInfo.RefreshToken,
		TokenExpiresAt: userInfo.TokenExpiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth account: %w", err)
	}

	return account, nil
}

func (s *oauthService) UnlinkOAuthAccount(ctx context.Context, userID int64, provider string) error {
	return s.oauthRepo.DeleteOAuthAccountByProvider(ctx, userID, provider)
}

func (s *oauthService) GetOAuthAccounts(ctx context.Context, userID int64) ([]domain.OAuthAccount, error) {
	return s.oauthRepo.GetOAuthAccountsByUserID(ctx, userID)
}

func (s *oauthService) GetOAuthAccountByProvider(ctx context.Context, provider, providerUserID string) (*domain.OAuthAccount, error) {
	return s.oauthRepo.GetOAuthAccountByProvider(ctx, provider, providerUserID)
}

func (s *oauthService) UpdateOAuthTokens(ctx context.Context, accountID, userID int64, accessToken, refreshToken *string, expiresAt *time.Time) (*domain.OAuthAccount, error) {
	return s.oauthRepo.UpdateOAuthTokens(ctx, domain.UpdateOAuthTokensAction{
		ID:             accountID,
		UserID:         userID,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		TokenExpiresAt: expiresAt,
	})
}

func (s *oauthService) AuthenticateWithOAuth(ctx context.Context, provider, providerUserID string, userInfo *OAuthUserInfo) (*domain.User, *domain.OAuthAccount, error) {
	// Try to find existing OAuth account
	oauthAccount, err := s.oauthRepo.GetOAuthAccountByProvider(ctx, provider, providerUserID)
	if err == nil {
		// OAuth account exists, get the user
		user, err := s.userRepo.GetUserByID(ctx, oauthAccount.UserID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get user: %w", err)
		}

		// Update OAuth account with latest info
		updatedAccount, err := s.oauthRepo.UpdateOAuthAccount(ctx, domain.UpdateOAuthAccountAction{
			ID:             oauthAccount.ID,
			UserID:         oauthAccount.UserID,
			Email:          userInfo.Email,
			Name:           userInfo.Name,
			AvatarURL:      userInfo.AvatarURL,
			AccessToken:    userInfo.AccessToken,
			RefreshToken:   userInfo.RefreshToken,
			TokenExpiresAt: userInfo.TokenExpiresAt,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update oauth account: %w", err)
		}

		return user, updatedAccount, nil
	}

	// OAuth account doesn't exist, try to find user by email
	if userInfo.Email != nil {
		user, err := s.userRepo.GetUserByEmail(ctx, *userInfo.Email)
		if err == nil {
			// User exists, link the OAuth account
			oauthAccount, err := s.LinkOAuthAccount(ctx, user.ID, provider, providerUserID, userInfo)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to link oauth account: %w", err)
			}
			return user, oauthAccount, nil
		}
	}

	// No existing user found, create new user
	user, err := s.userRepo.CreateUser(ctx, domain.CreateUserAction{
		Email:           *userInfo.Email,
		PrivacySettings: &domain.PrivacySettings{}, // Default privacy settings
		Password:        "",                        // OAuth users don't need password
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Link OAuth account to new user
	oauthAccount, err = s.LinkOAuthAccount(ctx, user.ID, provider, providerUserID, userInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to link oauth account: %w", err)
	}

	return user, oauthAccount, nil
}

func (s *oauthService) GenerateOAuthState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(bytes)
	s.stateStore[state] = time.Now().Add(10 * time.Minute) // 10 minute expiry

	return state, nil
}

func (s *oauthService) ValidateOAuthState(state string) bool {
	expiry, exists := s.stateStore[state]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(s.stateStore, state)
		return false
	}

	delete(s.stateStore, state)
	return true
}
