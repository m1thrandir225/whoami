// Package repositories
package repositories

import (
	"context"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req any) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	UpdateUserPrivacySettings(ctx context.Context, id int64, privacySettings domain.PrivacySettings) error
	DeactivateUser(ctx context.Context, id int64) error
	ActivateUser(ctx context.Context, id int64) error
	UpdateLastLogin(ctx context.Context, id int64) error
}

type userRepository struct {
	store *db.Store
}

func NewUserRepository(store *db.Store) UserRepository {
	return &userRepository{
		store: store,
	}
}

func (repo *userRepository) CreateUser(ctx context.Context, req any) (*domain.User, error) {
	return &domain.User{}, nil
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
