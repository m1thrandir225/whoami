package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
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
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		fmt.Printf("Failed to marshal audit details: %v\n", err)
		detailsJSON = []byte("{}")
	}

	ipAddress := s.getClientIP(r)
	userAgent := s.getUserAgent(r)

	_, err = s.auditRepo.CreateAuditLog(ctx, domain.CreateAuditLogAction{
		UserID:       &userID,
		Action:       action,
		ResourceType: &resourceType,
		ResourceID:   &resourceID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Details:      detailsJSON,
	})
	if err != nil {
		fmt.Printf("Failed to create audit log: %v\n", err)
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
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		fmt.Printf("Failed to marshal audit details: %v\n", err)
		detailsJSON = []byte("{}")
	}

	ipAddress := s.getClientIP(r)
	userAgent := s.getUserAgent(r)

	_, err = s.auditRepo.CreateAuditLog(ctx, domain.CreateAuditLogAction{
		UserID:       nil, // Anonymous action
		Action:       action,
		ResourceType: &resourceType,
		ResourceID:   &resourceID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Details:      detailsJSON,
	})
	if err != nil {
		fmt.Printf("Failed to create audit log: %v\n", err)
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
	// Try various headers for the real IP
	headers := []string{
		"CF-Connecting-IP",    // Cloudflare
		"True-Client-IP",      // Cloudflare Enterprise
		"X-Real-IP",           // Nginx
		"X-Forwarded-For",     // Standard
		"X-Client-IP",         // Apache
		"X-Cluster-Client-IP", // Cluster
	}

	for _, header := range headers {
		ip := r.Header.Get(header)
		if ip != "" {
			// For X-Forwarded-For, take the first IP
			if header == "X-Forwarded-For" {
				ips := strings.Split(ip, ",")
				ip = strings.TrimSpace(ips[0])
			}

			// Validate the IP address
			if net.ParseIP(ip) != nil {
				return &ip
			}
		}
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If SplitHostPort fails, use RemoteAddr as is
		ip = r.RemoteAddr
	}

	// Validate the IP
	if net.ParseIP(ip) != nil {
		return &ip
	}

	// If all else fails, return localhost
	localhost := "127.0.0.1"
	return &localhost
}

func (s *auditService) getUserAgent(r *http.Request) *string {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return nil
	}
	return &userAgent
}
