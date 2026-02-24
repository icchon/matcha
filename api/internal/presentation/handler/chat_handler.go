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

type ChatHandler struct {
	chatSvc service.ChatService
}

func NewChatHandler(chatSvc service.ChatService) *ChatHandler {
	return &ChatHandler{chatSvc: chatSvc}
}

func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized) // Use helper.HandleError
		return
	}

	chats, err := h.chatSvc.GetChatsForUser(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInternalServer) // Use helper.HandleError
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, chats) // Use helper.RespondWithJSON
}

func (h *ChatHandler) GetChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// get self id from context
	selfID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized) // Use helper.HandleError
		return
	}

	otherIDStr := chi.URLParam(r, "userID")
	otherID, err := uuid.Parse(otherIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput) // Use helper.HandleError
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20 // default limit
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			helper.HandleError(w, apperrors.ErrInvalidInput) // Use helper.HandleError
			return
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0 // default offset
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			helper.HandleError(w, apperrors.ErrInvalidInput) // Use helper.HandleError
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
		helper.HandleError(w, apperrors.ErrInternalServer) // Use helper.HandleError
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, messages) // Use helper.RespondWithJSON
}
