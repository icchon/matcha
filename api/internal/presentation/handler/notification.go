package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/service"
	"github.com/icchon/matcha/api/internal/presentation/helper" // Add this import
)

type NotificationHandler struct {
	notifSvc service.NotificationService
}

func NewNotificationHandler(notifSvc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifSvc: notifSvc}
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput) // Use helper.HandleError
		return
	}

	notifications, err := h.notifSvc.GetNotifications(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInternalServer) // Use helper.HandleError
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, notifications) // Use helper.RespondWithJSON
}
