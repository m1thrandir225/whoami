// Package services
package services

import "github.com/m1thrandir225/whoami/internal/repositories"

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}
