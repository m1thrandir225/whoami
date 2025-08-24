package repositories

import (
	"context"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type UserDevicesRepository interface {
	CreateUserDevice(ctx context.Context, req domain.CreateUserDeviceAction) (*domain.UserDevice, error)
	GetUserDevicesByUserID(ctx context.Context, userID int64) ([]domain.UserDevice, error)
	GetUserDeviceByID(ctx context.Context, id, userID int64) (*domain.UserDevice, error)
	UpdateUserDeviceLastUsed(ctx context.Context, id int64, lastUsedAt time.Time) (*domain.UserDevice, error)
	DeleteUserDevice(ctx context.Context, id, userID int64) error
	DeleteAllUserDevices(ctx context.Context, userID int64) error
	GetUserDeviceByDeviceID(ctx context.Context, userID int64, deviceID string) (*domain.UserDevice, error)
	UpdateUserDevice(ctx context.Context, req domain.UpdateUserDeviceAction) (*domain.UserDevice, error)
	MarkDeviceAsTrusted(ctx context.Context, id, userID int64, trusted bool) (*domain.UserDevice, error)
}

type userDevicesRepository struct {
	store db.Store
}

func NewUserDevicesRepository(store db.Store) UserDevicesRepository {
	return &userDevicesRepository{
		store: store,
	}
}

func (r *userDevicesRepository) CreateUserDevice(ctx context.Context, req domain.CreateUserDeviceAction) (*domain.UserDevice, error) {
	parsedIP, err := netip.ParseAddr(req.IPAddress)
	if err != nil {
		return nil, err
	}
	var deviceName pgtype.Text
	if req.DeviceName == "" {
		deviceName = pgtype.Text{String: "Unknown", Valid: true}
	} else {
		deviceName = pgtype.Text{String: req.DeviceName, Valid: true}
	}

	var deviceType pgtype.Text
	if req.DeviceType == "" {
		deviceType = pgtype.Text{String: "Unknown", Valid: true}
	} else {
		deviceType = pgtype.Text{String: req.DeviceType, Valid: true}
	}

	var trusted pgtype.Bool
	if req.Trusted {
		trusted = pgtype.Bool{Bool: req.Trusted, Valid: true}
	} else {
		trusted = pgtype.Bool{Bool: false, Valid: true}
	}

	dbDevice, err := r.store.CreateUserDevice(ctx, db.CreateUserDeviceParams{
		UserID:     req.UserID,
		DeviceID:   req.DeviceID,
		DeviceName: deviceName,
		DeviceType: deviceType,
		UserAgent:  req.UserAgent,
		IpAddress:  &parsedIP,
		Trusted:    trusted,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbDevice), nil
}

func (r *userDevicesRepository) GetUserDevicesByUserID(ctx context.Context, userID int64) ([]domain.UserDevice, error) {
	dbDevices, err := r.store.GetUserDevicesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	devices := make([]domain.UserDevice, len(dbDevices))
	for i, device := range dbDevices {
		devices[i] = *r.toDomain(device)
	}

	return devices, nil
}

func (r *userDevicesRepository) GetUserDeviceByID(ctx context.Context, id, userID int64) (*domain.UserDevice, error) {
	dbDevice, err := r.store.GetUserDeviceByID(ctx, db.GetUserDeviceByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbDevice), nil
}

func (r *userDevicesRepository) UpdateUserDeviceLastUsed(ctx context.Context, id int64, lastUsedAt time.Time) (*domain.UserDevice, error) {
	dbDevice, err := r.store.UpdateUserDeviceLastUsed(ctx, db.UpdateUserDeviceLastUsedParams{
		ID:         id,
		LastUsedAt: &lastUsedAt,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbDevice), nil
}

func (r *userDevicesRepository) DeleteUserDevice(ctx context.Context, id, userID int64) error {
	return r.store.DeleteUserDevice(ctx, db.DeleteUserDeviceParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *userDevicesRepository) DeleteAllUserDevices(ctx context.Context, userID int64) error {
	return r.store.DeleteAllUserDevices(ctx, userID)
}

func (r *userDevicesRepository) GetUserDeviceByDeviceID(ctx context.Context, userID int64, deviceID string) (*domain.UserDevice, error) {
	dbDevice, err := r.store.GetUserDeviceByDeviceID(ctx, db.GetUserDeviceByDeviceIDParams{
		UserID:   userID,
		DeviceID: deviceID,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbDevice), nil
}

func (r *userDevicesRepository) UpdateUserDevice(ctx context.Context, req domain.UpdateUserDeviceAction) (*domain.UserDevice, error) {
	var deviceName pgtype.Text
	if req.DeviceName == "" {
		deviceName = pgtype.Text{String: "Unknown", Valid: true}
	} else {
		deviceName = pgtype.Text{String: req.DeviceName, Valid: true}
	}

	var deviceType pgtype.Text
	if req.DeviceType == "" {
		deviceType = pgtype.Text{String: "Unknown", Valid: true}
	} else {
		deviceType = pgtype.Text{String: req.DeviceType, Valid: true}
	}

	var trusted pgtype.Bool
	trusted = pgtype.Bool{Bool: req.Trusted, Valid: true}

	dbDevice, err := r.store.UpdateUserDevice(ctx, db.UpdateUserDeviceParams{
		ID:         req.ID,
		UserID:     req.UserID,
		DeviceName: deviceName,
		DeviceType: deviceType,
		UserAgent:  req.UserAgent,
		Trusted:    trusted,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbDevice), nil
}

func (r *userDevicesRepository) MarkDeviceAsTrusted(ctx context.Context, id, userID int64, trusted bool) (*domain.UserDevice, error) {
	var trustedBool pgtype.Bool
	trustedBool = pgtype.Bool{Bool: trusted, Valid: true}

	dbDevice, err := r.store.MarkDeviceAsTrusted(ctx, db.MarkDeviceAsTrustedParams{
		ID:      id,
		UserID:  userID,
		Trusted: trustedBool,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbDevice), nil
}

func (r *userDevicesRepository) toDomain(dbDevice db.UserDevice) *domain.UserDevice {
	return &domain.UserDevice{
		ID:         dbDevice.ID,
		UserID:     dbDevice.UserID,
		DeviceID:   dbDevice.DeviceID,
		DeviceName: dbDevice.DeviceName.String,
		DeviceType: dbDevice.DeviceType.String,
		UserAgent:  dbDevice.UserAgent,
		IPAddress:  dbDevice.IpAddress.String(),
		Trusted:    dbDevice.Trusted.Bool,
		LastUsedAt: *dbDevice.LastUsedAt,
		CreatedAt:  dbDevice.CreatedAt,
	}
}
