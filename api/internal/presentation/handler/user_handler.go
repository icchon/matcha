package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/service"
	"github.com/icchon/matcha/api/internal/presentation/helper"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
)

type UserHandler struct {
	userService    service.UserService
	profileService service.ProfileService
}

func NewUserHandler(userService service.UserService, profileService service.ProfileService) *UserHandler {
	return &UserHandler{userService: userService, profileService: profileService}
}

type LikeUserResponse struct {
	Connection *entity.Connection `json:"connection"`
	Message    string             `json:"message"`
}

// /users/{userID}/like POST
func (h *UserHandler) LikeUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	likedIDStr := chi.URLParam(r, string(helper.UserIDUrlParam))
	likedID, err := uuid.Parse(likedIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	notif, err := h.userService.LikeUser(r.Context(), userID, likedID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := LikeUserResponse{
		Connection: notif,
	}
	if notif != nil {
		res.Message = "It's a match!"
	} else {
		res.Message = "User liked successfully"
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

// /users/{userID}/like DELETE
func (h *UserHandler) UnlikeUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	likedIDStr := chi.URLParam(r, string(helper.UserIDUrlParam))
	likedID, err := uuid.Parse(likedIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.userService.UnlikeUser(r.Context(), userID, likedID); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User unliked successfully"})
}

type GetMyLikedListResponse struct {
	Likes []*entity.Like `json:"likes"`
}

// /users/me/likes GET
func (h *UserHandler) GetMyLikedListHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	likes, err := h.userService.FindMyLikedList(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := GetMyLikedListResponse{
		Likes: likes,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

type GetMyViewedListResponse struct {
	Views []*entity.View `json:"views"`
}

// /users/me/views GET
func (h *UserHandler) GetMyViewedListHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	views, err := h.userService.FindMyViewedList(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := GetMyViewedListResponse{
		Views: views,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

// users/me DELETE
func (h *UserHandler) DeleteMyAccountHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if err := h.userService.DeleteUser(r.Context(), userID); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User account deleted successfully"})
}

type GetMyBlockedListResponse struct {
	Blocks []*entity.Block `json:"blocks"`
}

// /users/me/blocks GET
func (h *UserHandler) GetMyBlockedListHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	blocks, err := h.userService.FindBlockList(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := GetMyBlockedListResponse{
		Blocks: blocks,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

// /users/{userID}/block POST
func (h *UserHandler) BlockUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	blockedIDStr := chi.URLParam(r, string(helper.UserIDUrlParam))
	blockedID, err := uuid.Parse(blockedIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.userService.BlockUser(r.Context(), userID, blockedID); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User blocked successfully"})
}

// /users/{userID}/block DELETE
func (h *UserHandler) UnblockUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	blockedIDStr := chi.URLParam(r, string(helper.UserIDUrlParam))
	blockedID, err := uuid.Parse(blockedIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.userService.UnblockUser(r.Context(), userID, blockedID); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User unblocked successfully"})
}

// GetMyUserDataHandler handles GET /users/me/data
func (h *UserHandler) GetMyUserDataHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	userData, err := h.userService.GetUserData(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, userData)
}

// CreateMyUserDataHandler handles POST /users/me/data
func (h *UserHandler) CreateMyUserDataHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	var userData entity.UserData
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	userData.UserID = userID

	if err := h.userService.CreateUserData(r.Context(), &userData); err != nil {
		helper.HandleError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, userData)
}

type UpdateMyUserDataRequest struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// UpdateMyUserDataHandler handles PUT /users/me/data
func (h *UserHandler) UpdateMyUserDataHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	var req UpdateMyUserDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}

	userData := entity.UserData{
		UserID: userID,
	}
	if req.Latitude != nil {
		userData.Latitude = sql.NullFloat64{Float64: *req.Latitude, Valid: true}
	}
	if req.Longitude != nil {
		userData.Longitude = sql.NullFloat64{Float64: *req.Longitude, Valid: true}
	}

	if err := h.userService.UpdateUserData(r.Context(), &userData); err != nil {
		helper.HandleError(w, err)
		return
	}

	res := map[string]interface{}{
		"user_id": userData.UserID,
	}
	if userData.Latitude.Valid {
		res["latitude"] = userData.Latitude.Float64
	} else {
		res["latitude"] = nil
	}
	if userData.Longitude.Valid {
		res["longitude"] = userData.Longitude.Float64
	} else {
		res["longitude"] = nil
	}
	if userData.InternalScore.Valid {
		res["internal_score"] = userData.InternalScore.Int32
	} else {
		res["internal_score"] = nil
	}

	helper.RespondWithJSON(w, http.StatusOK, res)
}

// GetAllTagsHandler handles GET /tags
func (h *UserHandler) GetAllTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := h.userService.GetAllTags(r.Context())
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, tags)
}

// GetUserTagsHandler handles GET /users/me/tags
func (h *UserHandler) GetUserTagsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}
	tags, err := h.userService.GetUserTags(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, tags)
}

type AddUserTagRequest struct {
	TagID int32 `json:"tag_id"`
}

// AddUserTagHandler handles POST /users/me/tags
func (h *UserHandler) AddUserTagHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}
	var req AddUserTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.userService.AddUserTag(r.Context(), userID, req.TagID); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Tag added successfully"})
}

// DeleteUserTagHandler handles DELETE /users/me/tags/{tagID}
func (h *UserHandler) DeleteUserTagHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}
	tagIDStr := chi.URLParam(r, "tagID")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	if err := h.userService.DeleteUserTag(r.Context(), userID, int32(tagID)); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusNoContent, nil)
}
