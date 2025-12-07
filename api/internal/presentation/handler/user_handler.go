package handler

import (
	"encoding/json"
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

type UserHandler struct {
	userService    service.UserService
	profileService service.ProfileService
}

func NewUserHandler(userService service.UserService, profileService service.ProfileService) *UserHandler {
	return &UserHandler{userService: userService, profileService: profileService}
}

// UserDataResponse is the DTO for UserData responses.
type UserDataResponse struct {
	UserID        uuid.UUID `json:"user_id"`
	Latitude      *float64  `json:"latitude,omitempty"`
	Longitude     *float64  `json:"longitude,omitempty"`
	InternalScore *int32    `json:"internal_score,omitempty"`
}

// newUserDataResponse converts an entity.UserData to a UserDataResponse DTO.
func newUserDataResponse(ud *entity.UserData) *UserDataResponse {
	if ud == nil {
		return nil
	}
	res := &UserDataResponse{
		UserID: ud.UserID,
	}
	if ud.Latitude.Valid {
		res.Latitude = &ud.Latitude.Float64
	}
	if ud.Longitude.Valid {
		res.Longitude = &ud.Longitude.Float64
	}
	if ud.InternalScore.Valid {
		res.InternalScore = &ud.InternalScore.Int32
	}
	return res
}

// UserDataRequest is the DTO for creating/updating UserData.
type UserDataRequest struct {
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	InternalScore *int32   `json:"internal_score"`
}

type LikeResponse struct {
	LikerID   uuid.UUID `json:"liker_id"`
	LikedID   uuid.UUID `json:"liked_id"`
	CreatedAt time.Time `json:"created_at"`
}

func newLikeResponse(l *entity.Like) *LikeResponse {
	if l == nil {
		return nil
	}
	return &LikeResponse{
		LikerID:   l.LikerID,
		LikedID:   l.LikedID,
		CreatedAt: l.CreatedAt,
	}
}

type ViewResponse struct {
	ViewerID uuid.UUID `json:"viewer_id"`
	ViewedID uuid.UUID `json:"viewed_id"`
	ViewTime time.Time `json:"view_time"`
}

func newViewResponse(v *entity.View) *ViewResponse {
	if v == nil {
		return nil
	}
	return &ViewResponse{
		ViewerID: v.ViewerID,
		ViewedID: v.ViewedID,
		ViewTime: v.ViewTime,
	}
}

type BlockResponse struct {
	BlockerID uuid.UUID `json:"blocker_id"`
	BlockedID uuid.UUID `json:"blocked_id"`
}

func newBlockResponse(b *entity.Block) *BlockResponse {
	if b == nil {
		return nil
	}
	return &BlockResponse{
		BlockerID: b.BlockerID,
		BlockedID: b.BlockedID,
	}
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

// /users/me/likes GET
func (h *UserHandler) GetMyLikedListHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	likes, err := h.userService.FindMyLikedList(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := make([]*LikeResponse, len(likes))
	for i, l := range likes {
		res[i] = newLikeResponse(l)
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

// /users/me/views GET
func (h *UserHandler) GetMyViewedListHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	views, err := h.userService.FindMyViewedList(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := make([]*ViewResponse, len(views))
	for i, v := range views {
		res[i] = newViewResponse(v)
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


// /users/me/blocks GET
func (h *UserHandler) GetMyBlockedListHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	blocks, err := h.userService.FindBlockList(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := make([]*BlockResponse, len(blocks))
	for i, b := range blocks {
		res[i] = newBlockResponse(b)
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

// GetMyUserDataHandler handles GET /me/data
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

	helper.RespondWithJSON(w, http.StatusOK, newUserDataResponse(userData))
}

// CreateMyUserDataHandler handles POST /me/data
func (h *UserHandler) CreateMyUserDataHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	var req UserDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}

	userData := &entity.UserData{
		UserID:        userID,
		Latitude:      helper.NewNullFloat64(req.Latitude),
		Longitude:     helper.NewNullFloat64(req.Longitude),
		InternalScore: helper.NewNullInt32(req.InternalScore),
	}

	if err := h.userService.CreateUserData(r.Context(), userData); err != nil {
		helper.HandleError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, newUserDataResponse(userData))
}

// UpdateMyUserDataHandler handles PUT /me/data
func (h *UserHandler) UpdateMyUserDataHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrUnauthorized)
		return
	}

	var req UserDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}

	userData := &entity.UserData{
		UserID:        userID,
		Latitude:      helper.NewNullFloat64(req.Latitude),
		Longitude:     helper.NewNullFloat64(req.Longitude),
		InternalScore: helper.NewNullInt32(req.InternalScore),
	}

	if err := h.userService.UpdateUserData(r.Context(), userData); err != nil {
		helper.HandleError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, newUserDataResponse(userData))
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
