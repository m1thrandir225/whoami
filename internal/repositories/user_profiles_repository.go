package repositories

import (
	"context"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type UserProfilesRepository interface {
	CreateUserProfile(ctx context.Context, req domain.CreateUserProfileRequest) (*domain.UserProfile, error)
	GetUserProfile(ctx context.Context, userID int64) (*domain.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile domain.UserProfile) error
}

type userProfilesRepository struct {
	store db.Store
}

func NewUserProfilesRepository(store db.Store) UserProfilesRepository {
	return &userProfilesRepository{
		store: store,
	}
}

func (repo *userProfilesRepository) CreateUserProfile(ctx context.Context, req domain.CreateUserProfileRequest) (*domain.UserProfile, error) {
	panic("unimplemented")
}

func (repo *userProfilesRepository) GetUserProfile(ctx context.Context, userID int64) (*domain.UserProfile, error) {
	panic("unimplemented")
}

func (repo *userProfilesRepository) UpdateUserProfile(ctx context.Context, profile domain.UserProfile) error {
	panic("unimplemented")
}
