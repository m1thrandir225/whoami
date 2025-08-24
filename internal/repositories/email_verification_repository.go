package repositories

import (
	"context"
	"time"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type EmailVerificationRepository interface {
	CreateEmailVerification(ctx context.Context, req domain.CreateEmailVerificationAction) error
	GetEmailVerificationByToken(ctx context.Context, token string) (*domain.EmailVerification, error)
	MarkEmailVerified(ctx context.Context, id int64) error
	DeleteUnverifiedTokens(ctx context.Context, userID int64) error
	GetUnverifiedVerifications(ctx context.Context, userID int64) ([]domain.EmailVerification, error)
}

type emailVerificationRepository struct {
	store db.Store
}

func NewEmailVerificationRepository(store db.Store) EmailVerificationRepository {
	return &emailVerificationRepository{
		store: store,
	}
}

func (r *emailVerificationRepository) CreateEmailVerification(ctx context.Context, req domain.CreateEmailVerificationAction) error {

	_, err := r.store.CreateEmailVerification(ctx, db.CreateEmailVerificationParams{
		UserID:    req.UserID,
		TokenHash: req.Token,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})
	return err
}

func (r *emailVerificationRepository) GetEmailVerificationByToken(ctx context.Context, token string) (*domain.EmailVerification, error) {
	dbVerification, err := r.store.GetEmailVerificationByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbVerification), nil
}

func (r *emailVerificationRepository) MarkEmailVerified(ctx context.Context, id int64) error {
	err := r.store.MarkEmailVerificationAsUsed(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *emailVerificationRepository) DeleteUnverifiedTokens(ctx context.Context, userID int64) error {
	return r.store.DeleteUnverifiedTokens(ctx, userID)
}

func (r *emailVerificationRepository) GetUnverifiedVerifications(ctx context.Context, userID int64) ([]domain.EmailVerification, error) {
	dbVerifications, err := r.store.GetUnverifiedVerifications(ctx, userID)
	if err != nil {
		return nil, err
	}

	verifications := make([]domain.EmailVerification, len(dbVerifications))
	for i, verification := range dbVerifications {
		verifications[i] = *r.toDomain(verification)
	}

	return verifications, nil
}

func (r *emailVerificationRepository) toDomain(dbVerification db.EmailVerification) *domain.EmailVerification {
	return &domain.EmailVerification{
		ID:        dbVerification.ID,
		UserID:    dbVerification.UserID,
		TokenHash: dbVerification.TokenHash,
		ExpiresAt: dbVerification.ExpiresAt,
		CreatedAt: dbVerification.CreatedAt,
		UsedAt:    dbVerification.UsedAt,
	}
}
