package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) GetAuditLogsByUserID(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := int32(50) // Default limit
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	logs, err := h.auditService.GetAuditLogsByUserID(ctx, userID, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
	})
}

func (h *HTTPHandler) GetAuditLogsByAction(ctx *gin.Context) {
	action := ctx.Param("action")
	if action == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("action is required")))
		return
	}

	limit := int32(50) // Default limit
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	logs, err := h.auditService.GetAuditLogsByAction(ctx, action, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
	})
}

func (h *HTTPHandler) GetAuditLogsByResourceType(ctx *gin.Context) {
	resourceType := ctx.Param("resource_type")
	if resourceType == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("resource_type is required")))
		return
	}

	limit := int32(50) // Default limit
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	logs, err := h.auditService.GetAuditLogsByResourceType(ctx, resourceType, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
	})
}

func (h *HTTPHandler) GetAuditLogsByResourceID(ctx *gin.Context) {
	resourceType := ctx.Param("resource_type")
	if resourceType == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("resource_type is required")))
		return
	}

	resourceID, err := strconv.ParseInt(ctx.Param("resource_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := int32(50) // Default limit
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	logs, err := h.auditService.GetAuditLogsByResourceID(ctx, resourceType, resourceID, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
	})
}

func (h *HTTPHandler) GetAuditLogsByIP(ctx *gin.Context) {
	ipAddress := ctx.Param("ip")
	if ipAddress == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("ip address is required")))
		return
	}

	limit := int32(50) // Default limit
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	logs, err := h.auditService.GetAuditLogsByIP(ctx, ipAddress, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
	})
}

func (h *HTTPHandler) GetAuditLogsByDateRange(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	if startDate == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("start_date is required")))
		return
	}

	endDate := ctx.Query("end_date")
	if endDate == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("end_date is required")))
		return
	}

	limit := int32(50) // Default limit
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	logs, err := h.auditService.GetAuditLogsByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
	})
}

func (h *HTTPHandler) GetRecentAuditLogs(ctx *gin.Context) {
	limit := int32(50) // Default limit
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	logs, err := h.auditService.GetRecentAuditLogs(ctx, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
	})
}

func (h *HTTPHandler) CleanupOldAuditLogs(ctx *gin.Context) {
	err := h.auditService.CleanupOldAuditLogs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messageResponse("Old audit logs cleaned up successfully"))
}
