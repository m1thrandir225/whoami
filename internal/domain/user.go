// Package domain
package domain

import (
	"time"
)

type User struct {
	ID                int64           `json:"id"`
	Email             string          `json:"email"`
	Username          string          `json:"username"`
	Password          string          `json:"-"`
	EmailVerified     bool            `json:"email_verified"`
	Role              string          `json:"role"`
	Active            bool            `json:"active"`
	PrivacySettings   PrivacySettings `json:"privacy_settings"`
	LastLoginAt       *time.Time      `json:"last_login_at"`
	PasswordChangedAt *time.Time      `json:"password_changed_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	CreatedAt         time.Time       `json:"created_at"`
}

type PrivacySettings struct {
	ShowEmail        bool `json:"show_email"`
	ShowLastLogin    bool `json:"show_last_login"`
	TwoFactorEnabled bool `json:"two_factor_enabled"`
}

type CreateUserAction struct {
	Email           string           `json:"email"`
	Password        string           `json:"password"`
	Username        *string          `json:"username"`
	PrivacySettings *PrivacySettings `json:"privacy_settings"`
}
