package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/util"
)

type registerRequest struct {
	Email           string                  `json:"email"`
	Password        string                  `json:"password"`
	Username        *string                 `json:"username"`
	PrivacySettings *domain.PrivacySettings `json:"privacy_settings"`
}

type registerResponse struct {
	User         domain.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

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

	// TODO: return response with tokens

	ctx.JSON(http.StatusOK, user)
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

	// TODO: return response with tokens

	ctx.JSON(http.StatusOK, user)
}

func (h *HTTPHandler) GetCurrentUser(ctx *gin.Context) {
	// TODO verify if current user is the same as the one with the JWT token
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
