package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/domain"
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

	// Log suspicious activities retrieval
	h.auditService.LogUserAction(ctx, payload.UserID, domain.AuditActionSuspiciousActivity, domain.AuditResourceTypeAccount, payload.UserID, ctx.Request, map[string]interface{}{
		"action":  "get_suspicious_activities",
		"count":   len(activities),
		"success": true,
	})

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

	// Log suspicious activity resolution
	h.auditService.LogSystemAction(ctx, domain.AuditActionSuspiciousActivity, domain.AuditResourceTypeAccount, req.ActivityID, ctx.Request, map[string]interface{}{
		"activity_id": req.ActivityID,
		"action":      "resolve",
		"success":     true,
	})

	ctx.JSON(http.StatusOK, messageResponse("Activity resolved successfully"))
}

func (h *HTTPHandler) CleanupExpiredLockouts(ctx *gin.Context) {
	err := h.securityService.CleanupExpiredLockouts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log cleanup action
	h.auditService.LogSystemAction(ctx, "cleanup_expired_lockouts", domain.AuditResourceTypeAccount, 0, ctx.Request, map[string]interface{}{
		"action":  "cleanup_expired_lockouts",
		"success": true,
	})

	ctx.JSON(http.StatusOK, messageResponse("Expired lockouts cleaned up successfully"))
}

func (h *HTTPHandler) ResetRateLimit(ctx *gin.Context) {
	if h.config.Environment != "development" {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("this endpoint is only available in development mode")))
		return
	}

	var req resetRateLimitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var err error

	switch req.Type {
	case "all":
		err = h.rateLimiter.ResetAllRateLimits(ctx)
	case "ip":
		if req.IP == "" {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("IP address required for IP-based reset")))
			return
		}
		err = h.rateLimiter.ResetRateLimitByIP(ctx, req.IP)
	case "user":
		if req.UserID == 0 {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user ID required for user-based reset")))
			return
		}
		err = h.rateLimiter.ResetRateLimitByUser(ctx, req.UserID)
	default:
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid reset type. Use 'all', 'ip', or 'user'")))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("Rate limits reset successfully"))
}
