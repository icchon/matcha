package auth

import (
	"github.com/icchon/matcha/api/internal/domain"
	"github.com/icchon/matcha/api/internal/infrastructure/uow"
)

type AuthService interface {
}

type authService struct {
	uow      uow.UnitOfWork
	authRepo domain.AuthQueryRepository
}

func NewAuthService(
	uow uow.UnitOfWork,
	authRepo domain.AuthQueryRepository) AuthService {
	return &authService{
		uow:      uow,
		authRepo: authRepo,
	}
}
