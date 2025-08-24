package domain

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID           int64           `json:"id"`
	UserID       *int64          `json:"user_id"`
	Action       string          `json:"action"`
	ResourceType *string         `json:"resource_type"`
	ResourceID   *int64          `json:"resource_id"`
	IPAddress    *string         `json:"ip_address"`
	UserAgent    *string         `json:"user_agent"`
	Details      json.RawMessage `json:"details"`
	CreatedAt    *time.Time      `json:"created_at"`
}

type CreateAuditLogAction struct {
	UserID       *int64
	Action       string
	ResourceType *string
	ResourceID   *int64
	IPAddress    *string
	UserAgent    *string
	Details      json.RawMessage
}

// Common audit actions
const (
	AuditActionUserLogin          = "user_login"
	AuditActionUserLogout         = "user_logout"
	AuditActionUserRegister       = "user_register"
	AuditActionUserUpdate         = "user_update"
	AuditActionUserDeactivate     = "user_deactivate"
	AuditActionUserActivate       = "user_activate"
	AuditActionPasswordChange     = "password_change"
	AuditActionPasswordReset      = "password_reset"
	AuditActionEmailVerify        = "email_verify"
	AuditActionEmailResend        = "email_resend"
	AuditActionSessionCreate      = "session_create"
	AuditActionSessionRevoke      = "session_revoke"
	AuditActionSessionRevokeAll   = "session_revoke_all"
	AuditActionAccountLockout     = "account_lockout"
	AuditActionSuspiciousActivity = "suspicious_activity"
	AuditActionDataExport         = "data_export"
	AuditActionPrivacySettings    = "privacy_settings"
)

// Common resource types
const (
	AuditResourceTypeUser     = "user"
	AuditResourceTypeSession  = "session"
	AuditResourceTypePassword = "password"
	AuditResourceTypeEmail    = "email"
	AuditResourceTypeAccount  = "account"
	AuditResourceTypeData     = "data"
	AuditResourceTypePrivacy  = "privacy"
	AuditResourceTypeDevice   = "device"
)
