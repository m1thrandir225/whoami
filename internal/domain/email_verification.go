package domain

import "time"

type EmailVerification struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	TokenHash string     `json:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt *time.Time `json:"created_at"`
	UsedAt    *time.Time `json:"used_at"`
}
