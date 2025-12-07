package handler

import (
	"encoding/json"
	"io"
	"log"
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

type ProfileHandler struct {
	profileSvc service.ProfileService
}

func NewProfileHandler(profileSvc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileSvc: profileSvc}
}

// ProfileResponse is the DTO for sending profile data in API responses.

type ProfileResponse struct {

	UserID           uuid.UUID  `json:"user_id"`

	FirstName        *string    `json:"first_name,omitempty"`

	LastName         *string    `json:"last_name,omitempty"`

	Username         *string    `json:"username,omitempty"`

	Gender           *string    `json:"gender,omitempty"`

	SexualPreference *string    `json:"sexual_preference,omitempty"`

	Birthday         *time.Time `json:"birthday,omitempty"`

	Occupation       *string    `json:"occupation,omitempty"`

	Biography        *string    `json:"biography,omitempty"`

	FameRating       *int32     `json:"fame_rating,omitempty"`

	LocationName     *string    `json:"location_name,omitempty"`

	Distance         *float64   `json:"distance,omitempty"`

}



type PictureResponse struct {

	ID           int32     `json:"id"`

	UserID       uuid.UUID `json:"user_id"`

	URL          string    `json:"url"`

	IsProfilePic *bool     `json:"is_profile_pic,omitempty"`

	CreatedAt    time.Time `json:"created_at"`

}



func newPictureResponse(p *entity.Picture) *PictureResponse {

	if p == nil {

		return nil

	}

	res := &PictureResponse{

		ID:        p.ID,

		UserID:    p.UserID,

		URL:       p.URL,

		CreatedAt: p.CreatedAt,

	}

	if p.IsProfilePic.Valid {

		res.IsProfilePic = &p.IsProfilePic.Bool

	}

	return res

}



// newProfileResponse converts a domain entity.UserProfile to a ProfileResponse DTO.

func newProfileResponse(p *entity.UserProfile) *ProfileResponse {

	if p == nil {

		return nil

	}

	res := &ProfileResponse{

		UserID: p.UserID,

	}

	if p.FirstName.Valid {

		res.FirstName = &p.FirstName.String

	}

	if p.LastName.Valid {

		res.LastName = &p.LastName.String

	}

	if p.Username.Valid {

		res.Username = &p.Username.String

	}

	if p.Gender.Valid {

		res.Gender = &p.Gender.String

	}

	if p.SexualPreference.Valid {

		res.SexualPreference = &p.SexualPreference.String

	}

	if p.Birthday.Valid {

		res.Birthday = &p.Birthday.Time

	}

	if p.Occupation.Valid {

		res.Occupation = &p.Occupation.String

	}

	if p.Biography.Valid {

		res.Biography = &p.Biography.String

	}

	if p.FameRating.Valid {

		res.FameRating = &p.FameRating.Int32

	}

	if p.LocationName.Valid {

		res.LocationName = &p.LocationName.String

	}

	if p.Distance.Valid {

		res.Distance = &p.Distance.Float64

	}

	return res

}



// CreateProfileRequest is the DTO for creating a profile.

type CreateProfileRequest struct {

	FirstName        *string    `json:"first_name"`

	LastName         *string    `json:"last_name"`

	Username         *string    `json:"username"`

	Gender           *string    `json:"gender"`

	SexualPreference *string    `json:"sexual_preference"`

	Birthday         *time.Time `json:"birthday"`

	Occupation       *string    `json:"occupation"`

	Biography        *string    `json:"biography"`

	LocationName     *string    `json:"location_name"`

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

		FirstName:        helper.NewNullString(req.FirstName),

		LastName:         helper.NewNullString(req.LastName),

		Username:         helper.NewNullString(req.Username),

		Gender:           helper.NewNullString(req.Gender),

		SexualPreference: helper.NewNullString(req.SexualPreference),

		Birthday:         helper.NewNullTime(req.Birthday),

		Occupation:       helper.NewNullString(req.Occupation),

		Biography:        helper.NewNullString(req.Biography),

		LocationName:     helper.NewNullString(req.LocationName),

	}

	createdProfile, err := h.profileSvc.CreateProfile(r.Context(), profile)

	if err != nil {

		helper.HandleError(w, err)

		return

	}

	helper.RespondWithJSON(w, http.StatusCreated, newProfileResponse(createdProfile))

}



// UpdateProfileRequest is the DTO for updating a profile.

type UpdateProfileRequest struct {

	FirstName        *string    `json:"first_name"`

	LastName         *string    `json:"last_name"`

	Username         *string    `json:"username"`

	Gender           *string    `json:"gender"`

	SexualPreference *string    `json:"sexual_preference"`

	Birthday         *time.Time `json:"birthday"`

	Occupation       *string    `json:"occupation"`

	Biography        *string    `json:"biography"`

	LocationName     *string    `json:"location_name"`

}



// /profile GET

func (h *ProfileHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)

	if !ok {

		helper.HandleError(w, apperrors.ErrInternalServer)

		return

	}

	profile, err := h.profileSvc.FindProfile(r.Context(), userID)

	if err != nil {

		helper.HandleError(w, err)

		return

	}

	helper.RespondWithJSON(w, http.StatusOK, newProfileResponse(profile))

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



	// This conversion logic should ideally be in the service layer,

	// but for now, we'll do it here to update the entity.

	profileToUpdate := &entity.UserProfile{

		FirstName:        helper.NewNullString(req.FirstName),

		LastName:         helper.NewNullString(req.LastName),

		Username:         helper.NewNullString(req.Username),

		Gender:           helper.NewNullString(req.Gender),

		SexualPreference: helper.NewNullString(req.SexualPreference),

		Birthday:         helper.NewNullTime(req.Birthday),

		Occupation:       helper.NewNullString(req.Occupation),

		Biography:        helper.NewNullString(req.Biography),

		LocationName:     helper.NewNullString(req.LocationName),

	}



	updatedProfile, err := h.profileSvc.UpdateProfile(r.Context(), userID, profileToUpdate)

	if err != nil {

		helper.HandleError(w, err)

		return

	}

	helper.RespondWithJSON(w, http.StatusOK, newProfileResponse(updatedProfile))

}



// /me/profile/pictures GET

func (h *ProfileHandler) GetMyPicturesHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)

	if !ok {

		helper.HandleError(w, apperrors.ErrInternalServer)

		return

	}

	pictures, err := h.profileSvc.FindPictures(r.Context(), userID)

	if err != nil {

		helper.HandleError(w, err)

		return

	}

	res := make([]*PictureResponse, len(pictures))

	for i, p := range pictures {

		res[i] = newPictureResponse(p)

	}

	helper.RespondWithJSON(w, http.StatusOK, res)

}



// /profile/pictures POST

func (h *ProfileHandler) UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB

		helper.HandleError(w, apperrors.ErrInvalidInput)

		return

	}



	file, _, err := r.FormFile("image")

	if err != nil {

		helper.HandleError(w, apperrors.ErrInvalidInput)

		return

	}

	defer file.Close()



	img, err := io.ReadAll(file)

	if err != nil {

		helper.HandleError(w, apperrors.ErrInvalidInput)

		return

	}



	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)

	if !ok {

		helper.HandleError(w, apperrors.ErrInternalServer)

		return

	}



	picture, err := h.profileSvc.UploadPicture(r.Context(), userID, img)

	if err != nil {

		helper.HandleError(w, err)

		return

	}



	helper.RespondWithJSON(w, http.StatusCreated, newPictureResponse(picture))

}



// /profile/pictures/{pictureID} DELETE

func (h *ProfileHandler) DeleteProfilePictureHandler(w http.ResponseWriter, r *http.Request) {

	pictureIDStr := chi.URLParam(r, "pictureID")

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

		w.WriteHeader(http.StatusNoContent)

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



	res := make([]*ProfileResponse, len(profiles))

	for i, p := range profiles {

		res[i] = newProfileResponse(p)

	}



	helper.RespondWithJSON(w, http.StatusOK, res)

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



		helper.RespondWithJSON(w, http.StatusOK, newProfileResponse(profile))



	



	}



	



	// GetUserPicturesHandler handles the request to get a specific user's pictures.



	// GET /users/{userID}/pictures



	func (h *ProfileHandler) GetUserPicturesHandler(w http.ResponseWriter, r *http.Request) {



		userIDStr := chi.URLParam(r, "userID")



		userID, err := uuid.Parse(userIDStr)



		if err != nil {



			helper.HandleError(w, apperrors.ErrInvalidInput)



			return



		}



	



		pictures, err := h.profileSvc.FindPictures(r.Context(), userID)



		if err != nil {



			helper.HandleError(w, err)



			return



		}



		res := make([]*PictureResponse, len(pictures))



		for i, p := range pictures {



			res[i] = newPictureResponse(p)



		}



		helper.RespondWithJSON(w, http.StatusOK, res)



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

	res := make([]*ProfileResponse, len(profiles))

	for i, p := range profiles {

		res[i] = newProfileResponse(p)

	}

	helper.RespondWithJSON(w, http.StatusOK, res)

}



type UpdatePictureStatusRequest struct {

	IsProfilePic bool `json:"is_profile_pic"`

}



// /profile/pictures/{pictureID}/status PUT

func (h *ProfileHandler) UpdatePictureStatusHandler(w http.ResponseWriter, r *http.Request) {

	pictureIDStr := chi.URLParam(r, "pictureID")

	pictureID, err := strconv.Atoi(pictureIDStr)

	if err != nil {

		helper.HandleError(w, apperrors.ErrInvalidInput)

		return

	}



	var req UpdatePictureStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		helper.HandleError(w, apperrors.ErrInvalidInput)

		return

	}



	userID, ok := r.Context().Value(middleware.UserIDContextKey).(uuid.UUID)

	if !ok {

		helper.HandleError(w, apperrors.ErrInternalServer)

		return

	}



	if err := h.profileSvc.UpdatePictureStatus(r.Context(), userID, int32(pictureID), req.IsProfilePic); err != nil {

		helper.HandleError(w, err)

		return

	}



	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Picture status updated successfully"})

}


