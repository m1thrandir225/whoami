package handlers

import (
	"errors"
	"fmt"
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
		PrivacySettings: *requestData.PrivacySettings,
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

	response := registerResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid credentials")))
		return
	}

	// Record successful login
	if err := h.securityService.RecordSuccessfulLogin(ctx, user.ID, requestData.Email, clientIP, userAgent); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Update last login time
	if err := h.userService.UpdateLastLogin(ctx, user.ID); err != nil {
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

	response := loginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

	ctx.Status(http.StatusOK)
}

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
	ctx.Status(http.StatusOK)
}

func (h *HTTPHandler) Logout(ctx *gin.Context) {
	_, err := GetCurrentUserPayload(ctx)
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

	if err := h.sessionService.RevokeSession(ctx, token); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

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
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := h.tokenMaker.CreateToken(payload.UserID, h.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := refreshTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   payload.ExpiredAt,
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

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

	ctx.JSON(http.StatusOK, messageResponse("Password updated successfully"))
}
