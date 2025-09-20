package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type registerResponse struct {
	User                  domain.User        `json:"user"`
	AccessToken           string             `json:"access_token"`
	AccessTokenExpiresAt  time.Time          `json:"access_token_expires_at"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time          `json:"refresh_token_expires_at"`
	Device                *domain.UserDevice `json:"device"`
}

type loginResponse struct {
	User                  domain.User        `json:"user"`
	AccessToken           string             `json:"access_token"`
	AccessTokenExpiresAt  time.Time          `json:"access_token_expires_at"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time          `json:"refresh_token_expires_at"`
	Device                *domain.UserDevice `json:"device"`
}

type refreshTokenResponse struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func messageResponse(message string) gin.H {
	return gin.H{"message": message}
}

type healthResponse struct {
	Status    string            `json:"status"`
	Services  map[string]string `json:"services"`
	Timestamp string            `json:"timestamp"`
}
