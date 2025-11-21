package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
	"github.com/icchon/matcha/api/internal/infrastructure/uow"
)

type userService struct {
	uow      uow.UnitOfWork
	userRepo repo.UserRepository
}

var _ service.UserService = (*userService)(nil)

func NewUserService(uow uow.UnitOfWork, userRepo repo.UserRepository) *userService {
	return &userService{
		uow:      uow,
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return nil, nil
}
