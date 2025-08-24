package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/m1thrandir225/whoami/internal/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider(config Config) *GoogleProvider {
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       config.Scopes,
		Endpoint:     google.Endpoint,
	}

	return &GoogleProvider{
		config: oauthConfig,
	}
}

func (p *GoogleProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *GoogleProvider) ExchangeCode(ctx context.Context, code string) (*services.OAuthUserInfo, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	userInfo, err := p.GetUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}

	// Add token information
	userInfo.AccessToken = &token.AccessToken
	if token.RefreshToken != "" {
		userInfo.RefreshToken = &token.RefreshToken
	}
	if !token.Expiry.IsZero() {
		userInfo.TokenExpiresAt = &token.Expiry
	}

	return userInfo, nil
}

func (p *GoogleProvider) GetUserInfo(ctx context.Context, accessToken string) (*services.OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &services.OAuthUserInfo{
		ProviderUserID: googleUser.ID,
		Email:          &googleUser.Email,
		Name:           &googleUser.Name,
		AvatarURL:      &googleUser.Picture,
	}, nil
}

func (p *GoogleProvider) RefreshToken(ctx context.Context, refreshToken string) (*services.OAuthUserInfo, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := p.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	userInfo, err := p.GetUserInfo(ctx, newToken.AccessToken)
	if err != nil {
		return nil, err
	}

	// Update token information
	userInfo.AccessToken = &newToken.AccessToken
	if newToken.RefreshToken != "" {
		userInfo.RefreshToken = &newToken.RefreshToken
	}
	if !newToken.Expiry.IsZero() {
		userInfo.TokenExpiresAt = &newToken.Expiry
	}

	return userInfo, nil
}

func (p *GoogleProvider) GetProviderName() string {
	return "google"
}
