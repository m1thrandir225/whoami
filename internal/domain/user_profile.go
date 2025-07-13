package domain

import "time"

type UserProfile struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	AvatarUrl string    `json:"avatar_url"`
	Bio       string    `json:"bio"`
	Timezone  string    `json:"timezone"`
	Locale    string    `json:"locale"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserProfileRequest struct {
	UserID    int64
	FirstName string
	LastName  string
	Phone     string
	AvatarUrl string
	Bio       string
	Timezone  string
	Locale    string
}
