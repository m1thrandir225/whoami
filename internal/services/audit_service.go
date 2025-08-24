package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/repositories"
)

type AuditService interface {
	LogUserAction(ctx context.Context, userID int64, action string, resourceType string, resourceID int64, r *http.Request, details map[string]interface{}) error
	LogSystemAction(ctx context.Context, action string, resourceType string, resourceID int64, r *http.Request, details map[string]interface{}) error
	LogAnonymousAction(ctx context.Context, action string, resourceType string, resourceID int64, r *http.Request, details map[string]interface{}) error
	GetAuditLogsByUserID(ctx context.Context, userID int64, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByAction(ctx context.Context, action string, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByResourceType(ctx context.Context, resourceType string, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByResourceID(ctx context.Context, resourceType string, resourceID int64, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByDateRange(ctx context.Context, startDate, endDate string, limit int32) ([]domain.AuditLog, error)
	GetRecentAuditLogs(ctx context.Context, limit int32) ([]domain.AuditLog, error)
	CleanupOldAuditLogs(ctx context.Context) error
}

type auditService struct {
	auditRepo repositories.AuditLogsRepository
}

func NewAuditService(auditRepo repositories.AuditLogsRepository) AuditService {
	return &auditService{
		auditRepo: auditRepo,
	}
}

func (s *auditService) LogUserAction(ctx context.Context, userID int64, action string, resourceType string, resourceID int64, r *http.Request, details map[string]interface{}) error {
	detailsJSON, _ := json.Marshal(details)

	_, err := s.auditRepo.CreateAuditLog(ctx, domain.CreateAuditLogAction{
		UserID:       &userID,
		Action:       action,
		ResourceType: &resourceType,
		ResourceID:   &resourceID,
		IPAddress:    s.getClientIP(r),
		UserAgent:    s.getUserAgent(r),
		Details:      detailsJSON,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *auditService) LogSystemAction(ctx context.Context, action string, resourceType string, resourceID int64, r *http.Request, details map[string]interface{}) error {
	detailsJSON, _ := json.Marshal(details)

	_, err := s.auditRepo.CreateAuditLog(ctx, domain.CreateAuditLogAction{
		UserID:       nil, // System action, no user
		Action:       action,
		ResourceType: &resourceType,
		ResourceID:   &resourceID,
		IPAddress:    s.getClientIP(r),
		UserAgent:    s.getUserAgent(r),
		Details:      detailsJSON,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *auditService) LogAnonymousAction(ctx context.Context, action string, resourceType string, resourceID int64, r *http.Request, details map[string]interface{}) error {
	detailsJSON, _ := json.Marshal(details)

	_, err := s.auditRepo.CreateAuditLog(ctx, domain.CreateAuditLogAction{
		UserID:       nil, // Anonymous action
		Action:       action,
		ResourceType: &resourceType,
		ResourceID:   &resourceID,
		IPAddress:    s.getClientIP(r),
		UserAgent:    s.getUserAgent(r),
		Details:      detailsJSON,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *auditService) GetAuditLogsByUserID(ctx context.Context, userID int64, limit int32) ([]domain.AuditLog, error) {
	return s.auditRepo.GetAuditLogsByUserID(ctx, userID, limit)
}

func (s *auditService) GetAuditLogsByAction(ctx context.Context, action string, limit int32) ([]domain.AuditLog, error) {
	return s.auditRepo.GetAuditLogsByAction(ctx, action, limit)
}

func (s *auditService) GetAuditLogsByResourceType(ctx context.Context, resourceType string, limit int32) ([]domain.AuditLog, error) {
	return s.auditRepo.GetAuditLogsByResourceType(ctx, resourceType, limit)
}

func (s *auditService) GetAuditLogsByResourceID(ctx context.Context, resourceType string, resourceID int64, limit int32) ([]domain.AuditLog, error) {
	return s.auditRepo.GetAuditLogsByResourceID(ctx, resourceType, resourceID, limit)
}

func (s *auditService) GetAuditLogsByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.AuditLog, error) {
	return s.auditRepo.GetAuditLogsByIP(ctx, ipAddress, limit)
}

func (s *auditService) GetAuditLogsByDateRange(ctx context.Context, startDate, endDate string, limit int32) ([]domain.AuditLog, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, err
	}

	return s.auditRepo.GetAuditLogsByDateRange(ctx, start, end, limit)
}

func (s *auditService) GetRecentAuditLogs(ctx context.Context, limit int32) ([]domain.AuditLog, error) {
	return s.auditRepo.GetRecentAuditLogs(ctx, limit)
}

func (s *auditService) CleanupOldAuditLogs(ctx context.Context) error {
	return s.auditRepo.DeleteOldAuditLogs(ctx)
}

func (s *auditService) getClientIP(r *http.Request) *string {
	if r == nil {
		return nil
	}

	// Try to get real IP from headers
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	return &ip
}

func (s *auditService) getUserAgent(r *http.Request) *string {
	if r == nil {
		return nil
	}

	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return nil
	}

	return &userAgent
}
