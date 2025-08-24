package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type PasswordResetRepository interface {
	CreatePasswordReset(ctx context.Context, req domain.CreatePasswordResetAction) error
	GetPasswordResetByToken(ctx context.Context, token string) (*domain.PasswordReset, error)
	MarkPasswordResetAsUsed(ctx context.Context, id int64) error
	DeleteUnusedPasswordResets(ctx context.Context, userID int64) error
	GetUnusedPasswordResets(ctx context.Context, userID int64) ([]domain.PasswordReset, error)
}

type passwordResetRepository struct {
	store db.Store
}

func NewPasswordResetRepository(store db.Store) PasswordResetRepository {
	return &passwordResetRepository{
		store: store,
	}
}

func (r *passwordResetRepository) CreatePasswordReset(ctx context.Context, req domain.CreatePasswordResetAction) error {
	_, err := r.store.CreatePasswordReset(ctx, db.CreatePasswordResetParams{
		UserID:     req.UserID,
		TokenHash:  req.TokenHash,
		HotpSecret: req.HotpSecret,
		Counter:    pgtype.Int8{Int64: req.Counter, Valid: true},
		ExpiresAt:  time.Now().Add(time.Hour * 1),
	})
	return err
}

func (r *passwordResetRepository) GetPasswordResetByToken(ctx context.Context, token string) (*domain.PasswordReset, error) {
	dbReset, err := r.store.GetPasswordResetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbReset), nil
}

func (r *passwordResetRepository) MarkPasswordResetAsUsed(ctx context.Context, id int64) error {
	return r.store.MarkPasswordResetAsUsed(ctx, id)
}

func (r *passwordResetRepository) DeleteUnusedPasswordResets(ctx context.Context, userID int64) error {
	return r.store.DeleteUnusedPasswordResets(ctx, userID)
}

func (r *passwordResetRepository) GetUnusedPasswordResets(ctx context.Context, userID int64) ([]domain.PasswordReset, error) {
	dbResets, err := r.store.GetUnusedPasswordResets(ctx, userID)
	if err != nil {
		return nil, err
	}

	resets := make([]domain.PasswordReset, len(dbResets))
	for i, reset := range dbResets {
		resets[i] = *r.toDomain(reset)
	}

	return resets, nil
}

func (r *passwordResetRepository) toDomain(dbReset db.PasswordReset) *domain.PasswordReset {
	return &domain.PasswordReset{
		ID:         dbReset.ID,
		UserID:     dbReset.UserID,
		TokenHash:  dbReset.TokenHash,
		HotpSecret: dbReset.HotpSecret,
		Counter:    dbReset.Counter.Int64,
		ExpiresAt:  dbReset.ExpiresAt,
		CreatedAt:  dbReset.CreatedAt,
		UsedAt:     dbReset.UsedAt,
	}
}
