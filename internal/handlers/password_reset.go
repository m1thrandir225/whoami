package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) RequestPasswordReset(ctx *gin.Context) {
	var req requestPasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Request password reset (this will send email if user exists)
	if err := h.passwordResetService.RequestPasswordReset(ctx, req.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Always return success to prevent user enumeration
	ctx.JSON(http.StatusOK, messageResponse("If an account with that email exists, a password reset link has been sent"))
}

func (h *HTTPHandler) ResetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Reset the password
	if err := h.passwordResetService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("Password reset successfully"))
}

func (h *HTTPHandler) VerifyResetToken(ctx *gin.Context) {
	var req verifyResetTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify the reset token
	reset, err := h.passwordResetService.VerifyResetToken(ctx, req.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"valid":      true,
		"expires_at": reset.ExpiresAt,
	})
}
