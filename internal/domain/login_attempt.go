package domain

import "time"

type LoginAttempt struct {
	ID            int64      `json:"id"`
	UserID        *int64     `json:"user_id"`
	Email         string     `json:"email"`
	IPAddress     string     `json:"ip_address"`
	UserAgent     *string    `json:"user_agent"`
	Success       bool       `json:"success"`
	FailureReason *string    `json:"failure_reason"`
	CreatedAt     *time.Time `json:"created_at"`
}

type CreateLoginAttemptAction struct {
	UserID        *int64
	Email         string
	IPAddress     string
	UserAgent     *string
	Success       bool
	FailureReason *string
}
