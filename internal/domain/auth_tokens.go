package domain

import "time"

type RefreshToken struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type AccessToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    int64     `json:"user_id"`
}

type CreateRefreshTokenAction struct {
	UserID     int64
	Token      string
	DeviceInfo []byte
	ExpiresAt  time.Time
}
