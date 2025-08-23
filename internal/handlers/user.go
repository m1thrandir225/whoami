package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/util"
)

func (h *HTTPHandler) Register(ctx *gin.Context) {
	var requestData registerRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
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

	accessToken, _, err := h.tokenMaker.CreateToken(user.ID, h.config.AccessTokenDuration)
	refreshToken, _, err := h.tokenMaker.CreateToken(user.ID, h.config.RefreshTokenDuration)

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

	user, err := h.userService.GetUserByEmail(ctx, requestData.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = util.ComparePassword(user.Password, requestData.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
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

	// TODO:
	// 1. Add the token to a blacklist in Redis
	// 2. Invalidate refresh tokens
	// 3. Log the logout event

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
	//TODO: implement

	ctx.JSON(http.StatusOK, messageResponse("Email verified successfully"))
}

func (h *HTTPHandler) ResendVerificationEmail(ctx *gin.Context) {
	//TODO: implement

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

	user, err := h.userService.GetUserByID(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	if err := util.ComparePassword(user.Password, req.CurrentPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	newPasswordHash, err := util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user.Password = newPasswordHash
	user.PasswordChangedAt = time.Now()

	err = h.userService.UpdateUser(ctx, *user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, messageResponse("Password updated successfully"))
}
