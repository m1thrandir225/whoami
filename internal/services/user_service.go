// Package services
package services

import (
	"context"

	"github.com/m1thrandir225/whoami/internal/domain"
	"github.com/m1thrandir225/whoami/internal/repositories"
)

type UserService struct {
	repository repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req domain.CreateUserAction) (*domain.User, error) {
	return s.repository.CreateUser(ctx, req)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repository.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.repository.GetUserByID(ctx, id)
}

func (s *UserService) ActivateUser(ctx context.Context, id int64) error {
	return s.repository.ActivateUser(ctx, id)
}

func (s *UserService) DeactivateUser(ctx context.Context, id int64) error {
	return s.repository.DeactivateUser(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, user domain.User) error {
	return s.repository.UpdateUser(ctx, &user)
}

func (s *UserService) UpdateUserPrivacySettings(ctx context.Context, id int64, privacySettings domain.PrivacySettings) error {
	return s.repository.UpdateUserPrivacySettings(ctx, id, privacySettings)
}

func (s *UserService) UpdateLastLogin(ctx context.Context, id int64) error {
	return s.repository.UpdateLastLogin(ctx, id)
}
