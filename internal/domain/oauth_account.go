package domain

import "time"

type OAuthAccount struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"user_id"`
	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	Email          *string    `json:"email,omitempty"`
	Name           *string    `json:"name,omitempty"`
	AvatarURL      *string    `json:"avatar_url,omitempty"`
	AccessToken    *string    `json:"-"` // Never expose in JSON
	RefreshToken   *string    `json:"-"` // Never expose in JSON
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type CreateOAuthAccountAction struct {
	UserID         int64
	Provider       string
	ProviderUserID string
	Email          *string
	Name           *string
	AvatarURL      *string
	AccessToken    *string
	RefreshToken   *string
	TokenExpiresAt *time.Time
}

type UpdateOAuthAccountAction struct {
	ID             int64
	UserID         int64
	Email          *string
	Name           *string
	AvatarURL      *string
	AccessToken    *string
	RefreshToken   *string
	TokenExpiresAt *time.Time
}

type UpdateOAuthTokensAction struct {
	ID             int64
	UserID         int64
	AccessToken    *string
	RefreshToken   *string
	TokenExpiresAt *time.Time
}

const (
	OAuthProviderGoogle  = "google"
	OAuthProviderGitHub  = "github"
	OAuthProviderDiscord = "discord"
	OAuthProviderTwitter = "twitter"
)
