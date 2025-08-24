package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) GetUserSessions(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	sessions, err := h.sessionService.GetUserSessions(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}

func (h *HTTPHandler) RevokeSession(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req revokeSessionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify the session belongs to the user
	session, err := h.sessionService.GetSession(ctx, req.Token)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("session not found")))
		return
	}
	if session.UserID != payload.UserID {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("not authorized to revoke this session")))
		return
	}

	if err := h.sessionService.RevokeSession(ctx, req.Token); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("Session revoked successfully"))
}

func (h *HTTPHandler) RevokeAllSessions(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req revokeAllSessionsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	reason := "User requested logout from all devices"
	if req.Reason != "" {
		reason = req.Reason
	}
	if err := h.sessionService.RevokeAllUserSessions(ctx, payload.UserID, reason); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("All sessions revoked successfully"))
}
