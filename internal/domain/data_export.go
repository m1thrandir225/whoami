package domain

import "time"

type DataExport struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	ExportType  string     `json:"export_type"`
	Status      string     `json:"status"`
	FilePath    *string    `json:"file_path,omitempty"`
	FileSize    *int64     `json:"file_size,omitempty"`
	ExpiresAt   time.Time  `json:"expires_at"`
	CreatedAt   *time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type CreateDataExportAction struct {
	UserID     int64
	ExportType string
	ExpiresAt  time.Time
}

type UpdateDataExportStatusAction struct {
	ID     int64
	UserID int64
	Status string
}

type UpdateDataExportFileAction struct {
	ID       int64
	UserID   int64
	FilePath string
	FileSize int64
}

const (
	DataExportStatusPending   = "pending"
	DataExportStatusCompleted = "completed"
	DataExportStatusFailed    = "failed"
	DataExportStatusExpired   = "expired"

	DataExportTypeUserData     = "user_data"
	DataExportTypeAuditLogs    = "audit_logs"
	DataExportTypeLoginHistory = "login_history"
	DataExportTypeComplete     = "complete"
)
