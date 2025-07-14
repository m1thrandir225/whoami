package domain

import (
	"time"
)

type ExportStatus string

const (
	PendingStatus    ExportStatus = "pending"
	ProcessingStatus ExportStatus = "processing"
	CompletedStatus  ExportStatus = "completed"
	FailedStatus     ExportStatus = "failed"
)

type DataExport struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Status      string     `json:"status"`
	ExportType  string     `json:"export_type"`
	FilePath    *string    `json:"file_path"`
	ExpiresAt   time.Time  `json:"expires_at"`
	CreatedAt   *time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type CreateDataExport struct {
	UserID     int64
	Status     ExportStatus
	ExportType string
	FilePath   *string
}
