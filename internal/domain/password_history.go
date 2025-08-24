package domain

import "time"

type PasswordHistory struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreatePasswordHistory struct {
	UserID       int64
	PasswordHash string
}
