package user

import (
	"github.com/icchon/matcha/api/internal/domain"
	"github.com/icchon/matcha/api/internal/infrastructure/uow"
)

type UserService interface {
	// GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type userService struct {
	uow      uow.UnitOfWork
	userRepo domain.UserQueryRepository
}

func NewUserService(
	uow uow.UnitOfWork,
	userRepo domain.UserQueryRepository,
) UserService {
	return &userService{
		uow:      uow,
		userRepo: userRepo,
	}
}
