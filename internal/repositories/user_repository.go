// Package repositories
package repositories

import (
	"context"
	"encoding/json"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	GetUserByID(ctx context.Context, id int64) (*db.User, error)
	UpdateUser(ctx context.Context, user *db.User) error
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

func (repo *userRepository) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*db.User, error) {
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

	return &user, nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := repo.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) GetUserByID(ctx context.Context, id int64) (*db.User, error) {
	return &db.User{}, nil
}

func (repo *userRepository) UpdateUser(ctx context.Context, user *db.User) error {
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
