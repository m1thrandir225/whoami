package domain

import "time"

type PasswordReset struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	TokenHash  string     `json:"token_hash"`
	HotpSecret string     `json:"-"`
	Counter    int64      `json:"-"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  *time.Time `json:"created_at"`
	UsedAt     *time.Time `json:"used_at"`
}

type CreatePasswordResetAction struct {
	UserID int64
}
