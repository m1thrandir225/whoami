package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) GetSuspiciousActivities(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	activities, err := h.securityService.GetSuspiciousActivities(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"activities": activities,
	})
}

func (h *HTTPHandler) ResolveSuspiciousActivity(ctx *gin.Context) {
	_, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req resolveSuspiciousActivityRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.securityService.ResolveSuspiciousActivity(ctx, req.ActivityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("Activity resolved successfully"))
}

func (h *HTTPHandler) CleanupExpiredLockouts(ctx *gin.Context) {
	err := h.securityService.CleanupExpiredLockouts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("Expired lockouts cleaned up successfully"))
}
