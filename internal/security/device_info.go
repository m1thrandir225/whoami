package security

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeviceInfo struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	UserAgent  string `json:"user_agent"`
	IPAddress  string `json:"ip_address"`
}

func ExtractDeviceInfo(ctx *gin.Context) *DeviceInfo {
	userAgent := ctx.GetHeader("User-Agent")
	ipAddress := GetClientIP(ctx)

	// Generate or extract device ID from headers
	deviceID := ctx.GetHeader("X-Device-ID")
	if deviceID == "" {
		deviceID = generateDeviceID()
	}

	// Parse user agent for device info
	deviceName, deviceType := parseUserAgent(userAgent)

	return &DeviceInfo{
		DeviceID:   deviceID,
		DeviceName: deviceName,
		DeviceType: deviceType,
		UserAgent:  userAgent,
		IPAddress:  ipAddress,
	}
}

func generateDeviceID() string {
	return uuid.New().String()
}

func parseUserAgent(userAgent string) (deviceName, deviceType string) {
	ua := strings.ToLower(userAgent)

	// Detect device type
	switch {
	case strings.Contains(ua, "mobile"):
		deviceType = "mobile"
	case strings.Contains(ua, "tablet"):
		deviceType = "tablet"
	default:
		deviceType = "desktop"
	}

	// Detect OS and browser for device name
	var osInfo, browserInfo string

	// Detect OS
	switch {
	case strings.Contains(ua, "windows"):
		osInfo = "Windows"
	case strings.Contains(ua, "mac os"):
		osInfo = "macOS"
	case strings.Contains(ua, "linux"):
		osInfo = "Linux"
	case strings.Contains(ua, "android"):
		osInfo = "Android"
	case strings.Contains(ua, "ios"):
		osInfo = "iOS"
	default:
		osInfo = "Unknown"
	}

	// Detect browser
	switch {
	case strings.Contains(ua, "chrome"):
		browserInfo = "Chrome"
	case strings.Contains(ua, "firefox"):
		browserInfo = "Firefox"
	case strings.Contains(ua, "safari"):
		browserInfo = "Safari"
	case strings.Contains(ua, "edge"):
		browserInfo = "Edge"
	default:
		browserInfo = "Unknown"
	}

	// Generate device name
	deviceName = fmt.Sprintf("%s on %s", browserInfo, osInfo)

	return deviceName, deviceType
}
