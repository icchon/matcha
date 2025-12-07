package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/service"
	"github.com/icchon/matcha/api/internal/presentation/helper"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
)

type ChatHandler struct {
	chatSvc service.ChatService
}

func NewChatHandler(chatSvc service.ChatService) *ChatHandler {
	return &ChatHandler{chatSvc: chatSvc}
}

// MessageResponse is the DTO for a chat message.
type MessageResponse struct {
	ID          int64     `json:"id"`
	SenderID    uuid.UUID `json:"sender_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Content     string    `json:"content"`
	SentAt      time.Time `json:"sent_at"`
	IsRead      *bool     `json:"is_read,omitempty"`
}

func newMessageResponse(m *entity.Message) *MessageResponse {
	if m == nil {
		return nil
	}
	res := &MessageResponse{
		ID:          m.ID,
		SenderID:    m.SenderID,
		RecipientID: m.RecipientID,
		Content:     m.Content,
		SentAt:      m.SentAt,
	}
	if m.IsRead.Valid {
		res.IsRead = &m.IsRead.Bool
	}
	return res
}

// ChatResponse is the DTO for a chat list item.
type ChatResponse struct {
	OtherUser   *ProfileResponse `json:"other_user"`
	LastMessage *MessageResponse `json:"last_message"`
}

func newChatResponse(c *service.Chat) *ChatResponse {
	if c == nil {
		return nil
	}
	return &ChatResponse{
		OtherUser:   newProfileResponse(&c.OtherUser),
		LastMessage: newMessageResponse(c.LastMessage),
	}
}

func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	chats, err := h.chatSvc.GetChatsForUser(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}

	res := make([]*ChatResponse, len(chats))
	for i, c := range chats {
		res[i] = newChatResponse(c)
	}

	helper.RespondWithJSON(w, http.StatusOK, res)
}

func (h *ChatHandler) GetChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	selfID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	otherIDStr := chi.URLParam(r, "userID")
	otherID, err := uuid.Parse(otherIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20 // default limit
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			helper.HandleError(w, apperrors.ErrInvalidInput)
			return
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0 // default offset
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			helper.HandleError(w, apperrors.ErrInvalidInput)
			return
		}
	}

	params := &service.GetChatMessagesParams{
		UserID1: selfID,
		UserID2: otherID,
		Limit:   limit,
		Offset:  offset,
	}

	messages, err := h.chatSvc.GetChatMessages(r.Context(), params)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}

	res := make([]*MessageResponse, len(messages))
	for i, m := range messages {
		res[i] = newMessageResponse(m)
	}

	helper.RespondWithJSON(w, http.StatusOK, res)
}
