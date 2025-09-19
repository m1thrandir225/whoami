package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/m1thrandir225/whoami/internal/services"
	"golang.org/x/oauth2"
)

type GitHubProvider struct {
	config *oauth2.Config
}

func NewGitHubProvider(config Config) *GitHubProvider {
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	return &GitHubProvider{
		config: oauthConfig,
	}
}

func (p *GitHubProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func (p *GitHubProvider) ExchangeCode(ctx context.Context, code string) (*services.OAuthUserInfo, error) {
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

func (p *GitHubProvider) GetUserInfo(ctx context.Context, accessToken string) (*services.OAuthUserInfo, error) {
	// First get user basic info
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// If email is empty (private), try to get primary email
	var email string
	if githubUser.Email != "" {
		email = githubUser.Email
	} else {
		// Get user's primary email address
		primaryEmail, err := p.getPrimaryEmail(ctx, accessToken)
		if err != nil {
			// If we can't get email, we cannot proceed
			return nil, fmt.Errorf("unable to get user email address - please make sure your GitHub email is public or the email scope is granted: %w", err)
		}
		email = primaryEmail
	}

	// Use name if available, otherwise use login
	name := githubUser.Name
	if name == "" {
		name = githubUser.Login
	}

	providerUserID := fmt.Sprintf("%d", githubUser.ID)

	return &services.OAuthUserInfo{
		ProviderUserID: providerUserID,
		Email:          &email,
		Name:           &name,
		AvatarURL:      &githubUser.AvatarURL,
	}, nil
}

// Add this helper method to get primary email
func (p *GitHubProvider) getPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get user emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user emails: %s", resp.Status)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode user emails: %w", err)
	}

	// Find primary verified email
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	// If no primary verified email, return any verified email
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email address found")
}

func (p *GitHubProvider) RefreshToken(ctx context.Context, refreshToken string) (*services.OAuthUserInfo, error) {
	// GitHub doesn't support refresh tokens in the standard OAuth2 flow
	// This would need to be implemented differently for GitHub
	return nil, fmt.Errorf("github doesn't support refresh tokens in standard oauth2 flow")
}

func (p *GitHubProvider) GetProviderName() string {
	return "github"
}
