package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type UserService interface {
	UnlikeUser(ctx context.Context, likerID, likedID uuid.UUID) error
	FindMyLikedList(ctx context.Context, userID uuid.UUID) ([]*entity.Like, error)
	FindMyViewedList(ctx context.Context, userID uuid.UUID) ([]*entity.View, error)
	FindConnections(ctx context.Context, userID uuid.UUID) ([]*entity.Connection, error)
	LikeUser(ctx context.Context, likerID, likedID uuid.UUID) (*entity.Connection, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	FindBlockList(ctx context.Context, userID uuid.UUID) ([]*entity.Block, error)
	GetUserData(ctx context.Context, userID uuid.UUID) (*entity.UserData, error)
	CreateUserData(ctx context.Context, userData *entity.UserData) error
	UpdateUserData(ctx context.Context, userData *entity.UserData) error
	GetAllTags(ctx context.Context) ([]*entity.Tag, error)
	GetUserTags(ctx context.Context, userID uuid.UUID) ([]*entity.Tag, error)
	AddUserTag(ctx context.Context, userID uuid.UUID, tagID int32) error
	DeleteUserTag(ctx context.Context, userID uuid.UUID, tagID int32) error
	ReportUser(ctx context.Context, reporterID, reportedID uuid.UUID, reason string) error
}
