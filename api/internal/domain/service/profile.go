package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type ProfileService interface {
	CreateProfile(ctx context.Context, profile *entity.UserProfile) (*entity.UserProfile, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, profile *entity.UserProfile) (*entity.UserProfile, error)
	ViewProfile(ctx context.Context, viewerID, viewedID uuid.UUID) error
	FindWhoViewedMeList(ctx context.Context, userID uuid.UUID) ([]*entity.View, error)
	DeletePicture(ctx context.Context, pictureID int32, userID uuid.UUID) error
	FindPicture(ctx context.Context, pictureID int32) (*entity.Picture, error)
	FindPictures(ctx context.Context, userID uuid.UUID) ([]*entity.Picture, error)
	UpdatePictureStatus(ctx context.Context, userID uuid.UUID, pictureID int32, isProfilePic bool) error
	UploadPicture(ctx context.Context, userID uuid.UUID, image []byte) (*entity.Picture, error)
	UploadPicutures(ctx context.Context, userID uuid.UUID, images [][]byte) ([]*entity.Picture, error)
	FindWhoLikedMeList(ctx context.Context, userID uuid.UUID) ([]*entity.Like, error)
	FindProfile(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	ListProfiles(ctx context.Context, selfUserID uuid.UUID, lat, lon, dist *float64, ageMin, ageMax *int, gender *entity.Gender) ([]*entity.UserProfile, error)
	RecommendProfiles(ctx context.Context, selfUserID uuid.UUID) ([]*entity.UserProfile, error)
}
