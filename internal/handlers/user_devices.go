package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type updateDeviceRequest struct {
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	Trusted    bool   `json:"trusted"`
}

type markDeviceTrustedRequest struct {
	Trusted bool `json:"trusted"`
}

func (h *HTTPHandler) GetUserDevices(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	devices, err := h.userDevicesService.GetUserDevices(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log device retrieval
	h.auditService.LogUserAction(ctx, payload.UserID, "get_devices", domain.AuditResourceTypeDevice, payload.UserID, ctx.Request, map[string]interface{}{
		"count":   len(devices),
		"success": true,
	})

	ctx.JSON(http.StatusOK, devices)
}

func (h *HTTPHandler) GetUserDevice(ctx *gin.Context) {
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

	deviceID, err := strconv.ParseInt(requestData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	device, err := h.userDevicesService.GetUserDevice(ctx, deviceID, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	// Log device retrieval
	h.auditService.LogUserAction(ctx, payload.UserID, "get_device", domain.AuditResourceTypeDevice, deviceID, ctx.Request, map[string]interface{}{
		"device_id": device.DeviceID,
		"success":   true,
	})

	ctx.JSON(http.StatusOK, device)
}

func (h *HTTPHandler) UpdateUserDevice(ctx *gin.Context) {
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

	var req updateDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deviceID, err := strconv.ParseInt(requestData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	device, err := h.userDevicesService.UpdateDevice(ctx, domain.UpdateUserDeviceAction{
		ID:         deviceID,
		UserID:     payload.UserID,
		DeviceName: req.DeviceName,
		DeviceType: req.DeviceType,
		Trusted:    req.Trusted,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log device update
	h.auditService.LogUserAction(ctx, payload.UserID, "update_device", domain.AuditResourceTypeDevice, deviceID, ctx.Request, map[string]interface{}{
		"device_id":   device.DeviceID,
		"device_name": req.DeviceName,
		"trusted":     req.Trusted,
		"success":     true,
	})

	ctx.JSON(http.StatusOK, device)
}

func (h *HTTPHandler) DeleteUserDevice(ctx *gin.Context) {
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

	deviceID, err := strconv.ParseInt(requestData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.userDevicesService.DeleteDevice(ctx, deviceID, payload.UserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log device deletion
	h.auditService.LogUserAction(ctx, payload.UserID, "delete_device", domain.AuditResourceTypeDevice, deviceID, ctx.Request, map[string]interface{}{
		"success": true,
	})

	ctx.JSON(http.StatusOK, messageResponse("Device deleted successfully"))
}

func (h *HTTPHandler) DeleteAllUserDevices(ctx *gin.Context) {
	payload, err := GetCurrentUserPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if err := h.userDevicesService.DeleteAllDevices(ctx, payload.UserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log all devices deletion
	h.auditService.LogUserAction(ctx, payload.UserID, "delete_all_devices", domain.AuditResourceTypeDevice, payload.UserID, ctx.Request, map[string]interface{}{
		"success": true,
	})

	ctx.JSON(http.StatusOK, messageResponse("All devices deleted successfully"))
}

func (h *HTTPHandler) MarkDeviceAsTrusted(ctx *gin.Context) {
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

	var req markDeviceTrustedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deviceID, err := strconv.ParseInt(requestData.ID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	device, err := h.userDevicesService.MarkDeviceAsTrusted(ctx, deviceID, payload.UserID, req.Trusted)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Log device trust status change
	h.auditService.LogUserAction(ctx, payload.UserID, "mark_device_trusted", domain.AuditResourceTypeDevice, deviceID, ctx.Request, map[string]interface{}{
		"device_id": device.DeviceID,
		"trusted":   req.Trusted,
		"success":   true,
	})

	ctx.JSON(http.StatusOK, device)
}
