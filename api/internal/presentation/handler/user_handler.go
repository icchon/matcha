package handler

import (
	// "net/http"

	// "github.com/go-chi/chi/v5"
	// "github.com/google/uuid"
	// "github.com/icchon/matcha/api/internal/apperrors"
	// "github.com/icchon/matcha/api/internal/presentation/helper"
	// "github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
