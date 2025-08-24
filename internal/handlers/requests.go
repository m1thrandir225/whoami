package handlers

import "github.com/m1thrandir225/whoami/internal/domain"

type registerRequest struct {
	Email           string                  `json:"email"`
	Password        string                  `json:"password"`
	Username        *string                 `json:"username"`
	PrivacySettings *domain.PrivacySettings `json:"privacy_settings"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type updatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type resolveSuspiciousActivityRequest struct {
	ActivityID int64 `json:"activity_id" binding:"required"`
}

type verifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}
