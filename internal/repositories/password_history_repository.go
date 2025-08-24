package repositories

import (
	"context"
	"time"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type PasswordHistoryRepository interface {
	CreatePasswordHistory(ctx context.Context, req domain.CreatePasswordHistory) error
	GetPasswordHistory(ctx context.Context, userID int64, limit int32) ([]domain.PasswordHistory, error)
	DeleteOldPasswordHistory(ctx context.Context, userID int64) error
	CheckPasswordInHistory(ctx context.Context, userID int64, passwordHash string) (bool, error)
}

type passwordHistoryRepository struct {
	store db.Store
}

func NewPasswordHistoryRepository(store db.Store) PasswordHistoryRepository {
	return &passwordHistoryRepository{
		store: store,
	}
}

func (r *passwordHistoryRepository) CreatePasswordHistory(ctx context.Context, req domain.CreatePasswordHistory) error {
	now := time.Now()
	_, err := r.store.CreatePasswordHistory(ctx, db.CreatePasswordHistoryParams{
		UserID:       req.UserID,
		PasswordHash: req.PasswordHash,
		CreatedAt:    &now,
	})
	return err
}

func (r *passwordHistoryRepository) GetPasswordHistory(ctx context.Context, userID int64, limit int32) ([]domain.PasswordHistory, error) {
	dbHistory, err := r.store.GetPasswordHistoryByUserID(ctx, db.GetPasswordHistoryByUserIDParams{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	history := make([]domain.PasswordHistory, len(dbHistory))
	for i, item := range dbHistory {
		history[i] = *r.toDomain(item)
	}

	return history, nil
}

func (r *passwordHistoryRepository) DeleteOldPasswordHistory(ctx context.Context, userID int64) error {
	return r.store.DeleteOldPasswordHistory(ctx, userID)
}

func (r *passwordHistoryRepository) CheckPasswordInHistory(ctx context.Context, userID int64, passwordHash string) (bool, error) {
	count, err := r.store.CheckPasswordInHistory(ctx, db.CheckPasswordInHistoryParams{
		UserID:       userID,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *passwordHistoryRepository) toDomain(dbHistory db.PasswordHistory) *domain.PasswordHistory {
	return &domain.PasswordHistory{
		ID:           dbHistory.ID,
		UserID:       dbHistory.UserID,
		PasswordHash: dbHistory.PasswordHash,
		CreatedAt:    *dbHistory.CreatedAt,
	}
}
