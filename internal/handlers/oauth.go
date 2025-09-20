package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/oauth"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/services"
)

type oauthCallbackRequest struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

type linkOAuthRequest struct {
	Provider string `json:"provider" binding:"required,oneof=google github"`
}

type exchangeTempTokenRequest struct {
	Token string `json:"token" binding:"required"`
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

	//Update last login time
	if err := h.userService.UpdateLastLogin(ctx, user.ID); err != nil {
		log.Printf("Warning: Failed to update last login time: %v", err)
	}

	// Generate tokens
	accessToken, accessPayload, err := h.tokenMaker.CreateToken(user.ID, h.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := h.tokenMaker.CreateToken(user.ID, h.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	deviceInfo := security.ExtractDeviceInfo(ctx)
	device, err := h.userDevicesService.GetOrCreateDevice(ctx, user.ID, deviceInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		device = &domain.UserDevice{
			DeviceID:   deviceInfo.DeviceID,
			DeviceName: deviceInfo.DeviceName,
			DeviceType: deviceInfo.DeviceType,
		}
	}

	deviceInfoMap := map[string]string{
		"device_id":   deviceInfo.DeviceID,
		"device_name": deviceInfo.DeviceName,
		"device_type": deviceInfo.DeviceType,
		"user_agent":  deviceInfo.UserAgent,
		"ip_address":  deviceInfo.IPAddress,
	}

	if err := h.sessionService.CreateSession(ctx, user.ID, accessToken, deviceInfoMap); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log successful OAuth authentication
	h.auditService.LogUserAction(ctx, user.ID, "oauth_authentication", domain.AuditResourceTypeUser, user.ID, ctx.Request, map[string]interface{}{
		"provider":      provider,
		"oauth_account": oauthAccount.ID,
		"success":       true,
	})

	tempAuthData := &services.TempOAuthData{
		User:                  *user,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		Device:                device,
	}

	tempToken, err := h.oauthTempService.StoreTemporaryAuthData(ctx, tempAuthData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Redirect to frontend callback with temporary token
	frontendURL := h.config.FrontendURL
	fmt.Println("Frontend URL: ", frontendURL)
	redirectURL := fmt.Sprintf("%s/oauth-callback?token=%s&success=true", frontendURL, tempToken)
	fmt.Println("Redirecting to: ", redirectURL)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)

}

func (h *HTTPHandler) ExchangeTempOAuthToken(ctx *gin.Context) {
	var req exchangeTempTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get auth data from temporary storage
	authData, err := h.oauthTempService.GetTemporaryAuthData(ctx, req.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid or expired token")))
		return
	}

	// Delete the temporary token (one-time use)
	if err := h.oauthTempService.DeleteTemporaryAuthData(ctx, req.Token); err != nil {
		// Log error but don't fail the request
		h.auditService.LogAnonymousAction(ctx, "oauth_temp_token_cleanup_failed", domain.AuditResourceTypeUser, 0, ctx.Request, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Log successful token exchange
	h.auditService.LogUserAction(ctx, authData.User.ID, "oauth_token_exchange", domain.AuditResourceTypeUser, authData.User.ID, ctx.Request, map[string]interface{}{
		"success": true,
	})

	var device *domain.UserDevice
	if authData.Device != nil {
		if devicePtr, ok := authData.Device.(*domain.UserDevice); ok {
			device = devicePtr
		}
	}

	// Return the auth data as a login response
	response := loginResponse{
		User:                  authData.User,
		AccessToken:           authData.AccessToken,
		RefreshToken:          authData.RefreshToken,
		AccessTokenExpiresAt:  authData.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: authData.RefreshTokenExpiresAt,
		Device:                device,
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
