package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/service"
	"github.com/icchon/matcha/api/internal/presentation/helper"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
)

type NotificationHandler struct {
	notifSvc service.NotificationService
}

func NewNotificationHandler(notifSvc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifSvc: notifSvc}
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	notifications, err := h.notifSvc.GetNotifications(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, notifications)
}

// MarkNotificationAsReadHandler handles PUT /me/notifications/{id}/read
func (h *NotificationHandler) MarkNotificationAsReadHandler(w http.ResponseWriter, r *http.Request) {
	notificationIDStr := chi.URLParam(r, "id")
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}

	recipientID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	if err := h.notifSvc.MarkNotificationAsRead(r.Context(), notificationID, recipientID); err != nil {
		helper.HandleError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusNoContent, nil)
}

// MarkAllNotificationsAsReadHandler handles POST /me/notifications/read
func (h *NotificationHandler) MarkAllNotificationsAsReadHandler(w http.ResponseWriter, r *http.Request) {
	recipientID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	if err := h.notifSvc.MarkAllNotificationsAsRead(r.Context(), recipientID); err != nil {
		helper.HandleError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusNoContent, nil)
}
