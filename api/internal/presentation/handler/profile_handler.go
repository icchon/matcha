package handler

import (
	"log" // Add this import
	"net/http"

	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/service"
	"github.com/icchon/matcha/api/internal/presentation/helper"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
	"strconv"
)

type ProfileHandler struct {
	profileSvc service.ProfileService
}

func NewProfileHandler(profileSvc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileSvc: profileSvc}
}

type CreateProfileRequest struct {
	FirstName        sql.NullString `db:"first_name"`        // NULLを許容
	LastName         sql.NullString `db:"last_name"`         // NULLを許容
	Username         sql.NullString `db:"username"`          // NULLを許容
	Gender           sql.NullString `db:"gender"`            // ENUM型
	SexualPreference sql.NullString `db:"sexual_preference"` // ENUM型
	Birthday         sql.NullTime   `db:"birthday"`
	Occupation       sql.NullString `db:"occupation"`
	Biography        sql.NullString `db:"biography"`
	LocationName     sql.NullString `db:"location_name"`
}
type CreateProfileResponse struct {
	UserID           uuid.UUID      `json:"user_id"`
	FirstName        sql.NullString `json:"first_name"`
	LastName         sql.NullString `json:"last_name"`
	Username         sql.NullString `json:"username"`
	Gender           sql.NullString `json:"gender"`
	SexualPreference sql.NullString `json:"sexual_preference"`
	Birthday         sql.NullTime   `json:"birthday"`
	Occupation       sql.NullString `json:"occupation"`
	Biography        sql.NullString `json:"biography"`
	LocationName     sql.NullString `json:"location_name"`
}

// /profile POST
func (h *ProfileHandler) CreateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	profile := &entity.UserProfile{
		UserID:           userID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Username:         req.Username,
		Gender:           req.Gender,
		SexualPreference: req.SexualPreference,
		Birthday:         req.Birthday,
		Occupation:       req.Occupation,
		Biography:        req.Biography,
		LocationName:     req.LocationName,
	}
	profile, err := h.profileSvc.CreateProfile(r.Context(), profile)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := CreateProfileResponse{
		UserID:           profile.UserID,
		FirstName:        profile.FirstName,
		LastName:         profile.LastName,
		Username:         profile.Username,
		Gender:           profile.Gender,
		SexualPreference: profile.SexualPreference,
		Birthday:         profile.Birthday,
		Occupation:       profile.Occupation,
		Biography:        profile.Biography,
		LocationName:     profile.LocationName,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

type UpdateProfileRequest struct {
	FirstName        sql.NullString `db:"first_name"`        // NULLを許容
	LastName         sql.NullString `db:"last_name"`         // NULLを許容
	Username         sql.NullString `db:"username"`          // NULLを許容
	Gender           sql.NullString `db:"gender"`            // ENUM型
	SexualPreference sql.NullString `db:"sexual_preference"` // ENUM型
	Birthday         sql.NullTime   `db:"birthday"`
	Occupation       sql.NullString `db:"occupation"`
	Biography        sql.NullString `db:"biography"`
	LocationName     sql.NullString `db:"location_name"`
}
type UpdateProfileResponse struct {
	UserID           uuid.UUID      `json:"user_id"`
	FirstName        sql.NullString `json:"first_name"`
	LastName         sql.NullString `json:"last_name"`
	Username         sql.NullString `json:"username"`
	Gender           sql.NullString `json:"gender"`
	SexualPreference sql.NullString `json:"sexual_preference"`
	Birthday         sql.NullTime   `json:"birthday"`
	Occupation       sql.NullString `json:"occupation"`
	Biography        sql.NullString `json:"biography"`
	LocationName     sql.NullString `json:"location_name"`
}

// /profile PUT
func (h *ProfileHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	profile := &entity.UserProfile{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Username:         req.Username,
		Gender:           req.Gender,
		SexualPreference: req.SexualPreference,
		Birthday:         req.Birthday,
		Occupation:       req.Occupation,
		Biography:        req.Biography,
		LocationName:     req.LocationName,
	}
	profile, err := h.profileSvc.UpdateProfile(r.Context(), userID, profile)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := UpdateProfileResponse{
		UserID:           profile.UserID,
		FirstName:        profile.FirstName,
		LastName:         profile.LastName,
		Username:         profile.Username,
		Gender:           profile.Gender,
		SexualPreference: profile.SexualPreference,
		Birthday:         profile.Birthday,
		Occupation:       profile.Occupation,
		Biography:        profile.Biography,
		LocationName:     profile.LocationName,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

// /profile/pictures POST
func (h *ProfileHandler) UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
}

// /profile/pictures/{pictureID} DELETE
func (h *ProfileHandler) DeleteProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	pictureIDStr := chi.URLParam(r, string(helper.PictureIDParam))
	pictureID, err := strconv.Atoi(pictureIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	if err := h.profileSvc.DeletePicture(r.Context(), int32(pictureID), userID); err != nil {
		helper.HandleError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusNoContent, map[string]string{"message": "Picture deleted successfully"})
}

type GetWhoLikedMeListResponse struct {
	Likes []*entity.Like `json:"likes"`
}

// /profile/likes GET
func (h *ProfileHandler) GetWhoLikeMeListHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	likes, err := h.profileSvc.FindWhoLikedMeList(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := GetWhoLikedMeListResponse{
		Likes: likes,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

type GetWhoViewedMeListResponse struct {
	Views []*entity.View `json:"views"`
}

// /profile/views GET
func (h *ProfileHandler) GetWhoViewedMeListHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	views, err := h.profileSvc.FindWhoViewedMeList(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	res := struct {
		Views []*entity.View `json:"views"`
	}{
		Views: views,
	}
	helper.RespondWithJSON(w, http.StatusOK, res)
}

// ListProfilesHandler handles the request to list user profiles.
// GET /profiles
func (h *ProfileHandler) ListProfilesHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}

	var ageMin, ageMax *int
	if ageMinStr := r.URL.Query().Get("age_min"); ageMinStr != "" {
		if v, err := strconv.Atoi(ageMinStr); err == nil {
			ageMin = &v
		} else {
			helper.HandleError(w, apperrors.ErrInvalidInput)
			return
		}
	}
	if ageMaxStr := r.URL.Query().Get("age_max"); ageMaxStr != "" {
		if v, err := strconv.Atoi(ageMaxStr); err == nil {
			ageMax = &v
		} else {
			helper.HandleError(w, apperrors.ErrInvalidInput)
			return
		}
	}

	var gender *entity.Gender
	if genderStr := r.URL.Query().Get("gender"); genderStr != "" {
		g := entity.Gender(genderStr)
		gender = &g
	}

	// For now, we are not using location filters in this endpoint, so we pass nil.
	profiles, err := h.profileSvc.ListProfiles(r.Context(), userID, nil, nil, nil, ageMin, ageMax, gender)
	if err != nil {
		helper.HandleError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, profiles)
}

// GetUserProfileHandler handles the request to get a specific user's profile.
// GET /users/{userID}/profile
func (h *ProfileHandler) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		helper.HandleError(w, apperrors.ErrInvalidInput)
		return
	}

	// We can also add a view event here
	viewerID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)
	if !ok {
		helper.HandleError(w, apperrors.ErrInternalServer)
		return
	}
	if err := h.profileSvc.ViewProfile(r.Context(), viewerID, userID); err != nil {
		// Log this error but don't block the response
		log.Printf("failed to record view: %v", err)
	}

	profile, err := h.profileSvc.FindProfile(r.Context(), userID)
	if err != nil {
		helper.HandleError(w, err)
		return
	}
	if profile == nil {
		helper.HandleError(w, apperrors.ErrNotFound)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, profile)

}

// RecommendProfilesHandler handles the request to get recommended user profiles.

// GET /profiles/recommends

func (h *ProfileHandler) RecommendProfilesHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)

	if !ok {

		helper.HandleError(w, apperrors.ErrUnauthorized)

		return

	}

	profiles, err := h.profileSvc.RecommendProfiles(r.Context(), userID)

	if err != nil {

		helper.HandleError(w, err)

		return

	}

	helper.RespondWithJSON(w, http.StatusOK, profiles)

}
