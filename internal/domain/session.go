package domain

import "time"

type Session struct {
	ID           string            `json:"id"`
	UserID       int64             `json:"user_id"`
	Token        string            `json:"token"`         //current access token
	RefreshToken string            `json:"refresh_token"` // refresh token
	DeviceInfo   map[string]string `json:"device_info"`
	IPAddress    string            `json:"ip_address"`
	UserAgent    string            `json:"user_agent"`
	CreatedAt    time.Time         `json:"created_at"`
	LastActive   time.Time         `json:"last_active"`
	IsActive     bool              `json:"is_active"`
}
