package repositories

import (
	"context"
	"net/netip"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type LoginAttemptsRepository interface {
	CreateLoginAttempt(ctx context.Context, req domain.CreateLoginAttemptAction) (*domain.LoginAttempt, error)
	GetLoginAttemptsByUserID(ctx context.Context, userID int64, limit int32) ([]domain.LoginAttempt, error)
	GetLoginAttemptsByEmail(ctx context.Context, email string, limit int32) ([]domain.LoginAttempt, error)
	GetLoginAttemptsByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.LoginAttempt, error)
	GetFailedLoginAttemptsByUserID(ctx context.Context, userID int64, limit int32) ([]domain.LoginAttempt, error)
	GetFailedLoginAttemptsByEmail(ctx context.Context, email string, limit int32) ([]domain.LoginAttempt, error)
	GetFailedLoginAttemptsByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.LoginAttempt, error)
	GetRecentFailedAttemptsByUserID(ctx context.Context, userID int64) ([]domain.LoginAttempt, error)
	GetRecentFailedAttemptsByEmail(ctx context.Context, email string) ([]domain.LoginAttempt, error)
	GetRecentFailedAttemptsByIP(ctx context.Context, ipAddress string) ([]domain.LoginAttempt, error)
	DeleteOldLoginAttempts(ctx context.Context) error
}

type loginAttemptsRepository struct {
	store db.Store
}

func NewLoginAttemptsRepository(store db.Store) LoginAttemptsRepository {
	return &loginAttemptsRepository{
		store: store,
	}
}

func (r *loginAttemptsRepository) CreateLoginAttempt(ctx context.Context, req domain.CreateLoginAttemptAction) (*domain.LoginAttempt, error) {
	parsedIP, err := netip.ParseAddr(req.IPAddress)
	if err != nil {
		return nil, err
	}

	var pgUserID pgtype.Int8
	if req.UserID != nil {
		pgUserID = pgtype.Int8{Int64: *req.UserID, Valid: true}
	}

	var pgFailureReason pgtype.Text
	if req.FailureReason != nil {
		pgFailureReason = pgtype.Text{String: *req.FailureReason, Valid: true}
	}

	dbAttempt, err := r.store.CreateLoginAttempt(ctx, db.CreateLoginAttemptParams{
		UserID:        pgUserID,
		Email:         req.Email,
		IpAddress:     parsedIP,
		UserAgent:     req.UserAgent,
		Success:       req.Success,
		FailureReason: pgFailureReason,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbAttempt), nil
}

func (r *loginAttemptsRepository) GetLoginAttemptsByUserID(ctx context.Context, userID int64, limit int32) ([]domain.LoginAttempt, error) {
	dbAttempts, err := r.store.GetLoginAttemptsByUserID(ctx, db.GetLoginAttemptsByUserIDParams{
		UserID: pgtype.Int8{Int64: userID, Valid: true},
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	attempts := make([]domain.LoginAttempt, len(dbAttempts))
	for i, attempt := range dbAttempts {
		attempts[i] = *r.toDomain(attempt)
	}

	return attempts, nil
}

func (r *loginAttemptsRepository) GetLoginAttemptsByEmail(ctx context.Context, email string, limit int32) ([]domain.LoginAttempt, error) {
	dbAttempts, err := r.store.GetLoginAttemptsByEmail(ctx, db.GetLoginAttemptsByEmailParams{
		Email: email,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	attempts := make([]domain.LoginAttempt, len(dbAttempts))
	for i, attempt := range dbAttempts {
		attempts[i] = *r.toDomain(attempt)
	}

	return attempts, nil
}

func (r *loginAttemptsRepository) GetLoginAttemptsByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.LoginAttempt, error) {
	parsedIP, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil, err
	}

	dbAttempts, err := r.store.GetLoginAttemptsByIP(ctx, db.GetLoginAttemptsByIPParams{
		IpAddress: parsedIP,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	attempts := make([]domain.LoginAttempt, len(dbAttempts))
	for i, attempt := range dbAttempts {
		attempts[i] = *r.toDomain(attempt)
	}

	return attempts, nil
}

func (r *loginAttemptsRepository) GetFailedLoginAttemptsByUserID(ctx context.Context, userID int64, limit int32) ([]domain.LoginAttempt, error) {
	dbAttempts, err := r.store.GetFailedLoginAttemptsByUserID(ctx, db.GetFailedLoginAttemptsByUserIDParams{
		UserID: pgtype.Int8{Int64: userID, Valid: true},
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	attempts := make([]domain.LoginAttempt, len(dbAttempts))
	for i, attempt := range dbAttempts {
		attempts[i] = *r.toDomain(attempt)
	}

	return attempts, nil
}

func (r *loginAttemptsRepository) GetFailedLoginAttemptsByEmail(ctx context.Context, email string, limit int32) ([]domain.LoginAttempt, error) {
	dbAttempts, err := r.store.GetFailedLoginAttemptsByEmail(ctx, db.GetFailedLoginAttemptsByEmailParams{
		Email: email,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	attempts := make([]domain.LoginAttempt, len(dbAttempts))
	for i, attempt := range dbAttempts {
		attempts[i] = *r.toDomain(attempt)
	}

	return attempts, nil
}

func (r *loginAttemptsRepository) GetFailedLoginAttemptsByIP(ctx context.Context, ipAddress string, limit int32) ([]domain.LoginAttempt, error) {
	parsedIP, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil, err
	}

	dbAttempts, err := r.store.GetFailedLoginAttemptsByIP(ctx, db.GetFailedLoginAttemptsByIPParams{
		IpAddress: parsedIP,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	attempts := make([]domain.LoginAttempt, len(dbAttempts))
	for i, attempt := range dbAttempts {
		attempts[i] = *r.toDomain(attempt)
	}

	return attempts, nil
}

func (r *loginAttemptsRepository) GetRecentFailedAttemptsByUserID(ctx context.Context, userID int64) ([]domain.LoginAttempt, error) {
	dbAttempts, err := r.store.GetRecentFailedAttemptsByUserID(ctx, pgtype.Int8{Int64: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	attempts := make([]domain.LoginAttempt, len(dbAttempts))
	for i, attempt := range dbAttempts {
		attempts[i] = *r.toDomain(attempt)
	}

	return attempts, nil
}

func (r *loginAttemptsRepository) GetRecentFailedAttemptsByEmail(ctx context.Context, email string) ([]domain.LoginAttempt, error) {
	dbAttempts, err := r.store.GetRecentFailedAttemptsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	attempts := make([]domain.LoginAttempt, len(dbAttempts))
	for i, attempt := range dbAttempts {
		attempts[i] = *r.toDomain(attempt)
	}

	return attempts, nil
}

func (r *loginAttemptsRepository) GetRecentFailedAttemptsByIP(ctx context.Context, ipAddress string) ([]domain.LoginAttempt, error) {
	parsedIP, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil, err
	}

	dbAttempts, err := r.store.GetRecentFailedAttemptsByIP(ctx, parsedIP)
	if err != nil {
		return nil, err
	}

	attempts := make([]domain.LoginAttempt, len(dbAttempts))
	for i, attempt := range dbAttempts {
		attempts[i] = *r.toDomain(attempt)
	}

	return attempts, nil
}

func (r *loginAttemptsRepository) DeleteOldLoginAttempts(ctx context.Context) error {
	return r.store.DeleteOldLoginAttempts(ctx)
}

func (r *loginAttemptsRepository) toDomain(dbAttempt db.LoginAttempt) *domain.LoginAttempt {
	var userID *int64
	if dbAttempt.UserID.Valid {
		userID = &dbAttempt.UserID.Int64
	}

	var failureReason *string
	if dbAttempt.FailureReason.Valid {
		failureReason = &dbAttempt.FailureReason.String
	}

	return &domain.LoginAttempt{
		ID:            dbAttempt.ID,
		UserID:        userID,
		Email:         dbAttempt.Email,
		IPAddress:     dbAttempt.IpAddress.String(),
		UserAgent:     dbAttempt.UserAgent,
		Success:       dbAttempt.Success,
		FailureReason: failureReason,
		CreatedAt:     dbAttempt.CreatedAt,
	}
}
