package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/util"
)

func (h *HTTPHandler) Register(ctx *gin.Context) {
	var requestData registerRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate password for new user
	if err := h.passwordSecurityService.ValidateNewUserPassword(ctx, requestData.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordHash, err := util.HashPassword(requestData.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := h.userService.CreateUser(ctx, domain.CreateUserAction{
		Email:           requestData.Email,
		Username:        requestData.Username,
		PrivacySettings: requestData.PrivacySettings,
		Password:        passwordHash,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Add initial password to history
	if err := h.passwordSecurityService.AddInitialPasswordToHistory(ctx, user.ID, passwordHash); err != nil {
		fmt.Printf("Warning: Failed to add initial password to history: %v\n", err)
	}

	if err := h.emailService.SendVerificationEmail(ctx, user.ID, user.Email); err != nil {
		fmt.Printf("Warning: Failed to send verification email: %v\n", err)
	}

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

	// Extract device information
	deviceInfo := security.ExtractDeviceInfo(ctx)

	// Create device for new user
	device, err := h.userDevicesService.GetOrCreateDevice(ctx, user.ID, deviceInfo)
	if err != nil {
		// Log error but don't fail the registration
		fmt.Printf("Warning: Failed to create device record: %v\n", err)
	}

	// Create session with device info
	deviceInfoMap := map[string]string{
		"device_id":   deviceInfo.DeviceID,
		"device_name": deviceInfo.DeviceName,
		"device_type": deviceInfo.DeviceType,
		"user_agent":  deviceInfo.UserAgent,
		"ip_address":  deviceInfo.IPAddress,
	}

	if err := h.sessionService.CreateSession(ctx, user.ID, accessToken, deviceInfoMap); err != nil {
		fmt.Printf("Warning: Failed to create session: %v\n", err)
	}

	// Log successful registration
	h.auditService.LogUserAction(ctx, user.ID, domain.AuditActionUserRegister, domain.AuditResourceTypeUser, user.ID, ctx.Request, map[string]interface{}{
		"email":   user.Email,
		"success": true,
	})

	response := registerResponse{
		User:                  *user,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		Device:                device,
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) Login(ctx *gin.Context) {
	var requestData loginRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	clientIP := security.GetClientIP(ctx)
	userAgent := ctx.GetHeader("User-Agent")

	user, err := h.userService.GetUserByEmail(ctx, requestData.Email)
	if err != nil {
		// Record failed login attempt even if user doesn't exist
		// This prevents user enumeration attacks
		h.securityService.RecordFailedLogin(ctx, 0, requestData.Email, clientIP, userAgent)

		// Log failed login attempt
		h.auditService.LogAnonymousAction(ctx, domain.AuditActionUserLogin, domain.AuditResourceTypeUser, 0, ctx.Request, map[string]interface{}{
			"email":   requestData.Email,
			"success": false,
			"reason":  "user_not_found",
		})

		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid credentials")))
		return
	}

	// Check if account is locked
	if err := h.securityService.CheckAccountLockout(ctx, user.ID, clientIP); err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	// Verify password
	if err := util.ComparePassword(user.Password, requestData.Password); err != nil {
		// Record failed login attempt
		h.securityService.RecordFailedLogin(ctx, user.ID, requestData.Email, clientIP, userAgent)

		// Log failed login attempt
		h.auditService.LogUserAction(ctx, user.ID, domain.AuditActionUserLogin, domain.AuditResourceTypeUser, user.ID, ctx.Request, map[string]interface{}{
			"email":   requestData.Email,
			"success": false,
			"reason":  "invalid_password",
		})

		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid credentials")))
		return
	}

	// Record successful login
	if err := h.securityService.RecordSuccessfulLogin(ctx, user.ID, requestData.Email, clientIP, userAgent); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log successful login
	h.auditService.LogUserAction(ctx, user.ID, domain.AuditActionUserLogin, domain.AuditResourceTypeUser, user.ID, ctx.Request, map[string]interface{}{
		"email":   requestData.Email,
		"success": true,
	})

	// Update last login time
	if err := h.userService.UpdateLastLogin(ctx, user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
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

	// Extract device information
	deviceInfo := security.ExtractDeviceInfo(ctx)

	// Create or get device
	device, err := h.userDevicesService.GetOrCreateDevice(ctx, user.ID, deviceInfo)
	if err != nil {
		// Log error but don't fail the login
		fmt.Printf("Warning: Failed to create device record: %v\n", err)
	}

	// Update device last used time
	if device != nil {
		if _, err := h.userDevicesService.UpdateDeviceLastUsed(ctx, device.ID); err != nil {
			fmt.Printf("Warning: Failed to update device last used: %v\n", err)
		}
	}

	// Create session with device info
	deviceInfoMap := map[string]string{
		"device_id":   deviceInfo.DeviceID,
		"device_name": deviceInfo.DeviceName,
		"device_type": deviceInfo.DeviceType,
		"user_agent":  deviceInfo.UserAgent,
		"ip_address":  deviceInfo.IPAddress,
	}

	if err := h.sessionService.CreateSession(ctx, user.ID, accessToken, deviceInfoMap); err != nil {
		fmt.Printf("Warning: Failed to create session: %v\n", err)
	}

	response := loginResponse{
		User:                  *user,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		Device:                device,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) GetCurrentUser(ctx *gin.Context) {
	userPayload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, err := strconv.ParseInt(strconv.FormatInt(userPayload.UserID, 10), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// TODO: add verification that the user is the one deactivating themselves
func (h *HTTPHandler) DeactivateUser(ctx *gin.Context) {
	var requestData UriID
	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, err := strconv.ParseInt(requestData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.userService.DeactivateUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log user deactivation
	h.auditService.LogSystemAction(ctx, domain.AuditActionUserDeactivate, domain.AuditResourceTypeUser, userID, ctx.Request, map[string]interface{}{
		"user_id": userID,
		"success": true,
	})

	ctx.Status(http.StatusOK)
}

// TODO: add verification that the user is the one activating themselves
func (h *HTTPHandler) ActivateUser(ctx *gin.Context) {
	var requestData UriID
	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, err := strconv.ParseInt(requestData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.userService.ActivateUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log user activation
	h.auditService.LogSystemAction(ctx, domain.AuditActionUserActivate, domain.AuditResourceTypeUser, userID, ctx.Request, map[string]interface{}{
		"user_id": userID,
		"success": true,
	})

	ctx.Status(http.StatusOK)
}

func (h *HTTPHandler) UpdateUser(ctx *gin.Context) {
	var uriData UriID
	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var requestData updateUserRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userID, err := strconv.ParseInt(uriData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user.Email = requestData.Email
	user.Username = requestData.Username

	err = h.userService.UpdateUser(ctx, *user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log user update
	h.auditService.LogSystemAction(ctx, domain.AuditActionUserUpdate, domain.AuditResourceTypeUser, userID, ctx.Request, map[string]interface{}{
		"user_id":  userID,
		"email":    requestData.Email,
		"username": requestData.Username,
		"success":  true,
	})

	ctx.Status(http.StatusOK)
}

func (h *HTTPHandler) UpdateUserPrivacySettings(ctx *gin.Context) {
	var uriData UriID

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var requestData domain.PrivacySettings

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userID, err := strconv.ParseInt(uriData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.userService.UpdateUserPrivacySettings(ctx, userID, requestData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log privacy settings update
	h.auditService.LogSystemAction(ctx, domain.AuditActionPrivacySettings, domain.AuditResourceTypePrivacy, userID, ctx.Request, map[string]interface{}{
		"user_id":          userID,
		"privacy_settings": requestData,
		"success":          true,
	})

	ctx.Status(http.StatusOK)
}

func (h *HTTPHandler) Logout(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	authorizationHeader := ctx.GetHeader("authorization")
	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid authorization header")))
		return
	}
	token := fields[1]

	// Revoke the session (this will also blacklist the token)
	if err := h.sessionService.RevokeSession(ctx, token); err != nil {
		log.Printf("Warning: Failed to revoke session: %v", err)
	}

	// Log logout action
	h.auditService.LogUserAction(ctx, payload.UserID, domain.AuditActionUserLogout, domain.AuditResourceTypeSession, payload.UserID, ctx.Request, map[string]interface{}{
		"token":   token,
		"success": true,
	})

	ctx.JSON(http.StatusOK, messageResponse("Logged out successfully"))
}

func (h *HTTPHandler) RefreshToken(ctx *gin.Context) {
	var req refreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := h.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		// Log failed token refresh
		h.auditService.LogAnonymousAction(ctx, "token_refresh", domain.AuditResourceTypeSession, 0, ctx.Request, map[string]interface{}{
			"token":   req.RefreshToken,
			"success": false,
			"reason":  err.Error(),
		})

		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := h.tokenMaker.CreateToken(payload.UserID, h.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log successful token refresh
	h.auditService.LogUserAction(ctx, payload.UserID, "token_refresh", domain.AuditResourceTypeSession, payload.UserID, ctx.Request, map[string]interface{}{
		"success": true,
	})

	response := refreshTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) VerifyEmail(ctx *gin.Context) {
	var req verifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.emailService.VerifyEmailToken(ctx, req.Token); err != nil {
		// Log failed email verification
		h.auditService.LogAnonymousAction(ctx, domain.AuditActionEmailVerify, domain.AuditResourceTypeEmail, 0, ctx.Request, map[string]interface{}{
			"token":   req.Token,
			"success": false,
			"reason":  err.Error(),
		})

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Log successful email verification
	h.auditService.LogSystemAction(ctx, domain.AuditActionEmailVerify, domain.AuditResourceTypeEmail, 0, ctx.Request, map[string]interface{}{
		"token":   req.Token,
		"success": true,
	})

	ctx.JSON(http.StatusOK, messageResponse("Email verified successfully"))
}

func (h *HTTPHandler) ResendVerificationEmail(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := h.userService.GetUserByID(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.EmailVerified {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("email already verified")))
		return
	}

	if err := h.emailService.ResendVerificationEmail(ctx, user.ID, user.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log email resend action
	h.auditService.LogUserAction(ctx, payload.UserID, domain.AuditActionEmailResend, domain.AuditResourceTypeEmail, user.ID, ctx.Request, map[string]interface{}{
		"email":   user.Email,
		"success": true,
	})

	ctx.JSON(http.StatusOK, messageResponse("Verification email resent successfully"))
}

func (h *HTTPHandler) UpdatePassword(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req updatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify current password
	user, err := h.userService.GetUserByID(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.ComparePassword(user.Password, req.CurrentPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("current password is incorrect")))
		return
	}

	if err := h.passwordSecurityService.UpdatePassword(ctx, payload.UserID, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	h.auditService.LogUserAction(ctx, payload.UserID, domain.AuditActionPasswordChange, domain.AuditResourceTypeUser, payload.UserID, ctx.Request, map[string]interface{}{
		"email":   user.Email,
		"success": true,
	})

	ctx.JSON(http.StatusOK, messageResponse("Password updated successfully"))
}
