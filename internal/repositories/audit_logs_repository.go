package repositories

import (
	"context"
	"encoding/json"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type AuditLogsRepository interface {
	CreateAuditLog(ctx context.Context, req domain.CreateAuditLogAction) (*domain.AuditLog, error)
	GetAuditLogsByUserID(ctx context.Context, userID int64, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByAction(ctx context.Context, action string, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByResourceType(ctx context.Context, resourceType string, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByResourceID(ctx context.Context, resourceType string, resourceID int64, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.AuditLog, error)
	GetAuditLogsByDateRange(ctx context.Context, startDate, endDate time.Time, limit int32) ([]domain.AuditLog, error)
	GetRecentAuditLogs(ctx context.Context, limit int32) ([]domain.AuditLog, error)
	DeleteOldAuditLogs(ctx context.Context) error
}

type auditLogsRepository struct {
	store db.Store
}

func NewAuditLogsRepository(store db.Store) AuditLogsRepository {
	return &auditLogsRepository{
		store: store,
	}
}

func (r *auditLogsRepository) CreateAuditLog(ctx context.Context, req domain.CreateAuditLogAction) (*domain.AuditLog, error) {
	var pgUserID pgtype.Int8
	if req.UserID != nil {
		pgUserID = pgtype.Int8{Int64: *req.UserID, Valid: true}
	}

	var pgResourceType pgtype.Text
	if req.ResourceType != nil {
		pgResourceType = pgtype.Text{String: *req.ResourceType, Valid: true}
	}

	var pgResourceID pgtype.Int8
	if req.ResourceID != nil {
		pgResourceID = pgtype.Int8{Int64: *req.ResourceID, Valid: true}
	}

	var ipAddress netip.Addr
	if req.IPAddress != nil {
		parsedIP, err := netip.ParseAddr(*req.IPAddress)
		if err != nil {
			return nil, err
		}
		ipAddress = parsedIP
	}

	userAgent := ""
	if req.UserAgent != nil {
		userAgent = *req.UserAgent
	}

	var details []byte
	if req.Details != nil {
		details = req.Details
	}

	dbLog, err := r.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		UserID:       pgUserID,
		Action:       req.Action,
		ResourceType: pgResourceType,
		ResourceID:   pgResourceID,
		IpAddress:    &ipAddress,
		UserAgent:    userAgent,
		Details:      details,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbLog), nil
}

func (r *auditLogsRepository) GetAuditLogsByUserID(ctx context.Context, userID int64, limit int32) ([]domain.AuditLog, error) {
	dbLogs, err := r.store.GetAuditLogsByUserID(ctx, db.GetAuditLogsByUserIDParams{
		UserID: pgtype.Int8{Int64: userID, Valid: true},
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	logs := make([]domain.AuditLog, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = *r.toDomain(log)
	}

	return logs, nil
}

func (r *auditLogsRepository) GetAuditLogsByAction(ctx context.Context, action string, limit int32) ([]domain.AuditLog, error) {
	dbLogs, err := r.store.GetAuditLogsByAction(ctx, db.GetAuditLogsByActionParams{
		Action: action,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	logs := make([]domain.AuditLog, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = *r.toDomain(log)
	}

	return logs, nil
}

func (r *auditLogsRepository) GetAuditLogsByResourceType(ctx context.Context, resourceType string, limit int32) ([]domain.AuditLog, error) {
	dbLogs, err := r.store.GetAuditLogsByResourceType(ctx, db.GetAuditLogsByResourceTypeParams{
		ResourceType: pgtype.Text{String: resourceType, Valid: true},
		Limit:        limit,
	})
	if err != nil {
		return nil, err
	}

	logs := make([]domain.AuditLog, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = *r.toDomain(log)
	}

	return logs, nil
}

func (r *auditLogsRepository) GetAuditLogsByResourceID(ctx context.Context, resourceType string, resourceID int64, limit int32) ([]domain.AuditLog, error) {
	dbLogs, err := r.store.GetAuditLogsByResourceID(ctx, db.GetAuditLogsByResourceIDParams{
		ResourceType: pgtype.Text{String: resourceType, Valid: true},
		ResourceID:   pgtype.Int8{Int64: resourceID, Valid: true},
		Limit:        limit,
	})
	if err != nil {
		return nil, err
	}

	logs := make([]domain.AuditLog, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = *r.toDomain(log)
	}

	return logs, nil
}

func (r *auditLogsRepository) GetAuditLogsByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.AuditLog, error) {
	parsedIP, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil, err
	}

	dbLogs, err := r.store.GetAuditLogsByIP(ctx, db.GetAuditLogsByIPParams{
		IpAddress: &parsedIP,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	logs := make([]domain.AuditLog, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = *r.toDomain(log)
	}

	return logs, nil
}

func (r *auditLogsRepository) GetAuditLogsByDateRange(ctx context.Context, startDate, endDate time.Time, limit int32) ([]domain.AuditLog, error) {
	dbLogs, err := r.store.GetAuditLogsByDateRange(ctx, db.GetAuditLogsByDateRangeParams{
		CreatedAt:   &startDate,
		CreatedAt_2: &endDate,
		Limit:       limit,
	})
	if err != nil {
		return nil, err
	}

	logs := make([]domain.AuditLog, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = *r.toDomain(log)
	}

	return logs, nil
}

func (r *auditLogsRepository) GetRecentAuditLogs(ctx context.Context, limit int32) ([]domain.AuditLog, error) {
	dbLogs, err := r.store.GetRecentAuditLogs(ctx, limit)
	if err != nil {
		return nil, err
	}

	logs := make([]domain.AuditLog, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = *r.toDomain(log)
	}

	return logs, nil
}

func (r *auditLogsRepository) DeleteOldAuditLogs(ctx context.Context) error {
	return r.store.DeleteOldAuditLogs(ctx)
}

func (r *auditLogsRepository) toDomain(dbLog db.AuditLog) *domain.AuditLog {
	var userID *int64
	if dbLog.UserID.Valid {
		userID = &dbLog.UserID.Int64
	}

	var resourceType *string
	if dbLog.ResourceType.Valid {
		resourceType = &dbLog.ResourceType.String
	}

	var resourceID *int64
	if dbLog.ResourceID.Valid {
		resourceID = &dbLog.ResourceID.Int64
	}

	var ipAddress string
	if dbLog.IpAddress != nil {
		ipAddress = dbLog.IpAddress.String()
	}

	var userAgent *string
	if dbLog.UserAgent != "" {
		userAgent = &dbLog.UserAgent
	}

	var details json.RawMessage
	if len(dbLog.Details) > 0 {
		details = dbLog.Details
	}

	return &domain.AuditLog{
		ID:           dbLog.ID,
		UserID:       userID,
		Action:       dbLog.Action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IPAddress:    &ipAddress,
		UserAgent:    userAgent,
		Details:      details,
		CreatedAt:    dbLog.CreatedAt,
	}
}
