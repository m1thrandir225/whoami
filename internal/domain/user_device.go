package domain

import "time"

type DeviceType string

const (
	MobileDevice  DeviceType = "mobile"
	DesktopDevice DeviceType = "desktop"
	TabletDevice  DeviceType = "tablet"
)

type UserDevice struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	DeviceID   string     `json:"device_id"`
	DeviceName *string    `json:"device_name"`
	DeviceType *string    `json:"device_type"`
	UserAgent  *string    `json:"user_agent"`
	IPAddress  *string    `json:"ip_address"`
	Trusted    *bool      `json:"trusted"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  *time.Time `json:"created_at"`
}

type CreateUserDeviceAction struct {
	UserID     int64
	DeviceID   string
	DeviceName *string
	DeviceType *DeviceType
	UserAgent  *string
	IPAddress  *string
	Trusted    *bool
}
