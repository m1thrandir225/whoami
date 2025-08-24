package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/oauth"
)

type oauthCallbackRequest struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

type linkOAuthRequest struct {
	Provider string `json:"provider" binding:"required,oneof=google github"`
}

func (h *HTTPHandler) OAuthLogin(ctx *gin.Context) {
	provider := ctx.Param("provider")

	// Validate provider
	switch provider {
	case domain.OAuthProviderGoogle, domain.OAuthProviderGitHub:
	default:
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("unsupported oauth provider")))
		return
	}

	// Get OAuth provider
	oauthProvider, err := h.getOAuthProvider(provider)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Generate state for CSRF protection
	state, err := h.oauthService.GenerateOAuthState()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get authorization URL
	authURL := oauthProvider.GetAuthURL(state)

	// Log OAuth login attempt
	h.auditService.LogAnonymousAction(ctx, "oauth_login", domain.AuditResourceTypeUser, 0, ctx.Request, map[string]interface{}{
		"provider": provider,
		"state":    state,
		"success":  true,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *HTTPHandler) OAuthCallback(ctx *gin.Context) {
	provider := ctx.Param("provider")

	// Validate provider
	switch provider {
	case domain.OAuthProviderGoogle, domain.OAuthProviderGitHub:
	default:
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("unsupported oauth provider")))
		return
	}

	var req oauthCallbackRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate state
	if !h.oauthService.ValidateOAuthState(req.State) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid oauth state")))
		return
	}

	// Get OAuth provider
	oauthProvider, err := h.getOAuthProvider(provider)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Exchange code for token and get user info
	userInfo, err := oauthProvider.ExchangeCode(ctx, req.Code)
	if err != nil {
		// Log failed OAuth callback
		h.auditService.LogAnonymousAction(ctx, "oauth_callback", domain.AuditResourceTypeUser, 0, ctx.Request, map[string]interface{}{
			"provider": provider,
			"success":  false,
			"reason":   err.Error(),
		})

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Authenticate or create user
	user, oauthAccount, err := h.oauthService.AuthenticateWithOAuth(ctx, provider, userInfo.ProviderUserID, userInfo)
	if err != nil {
		// Log failed OAuth authentication
		h.auditService.LogAnonymousAction(ctx, "oauth_authentication", domain.AuditResourceTypeUser, 0, ctx.Request, map[string]interface{}{
			"provider": provider,
			"success":  false,
			"reason":   err.Error(),
		})

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Generate tokens
	accessToken, _, err := h.tokenMaker.CreateToken(user.ID, h.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, _, err := h.tokenMaker.CreateToken(user.ID, h.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log successful OAuth authentication
	h.auditService.LogUserAction(ctx, user.ID, "oauth_authentication", domain.AuditResourceTypeUser, user.ID, ctx.Request, map[string]interface{}{
		"provider":      provider,
		"oauth_account": oauthAccount.ID,
		"success":       true,
	})

	response := loginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) LinkOAuthAccount(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req linkOAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get OAuth provider
	oauthProvider, err := h.getOAuthProvider(req.Provider)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Generate state for CSRF protection
	state, err := h.oauthService.GenerateOAuthState()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get authorization URL
	authURL := oauthProvider.GetAuthURL(state)

	// Log OAuth account linking attempt
	h.auditService.LogUserAction(ctx, payload.UserID, "oauth_link_attempt", domain.AuditResourceTypeUser, payload.UserID, ctx.Request, map[string]interface{}{
		"provider": req.Provider,
		"state":    state,
		"success":  true,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *HTTPHandler) GetOAuthAccounts(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accounts, err := h.oauthService.GetOAuthAccounts(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log OAuth accounts retrieval
	h.auditService.LogUserAction(ctx, payload.UserID, "get_oauth_accounts", domain.AuditResourceTypeUser, payload.UserID, ctx.Request, map[string]interface{}{
		"count":   len(accounts),
		"success": true,
	})

	ctx.JSON(http.StatusOK, accounts)
}

func (h *HTTPHandler) UnlinkOAuthAccount(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	provider := ctx.Param("provider")

	// Validate provider
	switch provider {
	case domain.OAuthProviderGoogle, domain.OAuthProviderGitHub:
	default:
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("unsupported oauth provider")))
		return
	}

	if err := h.oauthService.UnlinkOAuthAccount(ctx, payload.UserID, provider); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log OAuth account unlinking
	h.auditService.LogUserAction(ctx, payload.UserID, "unlink_oauth_account", domain.AuditResourceTypeUser, payload.UserID, ctx.Request, map[string]interface{}{
		"provider": provider,
		"success":  true,
	})

	ctx.JSON(http.StatusOK, messageResponse("OAuth account unlinked successfully"))
}

func (h *HTTPHandler) getOAuthProvider(providerName string) (oauth.Provider, error) {
	switch providerName {
	case domain.OAuthProviderGoogle:
		return h.oauthProviders.Google, nil
	case domain.OAuthProviderGitHub:
		return h.oauthProviders.GitHub, nil
	default:
		return nil, errors.New("unsupported oauth provider")
	}
}
