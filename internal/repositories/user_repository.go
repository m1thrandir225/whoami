// Package repositories
package repositories

import (
	"context"
	"encoding/json"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req domain.CreateUserAction) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	MarkEmailVerified(ctx context.Context, id int64) error
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

func (r *userRepository) toDomain(dbUser db.User, privacySettings domain.PrivacySettings) *domain.User {
	return &domain.User{
		ID:                dbUser.ID,
		Email:             dbUser.Email,
		EmailVerified:     dbUser.EmailVerified,
		Password:          dbUser.PasswordHash,
		PasswordChangedAt: dbUser.PasswordChangedAt,
		Role:              dbUser.Role,
		Active:            dbUser.Active,
		PrivacySettings:   privacySettings,
		LastLoginAt:       dbUser.LastLoginAt,
		CreatedAt:         dbUser.CreatedAt,
		UpdatedAt:         dbUser.UpdatedAt,
	}
}

func (repo *userRepository) CreateUser(ctx context.Context, req domain.CreateUserAction) (*domain.User, error) {
	privacySettingsJSON, err := json.Marshal(req.PrivacySettings)
	if err != nil {
		return nil, err
	}

	user, err := repo.store.CreateUser(ctx, db.CreateUserParams{
		Email:           req.Email,
		PasswordHash:    req.Password,
		Role:            string(domain.RoleUser),
		PrivacySettings: privacySettingsJSON,
	})
	if err != nil {
		return nil, err
	}

	var privacySettings domain.PrivacySettings
	err = json.Unmarshal(user.PrivacySettings, &privacySettings)
	if err != nil {
		return nil, err
	}

	return repo.toDomain(user, privacySettings), nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := repo.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	var privacySettings domain.PrivacySettings
	err = json.Unmarshal(user.PrivacySettings, &privacySettings)
	if err != nil {
		return nil, err
	}

	return repo.toDomain(user, privacySettings), nil
}

func (repo *userRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := repo.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	var privacySettings domain.PrivacySettings
	err = json.Unmarshal(user.PrivacySettings, &privacySettings)
	if err != nil {
		return nil, err
	}

	return repo.toDomain(user, privacySettings), nil
}

func (repo *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return repo.store.UpdateUser(ctx, db.UpdateUserParams{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	})
}

func (repo *userRepository) UpdateUserPrivacySettings(ctx context.Context, id int64, privacySettings domain.PrivacySettings) error {
	privacySettingsJSON, err := json.Marshal(privacySettings)
	if err != nil {
		return err
	}

	return repo.store.UpdateUserPrivacySettings(ctx, db.UpdateUserPrivacySettingsParams{
		ID:              id,
		PrivacySettings: privacySettingsJSON,
	})
}

func (repo *userRepository) DeactivateUser(ctx context.Context, id int64) error {
	return repo.store.DeactivateUser(ctx, id)
}

func (repo *userRepository) ActivateUser(ctx context.Context, id int64) error {
	return repo.store.ActivateUser(ctx, id)
}

func (repo *userRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	return repo.store.UpdateLastLogin(ctx, id)
}

func (repo *userRepository) MarkEmailVerified(ctx context.Context, id int64) error {
	return repo.store.MarkEmailVerified(ctx, id)
}
