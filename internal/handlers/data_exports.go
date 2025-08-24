package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type requestDataExportRequest struct {
	ExportType string `json:"export_type" binding:"required,oneof=user_data audit_logs login_history complete"`
}

func (h *HTTPHandler) RequestDataExport(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req requestDataExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	export, err := h.dataExportsService.RequestDataExport(ctx, payload.UserID, req.ExportType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log data export request
	h.auditService.LogUserAction(ctx, payload.UserID, "request_data_export", domain.AuditResourceTypeData, export.ID, ctx.Request, map[string]interface{}{
		"export_type": req.ExportType,
		"export_id":   export.ID,
		"success":     true,
	})

	ctx.JSON(http.StatusOK, export)
}

func (h *HTTPHandler) GetDataExports(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	exports, err := h.dataExportsService.GetDataExports(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log data exports retrieval
	h.auditService.LogUserAction(ctx, payload.UserID, "get_data_exports", domain.AuditResourceTypeData, payload.UserID, ctx.Request, map[string]interface{}{
		"count":   len(exports),
		"success": true,
	})

	ctx.JSON(http.StatusOK, exports)
}

func (h *HTTPHandler) GetDataExport(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var requestData UriID
	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	exportID, err := strconv.ParseInt(requestData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	export, err := h.dataExportsService.GetDataExport(ctx, exportID, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	// Log data export retrieval
	h.auditService.LogUserAction(ctx, payload.UserID, "get_data_export", domain.AuditResourceTypeData, exportID, ctx.Request, map[string]interface{}{
		"export_type": export.ExportType,
		"status":      export.Status,
		"success":     true,
	})

	ctx.JSON(http.StatusOK, export)
}

func (h *HTTPHandler) DownloadDataExport(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var requestData UriID
	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	exportID, err := strconv.ParseInt(requestData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	export, err := h.dataExportsService.GetDataExport(ctx, exportID, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if export.Status != domain.DataExportStatusCompleted {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("export is not ready for download")))
		return
	}

	if export.FilePath == nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("export file not found")))
		return
	}

	// Log data export download
	h.auditService.LogUserAction(ctx, payload.UserID, "download_data_export", domain.AuditResourceTypeData, exportID, ctx.Request, map[string]interface{}{
		"export_type": export.ExportType,
		"file_size":   export.FileSize,
		"success":     true,
	})

	// Set headers for file download
	filename := fmt.Sprintf("export_%d_%s.json", export.UserID, export.ExportType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", "application/json")

	// Serve the file
	ctx.File(*export.FilePath)
}

func (h *HTTPHandler) DeleteDataExport(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var requestData UriID
	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	exportID, err := strconv.ParseInt(requestData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.dataExportsService.DeleteDataExport(ctx, exportID, payload.UserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log data export deletion
	h.auditService.LogUserAction(ctx, payload.UserID, "delete_data_export", domain.AuditResourceTypeData, exportID, ctx.Request, map[string]interface{}{
		"success": true,
	})

	ctx.JSON(http.StatusOK, messageResponse("Data export deleted successfully"))
}
