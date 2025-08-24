package services

import (
	"context"
	"fmt"
	"time"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/repositories"
	"github.com/m1thrandir225/whoami/internal/security"
)

type UserDevicesService interface {
	RegisterDevice(ctx context.Context, userID int64, deviceInfo *security.DeviceInfo) (*domain.UserDevice, error)
	GetUserDevices(ctx context.Context, userID int64) ([]domain.UserDevice, error)
	GetUserDevice(ctx context.Context, id, userID int64) (*domain.UserDevice, error)
	UpdateDeviceLastUsed(ctx context.Context, id int64) (*domain.UserDevice, error)
	UpdateDevice(ctx context.Context, req domain.UpdateUserDeviceAction) (*domain.UserDevice, error)
	DeleteDevice(ctx context.Context, id, userID int64) error
	DeleteAllDevices(ctx context.Context, userID int64) error
	GetOrCreateDevice(ctx context.Context, userID int64, deviceInfo *security.DeviceInfo) (*domain.UserDevice, error)
	MarkDeviceAsTrusted(ctx context.Context, id, userID int64, trusted bool) (*domain.UserDevice, error)
}

type userDevicesService struct {
	userDevicesRepo repositories.UserDevicesRepository
}

func NewUserDevicesService(userDevicesRepo repositories.UserDevicesRepository) UserDevicesService {
	return &userDevicesService{
		userDevicesRepo: userDevicesRepo,
	}
}

func (s *userDevicesService) RegisterDevice(ctx context.Context, userID int64, deviceInfo *security.DeviceInfo) (*domain.UserDevice, error) {
	device, err := s.userDevicesRepo.CreateUserDevice(ctx, domain.CreateUserDeviceAction{
		UserID:     userID,
		DeviceID:   deviceInfo.DeviceID,
		DeviceName: deviceInfo.DeviceName,
		DeviceType: deviceInfo.DeviceType,
		UserAgent:  deviceInfo.UserAgent,
		IPAddress:  deviceInfo.IPAddress,
		Trusted:    false, // Default to untrusted
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}

	return device, nil
}

func (s *userDevicesService) GetUserDevices(ctx context.Context, userID int64) ([]domain.UserDevice, error) {
	return s.userDevicesRepo.GetUserDevicesByUserID(ctx, userID)
}

func (s *userDevicesService) GetUserDevice(ctx context.Context, id, userID int64) (*domain.UserDevice, error) {
	return s.userDevicesRepo.GetUserDeviceByID(ctx, id, userID)
}

func (s *userDevicesService) UpdateDeviceLastUsed(ctx context.Context, id int64) (*domain.UserDevice, error) {
	return s.userDevicesRepo.UpdateUserDeviceLastUsed(ctx, id, time.Now())
}

func (s *userDevicesService) UpdateDevice(ctx context.Context, req domain.UpdateUserDeviceAction) (*domain.UserDevice, error) {
	return s.userDevicesRepo.UpdateUserDevice(ctx, req)
}

func (s *userDevicesService) DeleteDevice(ctx context.Context, id, userID int64) error {
	return s.userDevicesRepo.DeleteUserDevice(ctx, id, userID)
}

func (s *userDevicesService) DeleteAllDevices(ctx context.Context, userID int64) error {
	return s.userDevicesRepo.DeleteAllUserDevices(ctx, userID)
}

func (s *userDevicesService) GetOrCreateDevice(ctx context.Context, userID int64, deviceInfo *security.DeviceInfo) (*domain.UserDevice, error) {
	// Try to find existing device
	existingDevice, err := s.userDevicesRepo.GetUserDeviceByDeviceID(ctx, userID, deviceInfo.DeviceID)
	if err == nil {
		// Device exists, update last used
		return s.UpdateDeviceLastUsed(ctx, existingDevice.ID)
	}

	// Device doesn't exist, create new one
	return s.RegisterDevice(ctx, userID, deviceInfo)
}

func (s *userDevicesService) MarkDeviceAsTrusted(ctx context.Context, id, userID int64, trusted bool) (*domain.UserDevice, error) {
	return s.userDevicesRepo.MarkDeviceAsTrusted(ctx, id, userID, trusted)
}
