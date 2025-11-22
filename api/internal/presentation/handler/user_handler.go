package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/presentation/helper"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/service"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type UserHandler struct {
	userService   service.UserService
	profileService service.ProfileService
}

func NewUserHandler(userService service.UserService, profileService service.ProfileService) *UserHandler {
	return &UserHandler{userService: userService, profileService: profileService}
}


type LikeUserResponse struct{
	Connection *entity.Connection `json:"connection"`
	Message string `json:"message"`
}
// /users/{userID}/like POST
func (h *UserHandler) LikeUserHandler(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	likedIDStr := chi.URLParam(r, string(helper.UserIDUrlParam))
	likedID, err := uuid.Parse(likedIDStr)
	if err != nil{
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	notif, err := h.userService.LikeUser(r.Context(), userID, likedID)
	if err != nil{
		helper.HandleError(w, err)
		return
	}
	res := LikeUserResponse{
		Connection: notif,
	}
	if notif != nil{
		res.Message = "It's a match!"
	} else{
		res.Message = "User liked successfully"
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

// /users/{userID}/like DELETE
func (h *UserHandler) UnlikeUserHandler(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	likedIDStr := chi.URLParam(r, string(helper.UserIDUrlParam))
	likedID, err := uuid.Parse(likedIDStr)
	if err != nil{
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.userService.UnlikeUser(r.Context(), userID, likedID); err != nil{
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User unliked successfully"})
}


type GetMyLikedListResponse struct{
	Likes []*entity.Like `json:"likes"`
}
// /users/me/likes GET
func (h *UserHandler) GetMyLikedListHandler(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	likes, err := h.userService.FindMyLikedList(r.Context(), userID)
	if err != nil{
		helper.HandleError(w, err)
		return
	}
	res := GetMyLikedListResponse{
		Likes: likes,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

type GetMyViewedListResponse struct{
	Views []*entity.View `json:"views"`
}
// /users/me/views GET
func(h *UserHandler) GetMyViewedListHandler(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	views, err := h.userService.FindMyViewedList(r.Context(), userID)
	if err != nil{
		helper.HandleError(w, err)
		return
	}
	res := GetMyViewedListResponse{
		Views: views,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

// users/me DELETE
func(h *UserHandler) DeleteMyAccountHandler(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if err := h.userService.DeleteUser(r.Context(), userID); err != nil{
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User account deleted successfully"})
}


type GetMyBlockedListResponse struct{
	Blocks []*entity.Block `json:"blocks"`
}
// /users/me/blocks GET
func (h *UserHandler) GetMyBlockedListHandler(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	blocks, err := h.userService.FindBlockList(r.Context(), userID)
	if err != nil{
		helper.HandleError(w, err)
		return
	}
	res := GetMyBlockedListResponse{
		Blocks: blocks,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

// /users/{userID}/block POST
func (h *UserHandler) BlockUserHandler(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	blockedIDStr := chi.URLParam(r, string(helper.UserIDUrlParam))
	blockedID, err := uuid.Parse(blockedIDStr)
	if err != nil{
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.userService.BlockUser(r.Context(), userID, blockedID); err != nil{
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User blocked successfully"})
}

// /users/{userID}/block DELETE
func (h *UserHandler) UnblockUserHandler(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	blockedIDStr := chi.URLParam(r, string(helper.UserIDUrlParam))
	blockedID, err := uuid.Parse(blockedIDStr)
	if err != nil{
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.userService.UnblockUser(r.Context(), userID, blockedID); err != nil{
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User unblocked successfully"})	
}
