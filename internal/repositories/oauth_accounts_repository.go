package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type OAuthAccountsRepository interface {
	CreateOAuthAccount(ctx context.Context, req domain.CreateOAuthAccountAction) (*domain.OAuthAccount, error)
	GetOAuthAccountByID(ctx context.Context, id, userID int64) (*domain.OAuthAccount, error)
	GetOAuthAccountByProvider(ctx context.Context, provider, providerUserID string) (*domain.OAuthAccount, error)
	GetOAuthAccountsByUserID(ctx context.Context, userID int64) ([]domain.OAuthAccount, error)
	GetOAuthAccountByEmail(ctx context.Context, email, provider string) (*domain.OAuthAccount, error)
	UpdateOAuthAccount(ctx context.Context, req domain.UpdateOAuthAccountAction) (*domain.OAuthAccount, error)
	DeleteOAuthAccount(ctx context.Context, id, userID int64) error
	DeleteOAuthAccountByProvider(ctx context.Context, userID int64, provider string) error
	UpdateOAuthTokens(ctx context.Context, req domain.UpdateOAuthTokensAction) (*domain.OAuthAccount, error)
}

type oauthAccountsRepository struct {
	store db.Store
}

func NewOAuthAccountsRepository(store db.Store) OAuthAccountsRepository {
	return &oauthAccountsRepository{
		store: store,
	}
}

func (r *oauthAccountsRepository) CreateOAuthAccount(ctx context.Context, req domain.CreateOAuthAccountAction) (*domain.OAuthAccount, error) {
	var email, name, avatarURL, accessToken, refreshToken pgtype.Text
	var tokenExpiresAt pgtype.Timestamptz

	if req.Email != nil {
		email = pgtype.Text{String: *req.Email, Valid: true}
	}
	if req.Name != nil {
		name = pgtype.Text{String: *req.Name, Valid: true}
	}
	if req.AvatarURL != nil {
		avatarURL = pgtype.Text{String: *req.AvatarURL, Valid: true}
	}
	if req.AccessToken != nil {
		accessToken = pgtype.Text{String: *req.AccessToken, Valid: true}
	}
	if req.RefreshToken != nil {
		refreshToken = pgtype.Text{String: *req.RefreshToken, Valid: true}
	}
	if req.TokenExpiresAt != nil {
		tokenExpiresAt = pgtype.Timestamptz{Time: *req.TokenExpiresAt, Valid: true}
	}

	dbAccount, err := r.store.CreateOAuthAccount(ctx, db.CreateOAuthAccountParams{
		UserID:         req.UserID,
		Provider:       req.Provider,
		ProviderUserID: req.ProviderUserID,
		Email:          email,
		Name:           name,
		AvatarUrl:      avatarURL,
		AccessToken:    accessToken.String,
		RefreshToken:   refreshToken.String,
		TokenExpiresAt: &tokenExpiresAt.Time,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbAccount), nil
}

func (r *oauthAccountsRepository) GetOAuthAccountByID(ctx context.Context, id, userID int64) (*domain.OAuthAccount, error) {
	dbAccount, err := r.store.GetOAuthAccountByID(ctx, db.GetOAuthAccountByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbAccount), nil
}

func (r *oauthAccountsRepository) GetOAuthAccountByProvider(ctx context.Context, provider, providerUserID string) (*domain.OAuthAccount, error) {
	dbAccount, err := r.store.GetOAuthAccountByProvider(ctx, db.GetOAuthAccountByProviderParams{
		Provider:       provider,
		ProviderUserID: providerUserID,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbAccount), nil
}

func (r *oauthAccountsRepository) GetOAuthAccountsByUserID(ctx context.Context, userID int64) ([]domain.OAuthAccount, error) {
	dbAccounts, err := r.store.GetOAuthAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	accounts := make([]domain.OAuthAccount, len(dbAccounts))
	for i, account := range dbAccounts {
		accounts[i] = *r.toDomain(account)
	}

	return accounts, nil
}

func (r *oauthAccountsRepository) GetOAuthAccountByEmail(ctx context.Context, email, provider string) (*domain.OAuthAccount, error) {
	var pgEmail pgtype.Text
	if email != "" {
		pgEmail = pgtype.Text{String: email, Valid: true}
	}

	dbAccount, err := r.store.GetOAuthAccountByEmail(ctx, db.GetOAuthAccountByEmailParams{
		Email:    pgEmail,
		Provider: provider,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbAccount), nil
}

func (r *oauthAccountsRepository) UpdateOAuthAccount(ctx context.Context, req domain.UpdateOAuthAccountAction) (*domain.OAuthAccount, error) {
	var email, name, avatarURL, accessToken, refreshToken pgtype.Text
	var tokenExpiresAt pgtype.Timestamptz

	if req.Email != nil {
		email = pgtype.Text{String: *req.Email, Valid: true}
	}
	if req.Name != nil {
		name = pgtype.Text{String: *req.Name, Valid: true}
	}
	if req.AvatarURL != nil {
		avatarURL = pgtype.Text{String: *req.AvatarURL, Valid: true}
	}
	if req.AccessToken != nil {
		accessToken = pgtype.Text{String: *req.AccessToken, Valid: true}
	}
	if req.RefreshToken != nil {
		refreshToken = pgtype.Text{String: *req.RefreshToken, Valid: true}
	}
	if req.TokenExpiresAt != nil {
		tokenExpiresAt = pgtype.Timestamptz{Time: *req.TokenExpiresAt, Valid: true}
	}

	dbAccount, err := r.store.UpdateOAuthAccount(ctx, db.UpdateOAuthAccountParams{
		ID:             req.ID,
		UserID:         req.UserID,
		Email:          email,
		Name:           name,
		AvatarUrl:      avatarURL,
		AccessToken:    accessToken.String,
		RefreshToken:   refreshToken.String,
		TokenExpiresAt: &tokenExpiresAt.Time,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbAccount), nil
}

func (r *oauthAccountsRepository) DeleteOAuthAccount(ctx context.Context, id, userID int64) error {
	return r.store.DeleteOAuthAccount(ctx, db.DeleteOAuthAccountParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *oauthAccountsRepository) DeleteOAuthAccountByProvider(ctx context.Context, userID int64, provider string) error {
	return r.store.DeleteOAuthAccountByProvider(ctx, db.DeleteOAuthAccountByProviderParams{
		UserID:   userID,
		Provider: provider,
	})
}

func (r *oauthAccountsRepository) UpdateOAuthTokens(ctx context.Context, req domain.UpdateOAuthTokensAction) (*domain.OAuthAccount, error) {
	var accessToken, refreshToken pgtype.Text
	var tokenExpiresAt pgtype.Timestamptz

	if req.AccessToken != nil {
		accessToken = pgtype.Text{String: *req.AccessToken, Valid: true}
	}
	if req.RefreshToken != nil {
		refreshToken = pgtype.Text{String: *req.RefreshToken, Valid: true}
	}
	if req.TokenExpiresAt != nil {
		tokenExpiresAt = pgtype.Timestamptz{Time: *req.TokenExpiresAt, Valid: true}
	}

	dbAccount, err := r.store.UpdateOAuthTokens(ctx, db.UpdateOAuthTokensParams{
		ID:             req.ID,
		UserID:         req.UserID,
		AccessToken:    accessToken.String,
		RefreshToken:   refreshToken.String,
		TokenExpiresAt: &tokenExpiresAt.Time,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbAccount), nil
}

func (r *oauthAccountsRepository) toDomain(dbAccount db.OauthAccount) *domain.OAuthAccount {
	account := &domain.OAuthAccount{
		ID:             dbAccount.ID,
		UserID:         dbAccount.UserID,
		Provider:       dbAccount.Provider,
		ProviderUserID: dbAccount.ProviderUserID,
		CreatedAt:      dbAccount.CreatedAt,
		UpdatedAt:      dbAccount.UpdatedAt,
	}

	if dbAccount.Email.Valid {
		account.Email = &dbAccount.Email.String
	}
	if dbAccount.Name.Valid {
		account.Name = &dbAccount.Name.String
	}
	if dbAccount.AvatarUrl.Valid {
		account.AvatarURL = &dbAccount.AvatarUrl.String
	}
	account.AccessToken = &dbAccount.AccessToken
	account.RefreshToken = &dbAccount.RefreshToken
	account.TokenExpiresAt = dbAccount.TokenExpiresAt

	return account
}
