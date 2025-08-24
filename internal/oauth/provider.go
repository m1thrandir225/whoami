package oauth

import (
	"context"

	"github.com/m1thrandir225/whoami/internal/services"
)

type Provider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*services.OAuthUserInfo, error)
	GetUserInfo(ctx context.Context, accessToken string) (*services.OAuthUserInfo, error)
	RefreshToken(ctx context.Context, refreshToken string) (*services.OAuthUserInfo, error)
	GetProviderName() string
}

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}
