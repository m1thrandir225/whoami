package domain

import "time"

type AccountLockout struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	IPAddress   string `json:"-"`
	LockoutType string `json:"lockout_type"`
	ExpiresAt   time.Time
	CreatedAt   *time.Time
}

type LockoutType string

const (
	LockoutTypeAccount LockoutType = "account"
	LockoutTypeIP      LockoutType = "ip"
	LockoutTypeBoth    LockoutType = "both"
)

type CreateAccountLockoutAction struct {
	UserID      int64
	IPAddress   string
	LockoutType LockoutType
	ExpiresAt   string
}
