package handler

import (
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/presentation/helper"
	"github.com/icchon/matcha/api/internal/service/user"
	"net/http"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) FindUserHandler(w http.ResponseWriter, r *http.Request) {
	helper.HandleError(w, apperrors.ErrNotImplemented)
}
