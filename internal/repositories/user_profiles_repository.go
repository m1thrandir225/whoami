package repositories

import (
	"context"

	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/domain"
)

type UserProfilesRepository interface {
	CreateUserProfile(ctx context.Context, req domain.CreateUserProfileAction) (*domain.UserProfile, error)
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

func (repo *userProfilesRepository) toDomain(dbProfile db.UserProfile) *domain.UserProfile {
	return &domain.UserProfile{
		ID:        dbProfile.ID,
		UserID:    dbProfile.UserID,
		FirstName: dbProfile.FirstName,
		LastName:  dbProfile.LastName,
		Phone:     dbProfile.Phone,
		AvatarUrl: dbProfile.AvatarUrl,
		Bio:       dbProfile.Bio,
		Timezone:  dbProfile.Timezone,
		Locale:    dbProfile.Locale,
		CreatedAt: dbProfile.CreatedAt,
		UpdatedAt: dbProfile.CreatedAt,
	}
}

func (repo *userProfilesRepository) CreateUserProfile(ctx context.Context, req domain.CreateUserProfileAction) (*domain.UserProfile, error) {
	profile, err := repo.store.CreateUserProfile(ctx, db.CreateUserProfileParams{
		UserID:    req.UserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		AvatarUrl: req.AvatarURL,
		Bio:       req.Bio,
		Timezone:  req.Timezone,
		Locale:    req.Locale,
	})
	if err != nil {
		return nil, err
	}
	return repo.toDomain(profile), nil
}

func (repo *userProfilesRepository) GetUserProfile(ctx context.Context, userID int64) (*domain.UserProfile, error) {
	profile, err := repo.store.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	return repo.toDomain(profile), nil
}

func (repo *userProfilesRepository) UpdateUserProfile(ctx context.Context, profile domain.UserProfile) error {
	return repo.store.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		UserID:    profile.UserID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Phone:     profile.Phone,
		Bio:       profile.Bio,
		AvatarUrl: profile.AvatarUrl,
		Timezone:  profile.Timezone,
		Locale:    profile.Locale,
	})
}
