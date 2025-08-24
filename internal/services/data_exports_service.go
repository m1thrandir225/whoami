package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/repositories"
)

type DataExportsService interface {
	RequestDataExport(ctx context.Context, userID int64, exportType string) (*domain.DataExport, error)
	GetDataExports(ctx context.Context, userID int64) ([]domain.DataExport, error)
	GetDataExport(ctx context.Context, id, userID int64) (*domain.DataExport, error)
	DeleteDataExport(ctx context.Context, id, userID int64) error
	ProcessPendingExports(ctx context.Context) error
	CleanupExpiredExports(ctx context.Context) error
}

type dataExportsService struct {
	dataExportsRepo   repositories.DataExportsRepository
	userRepo          repositories.UserRepository
	auditRepo         repositories.AuditLogsRepository
	loginAttemptsRepo repositories.LoginAttemptsRepository
	exportPath        string
}

func NewDataExportsService(
	dataExportsRepo repositories.DataExportsRepository,
	userRepo repositories.UserRepository,
	auditRepo repositories.AuditLogsRepository,
	loginAttemptsRepo repositories.LoginAttemptsRepository,
	exportPath string,
) DataExportsService {
	return &dataExportsService{
		dataExportsRepo:   dataExportsRepo,
		userRepo:          userRepo,
		auditRepo:         auditRepo,
		loginAttemptsRepo: loginAttemptsRepo,
		exportPath:        exportPath,
	}
}

func (s *dataExportsService) RequestDataExport(ctx context.Context, userID int64, exportType string) (*domain.DataExport, error) {
	// Validate export type
	switch exportType {
	case domain.DataExportTypeUserData, domain.DataExportTypeAuditLogs, domain.DataExportTypeLoginHistory, domain.DataExportTypeComplete:
	default:
		return nil, fmt.Errorf("invalid export type: %s", exportType)
	}

	// Create export request
	export, err := s.dataExportsRepo.CreateDataExport(ctx, domain.CreateDataExportAction{
		UserID:     userID,
		ExportType: exportType,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // 7 days
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create data export: %w", err)
	}

	return export, nil
}

func (s *dataExportsService) GetDataExports(ctx context.Context, userID int64) ([]domain.DataExport, error) {
	return s.dataExportsRepo.GetDataExportsByUserID(ctx, userID)
}

func (s *dataExportsService) GetDataExport(ctx context.Context, id, userID int64) (*domain.DataExport, error) {
	return s.dataExportsRepo.GetDataExportByID(ctx, id, userID)
}

func (s *dataExportsService) DeleteDataExport(ctx context.Context, id, userID int64) error {
	// Get export to check if file exists
	export, err := s.dataExportsRepo.GetDataExportByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// Delete file if it exists
	if export.FilePath != nil {
		if err := os.Remove(*export.FilePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete export file: %w", err)
		}
	}

	return s.dataExportsRepo.DeleteDataExport(ctx, id, userID)
}

func (s *dataExportsService) ProcessPendingExports(ctx context.Context) error {
	pendingExports, err := s.dataExportsRepo.GetPendingDataExports(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending exports: %w", err)
	}

	for _, export := range pendingExports {
		if err := s.processExport(ctx, &export); err != nil {
			// Mark as failed
			s.dataExportsRepo.UpdateDataExportStatus(ctx, domain.UpdateDataExportStatusAction{
				ID:     export.ID,
				UserID: export.UserID,
				Status: domain.DataExportStatusFailed,
			})
			continue
		}
	}

	return nil
}

func (s *dataExportsService) CleanupExpiredExports(ctx context.Context) error {
	// Get expired exports before deletion
	expiredExports, err := s.dataExportsRepo.GetDataExportsByUserID(ctx, 0) // This will need a different approach
	if err != nil {
		return fmt.Errorf("failed to get expired exports: %w", err)
	}

	// Delete files
	for _, export := range expiredExports {
		if export.FilePath != nil {
			os.Remove(*export.FilePath)
		}
	}

	return s.dataExportsRepo.DeleteExpiredDataExports(ctx)
}

func (s *dataExportsService) processExport(ctx context.Context, export *domain.DataExport) error {
	var data interface{}

	switch export.ExportType {
	case domain.DataExportTypeUserData:
		data = s.generateUserDataExport(ctx, export.UserID)
	case domain.DataExportTypeAuditLogs:
		data = s.generateAuditLogsExport(ctx, export.UserID)
	case domain.DataExportTypeLoginHistory:
		data = s.generateLoginHistoryExport(ctx, export.UserID)
	case domain.DataExportTypeComplete:
		data = s.generateCompleteExport(ctx, export.UserID)
	default:
		return fmt.Errorf("unknown export type: %s", export.ExportType)
	}

	// Create file
	filename := fmt.Sprintf("export_%d_%s_%d.json", export.UserID, export.ExportType, time.Now().Unix())
	filepath := filepath.Join(s.exportPath, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	// Write data as JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode export data: %w", err)
	}

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Update export with file info
	_, err = s.dataExportsRepo.UpdateDataExportFile(ctx, domain.UpdateDataExportFileAction{
		ID:       export.ID,
		UserID:   export.UserID,
		FilePath: filepath,
		FileSize: fileInfo.Size(),
	})
	if err != nil {
		return fmt.Errorf("failed to update export file info: %w", err)
	}

	// Mark as completed
	_, err = s.dataExportsRepo.UpdateDataExportStatus(ctx, domain.UpdateDataExportStatusAction{
		ID:     export.ID,
		UserID: export.UserID,
		Status: domain.DataExportStatusCompleted,
	})
	if err != nil {
		return fmt.Errorf("failed to mark export as completed: %w", err)
	}

	return nil
}

func (s *dataExportsService) generateUserDataExport(ctx context.Context, userID int64) map[string]interface{} {
	user, _ := s.userRepo.GetUserByID(ctx, userID)

	return map[string]interface{}{
		"export_type": "user_data",
		"exported_at": time.Now(),
		"user_data":   user,
	}
}

func (s *dataExportsService) generateAuditLogsExport(ctx context.Context, userID int64) map[string]interface{} {
	logs, _ := s.auditRepo.GetAuditLogsByUserID(ctx, userID, 1000)

	return map[string]interface{}{
		"export_type": "audit_logs",
		"exported_at": time.Now(),
		"audit_logs":  logs,
	}
}

func (s *dataExportsService) generateLoginHistoryExport(ctx context.Context, userID int64) map[string]interface{} {
	attempts, _ := s.loginAttemptsRepo.GetLoginAttemptsByUserID(ctx, userID, 1000)

	return map[string]interface{}{
		"export_type":    "login_history",
		"exported_at":    time.Now(),
		"login_attempts": attempts,
	}
}

func (s *dataExportsService) generateCompleteExport(ctx context.Context, userID int64) map[string]interface{} {
	user, _ := s.userRepo.GetUserByID(ctx, userID)
	logs, _ := s.auditRepo.GetAuditLogsByUserID(ctx, userID, 1000)
	attempts, _ := s.loginAttemptsRepo.GetLoginAttemptsByUserID(ctx, userID, 1000)

	return map[string]interface{}{
		"export_type":    "complete",
		"exported_at":    time.Now(),
		"user_data":      user,
		"audit_logs":     logs,
		"login_attempts": attempts,
	}
}
