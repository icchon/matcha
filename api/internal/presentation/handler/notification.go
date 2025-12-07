package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
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

// NotificationResponse is the DTO for notification responses.
type NotificationResponse struct {
	ID          int64     `json:"id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	SenderID    *string   `json:"sender_id,omitempty"`
	Type        string    `json:"type"`
	IsRead      *bool     `json:"is_read,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// newNotificationResponse converts an entity.Notification to a NotificationResponse DTO.
func newNotificationResponse(n *entity.Notification) *NotificationResponse {
	if n == nil {
		return nil
	}
	res := &NotificationResponse{
		ID:          n.ID,
		RecipientID: n.RecipientID,
		Type:        string(n.Type),
		CreatedAt:   n.CreatedAt,
	}
	if n.SenderID.Valid {
		res.SenderID = &n.SenderID.String
	}
	if n.IsRead.Valid {
		res.IsRead = &n.IsRead.Bool
	}
	return res
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}

	notifications, err := h.notifSvc.GetNotifications(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}

	res := make([]*NotificationResponse, len(notifications))
	for i, n := range notifications {
		res[i] = newNotificationResponse(n)
	}

	helper.RespondWithJSON(w, http.StatusOK, res)
}
