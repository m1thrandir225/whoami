// Package domain
package domain

import (
	"encoding/json"
	"time"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
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
	PasswordChangedAt time.Time       `json:"password_changed_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	CreatedAt         time.Time       `json:"created_at"`
}

func NewUserFromDBRow(dbRow db.User) (User, error) {
	var privacySettings PrivacySettings
	err := json.Unmarshal(dbRow.PrivacySettings, &privacySettings)
	if err != nil {
		return User{}, nil
	}
	return User{
		ID:                dbRow.ID,
		Email:             dbRow.Email,
		EmailVerified:     dbRow.EmailVerified,
		Password:          dbRow.PasswordHash,
		PasswordChangedAt: dbRow.PasswordChangedAt,
		Role:              dbRow.Role,
		Active:            dbRow.Active,
		PrivacySettings:   privacySettings,
		LastLoginAt:       dbRow.LastLoginAt,
		CreatedAt:         dbRow.CreatedAt,
		UpdatedAt:         dbRow.UpdatedAt,
	}, nil
}

type PrivacySettings struct {
	ShowEmail        bool `json:"show_email"`
	ShowLastLogin    bool `json:"show_last_login"`
	TwoFactorEnabled bool `json:"two_factor_enabled"`
}

type CreateUserRequest struct {
	Email           string          `json:"email"`
	Password        string          `json:"password"`
	PrivacySettings PrivacySettings `json:"privacy_settings"`
}
