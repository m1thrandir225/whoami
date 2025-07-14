package domain

import "time"

type OAuthAccount struct {
	ID               int64      `json:"id"`
	UserID           int64      `json:"user_id"`
	Provider         string     `json:"provider"`
	ProviderUserID   string     `json:"provider_user_id"`
	ProviderUsername *string    `json:"provider_username"`
	ProviderEmail    *string    `json:"provider_email"`
	AccessToken      *string    `json:"access_token"`
	RefreshToken     *string    `json:"refresh_token"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

type CreateOAuthAccountAction struct {
	UserID           int64
	Provider         string
	ProviderUserID   string
	ProviderUsername *string
	ProviderEmail    *string
	AccessToken      *string
	RefreshToken     *string
}
