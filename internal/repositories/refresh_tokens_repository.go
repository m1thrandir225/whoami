package repositories

import (
	"context"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type RefreshTokensRepository interface {
	CreateRefreshToken(ctx context.Context, req domain.CreateRefreshTokenAction) (*domain.RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	GetActiveRefreshTokensByUser(ctx context.Context, userID int64) ([]domain.RefreshToken, error)
	RevokeAllUserRefreshTokens(ctx context.Context, userID int64) error
	UpdateRefreshTokenLastUsed(ctx context.Context, tokenHash string) error
	CleanupExpiredRefreshTokens(ctx context.Context) error
}

type refreshTokensRepository struct {
	store db.Store
}

func NewRefreshTokensRepository(store db.Store) RefreshTokensRepository {
	return &refreshTokensRepository{
		store: store,
	}
}

func (repo *refreshTokensRepository) toDomain(dbToken db.RefreshToken) *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:        dbToken.ID,
		UserID:    dbToken.UserID,
		Token:     dbToken.TokenHash,
		CreatedAt: dbToken.CreatedAt,
		ExpiresAt: dbToken.ExpiresAt,
		RevokedAt: dbToken.RevokedAt,
	}
}

func (repo *refreshTokensRepository) CreateRefreshToken(ctx context.Context, req domain.CreateRefreshTokenAction) (*domain.RefreshToken, error) {
	token, err := repo.store.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:     req.UserID,
		ExpiresAt:  req.ExpiresAt,
		TokenHash:  req.Token,
		DeviceInfo: req.DeviceInfo,
	})
	if err != nil {
		return nil, err
	}

	return repo.toDomain(token), nil
}

func (repo *refreshTokensRepository) CleanupExpiredRefreshTokens(ctx context.Context) error {
	return repo.store.CleanupExpiredRefreshTokens(ctx)
}

func (repo *refreshTokensRepository) GetActiveRefreshTokensByUser(ctx context.Context, userID int64) ([]domain.RefreshToken, error) {
	tokens, err := repo.store.GetActiveRefreshTokensByUser(ctx, userID)
	if err != nil {
		return make([]domain.RefreshToken, 0), err
	}

	domainTokens := make([]domain.RefreshToken, len(tokens))

	for i, token := range tokens {
		domainTokens[i] = *repo.toDomain(token)
	}

	return domainTokens, nil
}

func (repo *refreshTokensRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	token, err := repo.store.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	return repo.toDomain(token), nil
}

func (repo *refreshTokensRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID int64) error {
	return repo.store.RevokeAllUserRefreshTokens(ctx, userID)
}

func (repo *refreshTokensRepository) UpdateRefreshTokenLastUsed(ctx context.Context, tokenHash string) error {
	return repo.store.UpdateRefreshTokenLastUsed(ctx, tokenHash)
}
