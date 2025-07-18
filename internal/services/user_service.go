// Package services
package services

import (
	"context"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, req domain.CreateUserAction) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	ActivateUser(ctx context.Context, id int64) error
	DeactivateUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, user domain.User) error
	UpdateUserPrivacySettings(ctx context.Context, id int64, privacySettings domain.PrivacySettings) error
	UpdateLastLogin(ctx context.Context, id int64) error
}

type userService struct {
	repository repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repository: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req domain.CreateUserAction) (*domain.User, error) {
	return s.repository.CreateUser(ctx, req)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repository.GetUserByEmail(ctx, email)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.repository.GetUserByID(ctx, id)
}

func (s *userService) ActivateUser(ctx context.Context, id int64) error {
	return s.repository.ActivateUser(ctx, id)
}

func (s *userService) DeactivateUser(ctx context.Context, id int64) error {
	return s.repository.DeactivateUser(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user domain.User) error {
	return s.repository.UpdateUser(ctx, &user)
}

func (s *userService) UpdateUserPrivacySettings(ctx context.Context, id int64, privacySettings domain.PrivacySettings) error {
	return s.repository.UpdateUserPrivacySettings(ctx, id, privacySettings)
}

func (s *userService) UpdateLastLogin(ctx context.Context, id int64) error {
	return s.repository.UpdateLastLogin(ctx, id)
}
