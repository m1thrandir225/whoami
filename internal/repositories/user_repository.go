// Package repositories
package repositories

import (
	"context"
	"encoding/json"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	UpdateUserPrivacySettings(ctx context.Context, id int64, privacySettings domain.PrivacySettings) error
	DeactivateUser(ctx context.Context, id int64) error
	ActivateUser(ctx context.Context, id int64) error
	UpdateLastLogin(ctx context.Context, id int64) error
}

type userRepository struct {
	store db.Store
}

func NewUserRepository(store db.Store) UserRepository {
	return &userRepository{
		store: store,
	}
}

func (repo *userRepository) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	privacySettingsJson, err := json.Marshal(req.PrivacySettings)
	if err != nil {
		return nil, err
	}

	user, err := repo.store.CreateUser(ctx, db.CreateUserParams{
		Email:           req.Email,
		PasswordHash:    req.Password,
		Role:            string(domain.RoleUser),
		PrivacySettings: privacySettingsJson,
	})
	if err != nil {
		return nil, err
	}
	var privacySettings domain.PrivacySettings
	err = json.Unmarshal(user.PrivacySettings, &privacySettings)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:                user.ID,
		Email:             user.Email,
		EmailVerified:     user.EmailVerified,
		Password:          user.PasswordHash,
		PasswordChangedAt: user.PasswordChangedAt,
		Role:              user.Role,
		Active:            user.Active,
		PrivacySettings:   privacySettings,
		LastLoginAt:       user.LastLoginAt,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}, nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return &domain.User{}, nil
}

func (repo *userRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return &domain.User{}, nil
}

func (repo *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return nil
}

func (repo *userRepository) UpdateUserPrivacySettings(ctx context.Context, id int64, privacySettings domain.PrivacySettings) error {
	return nil
}

func (repo *userRepository) DeactivateUser(ctx context.Context, id int64) error {
	return nil
}

func (repo *userRepository) ActivateUser(ctx context.Context, id int64) error {
	return nil
}

func (repo *userRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	return nil
}
