package repositories

import (
	"context"
	"net/netip"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type SuspiciousActivityRepository interface {
	CreateActivity(ctx context.Context, req domain.CreateSuspiciousActivityAction) (*domain.SuspiciousActivity, error)
	GetActivitiesByUserID(ctx context.Context, userID int64, limit int32) ([]domain.SuspiciousActivity, error)
	GetActivitiesByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.SuspiciousActivity, error)
	GetUnresolvedActivities(ctx context.Context, limit int32) ([]domain.SuspiciousActivity, error)
	ResolveActivity(ctx context.Context, id int64) error
	GetActivityCountByUser(ctx context.Context, userID int64) (int64, error)
	GetActivityCountByIP(ctx context.Context, ipAddress string) (int64, error)
}

type suspiciousActivityRepository struct {
	store db.Store
}

func NewSuspiciousActivityRepository(store db.Store) SuspiciousActivityRepository {
	return &suspiciousActivityRepository{
		store: store,
	}
}

func (r *suspiciousActivityRepository) CreateActivity(ctx context.Context, req domain.CreateSuspiciousActivityAction) (*domain.SuspiciousActivity, error) {
	severity := "medium" // default severity
	if req.Severity != nil {
		severity = string(*req.Severity)
	}

	parsedIP, err := netip.ParseAddr(req.IPAddress)
	if err != nil {
		return nil, err
	}

	var pgSeverity pgtype.Text
	if req.Severity != nil {
		pgSeverity = pgtype.Text{String: string(*req.Severity), Valid: true}
	} else {
		pgSeverity = pgtype.Text{String: severity, Valid: true}
	}

	dbActivity, err := r.store.CreateSuspiciousActivity(ctx, db.CreateSuspiciousActivityParams{
		UserID:       pgtype.Int8{Int64: req.UserID, Valid: true},
		ActivityType: req.ActivityType,
		IpAddress:    parsedIP,
		UserAgent:    req.UserAgent,
		Description:  req.Description,
		Metadata:     req.Metadata,
		Severity:     pgSeverity,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbActivity), nil
}

func (r *suspiciousActivityRepository) GetActivitiesByUserID(ctx context.Context, userID int64, limit int32) ([]domain.SuspiciousActivity, error) {
	dbActivities, err := r.store.GetSuspiciousActivitiesByUserID(ctx, db.GetSuspiciousActivitiesByUserIDParams{
		UserID: pgtype.Int8{Int64: userID, Valid: true},
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	activities := make([]domain.SuspiciousActivity, len(dbActivities))
	for i, activity := range dbActivities {
		activities[i] = *r.toDomain(activity)
	}

	return activities, nil
}

func (r *suspiciousActivityRepository) GetActivitiesByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.SuspiciousActivity, error) {
	parsedIP, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil, err
	}

	dbActivities, err := r.store.GetSuspiciousActivitiesByIP(ctx, db.GetSuspiciousActivitiesByIPParams{
		IpAddress: parsedIP,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	activities := make([]domain.SuspiciousActivity, len(dbActivities))
	for i, activity := range dbActivities {
		activities[i] = *r.toDomain(activity)
	}

	return activities, nil
}

func (r *suspiciousActivityRepository) GetUnresolvedActivities(ctx context.Context, limit int32) ([]domain.SuspiciousActivity, error) {
	dbActivities, err := r.store.GetUnresolvedSuspiciousActivities(ctx, limit)
	if err != nil {
		return nil, err
	}

	activities := make([]domain.SuspiciousActivity, len(dbActivities))
	for i, activity := range dbActivities {
		activities[i] = *r.toDomain(activity)
	}

	return activities, nil
}

func (r *suspiciousActivityRepository) ResolveActivity(ctx context.Context, id int64) error {
	return r.store.ResolveSuspiciousActivity(ctx, id)
}

func (r *suspiciousActivityRepository) GetActivityCountByUser(ctx context.Context, userID int64) (int64, error) {
	return r.store.GetSuspiciousActivityCountByUser(ctx, pgtype.Int8{Int64: userID, Valid: true})
}

func (r *suspiciousActivityRepository) GetActivityCountByIP(ctx context.Context, ipAddress string) (int64, error) {
	parsedIP, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return 0, err
	}
	return r.store.GetSuspiciousActivityCountByIP(ctx, parsedIP)
}

func (r *suspiciousActivityRepository) toDomain(dbActivity db.SuspiciousActivity) *domain.SuspiciousActivity {
	return &domain.SuspiciousActivity{
		ID:           dbActivity.ID,
		UserID:       dbActivity.UserID.Int64,
		ActivityType: dbActivity.ActivityType,
		IPAddress:    dbActivity.IpAddress.String(),
		UserAgent:    dbActivity.UserAgent,
		Description:  dbActivity.Description,
		Metadata:     dbActivity.Metadata,
		Severity:     dbActivity.Severity.String,
		Resolved:     &dbActivity.Resolved.Bool,
		CreatedAt:    dbActivity.CreatedAt,
	}
}
