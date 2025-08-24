package repositories

import (
	"context"
	"net/netip"
	"time"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type AccountLockoutRepository interface {
	CreateLockout(ctx context.Context, req domain.CreateAccountLockoutAction) (*domain.AccountLockout, error)
	GetLockoutByUserID(ctx context.Context, userID int64) (*domain.AccountLockout, error)
	GetLockoutByIP(ctx context.Context, ipAddress string) (*domain.AccountLockout, error)
	GetLockoutByUserAndIP(ctx context.Context, userID int64, ipAddress string) (*domain.AccountLockout, error)
	DeleteExpiredLockouts(ctx context.Context) error
	DeleteLockoutByID(ctx context.Context, id int64) error
}

type accountLockoutRepository struct {
	store db.Store
}

func NewAccountLockoutRepository(store db.Store) AccountLockoutRepository {
	return &accountLockoutRepository{
		store: store,
	}
}
func (r *accountLockoutRepository) CreateLockout(ctx context.Context, req domain.CreateAccountLockoutAction) (*domain.AccountLockout, error) {
	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		return nil, err
	}

	parsedIP, err := netip.ParseAddr(req.IPAddress)
	if err != nil {
		return nil, err
	}

	dbLockout, err := r.store.CreateAccountLockout(ctx, db.CreateAccountLockoutParams{
		UserID:      req.UserID,
		IpAddress:   &parsedIP,
		LockoutType: string(req.LockoutType),
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbLockout), nil
}

func (r *accountLockoutRepository) GetLockoutByUserID(ctx context.Context, userID int64) (*domain.AccountLockout, error) {
	dbLockout, err := r.store.GetAccountLockoutByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbLockout), nil
}

func (r *accountLockoutRepository) GetLockoutByIP(ctx context.Context, ipAddress string) (*domain.AccountLockout, error) {
	parsedIP, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil, err
	}
	dbLockout, err := r.store.GetAccountLockoutByIP(ctx, &parsedIP)
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbLockout), nil
}

func (r *accountLockoutRepository) GetLockoutByUserAndIP(ctx context.Context, userID int64, ipAddress string) (*domain.AccountLockout, error) {
	parsedIP, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return nil, err
	}
	dbLockout, err := r.store.GetAccountLockoutByUserAndIP(ctx, db.GetAccountLockoutByUserAndIPParams{
		UserID:    userID,
		IpAddress: &parsedIP,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbLockout), nil
}

func (r *accountLockoutRepository) DeleteExpiredLockouts(ctx context.Context) error {
	return r.store.DeleteExpiredLockouts(ctx)
}

func (r *accountLockoutRepository) DeleteLockoutByID(ctx context.Context, id int64) error {
	return r.store.DeleteAccountLockoutByID(ctx, id)
}

func (r *accountLockoutRepository) toDomain(dbLockout db.AccountLockout) *domain.AccountLockout {
	return &domain.AccountLockout{
		ID:          dbLockout.ID,
		UserID:      dbLockout.UserID,
		IPAddress:   dbLockout.IpAddress.String(),
		LockoutType: dbLockout.LockoutType,
		ExpiresAt:   dbLockout.ExpiresAt,
		CreatedAt:   dbLockout.CreatedAt,
	}
}
