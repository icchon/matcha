package repo

import (
	"context"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(m RepositoryManager) error) error
}

type RepositoryManager interface {
	UserRepo() UserRepository
	AuthRepo() AuthRepository
	ConnectionRepo() ConnectionRepository
	MessageRepo() MessageRepository
	NotificationRepo() NotificationRepository
	PasswordResetRepo() PasswordResetRepository
	PictureRepo() PictureRepository
	RefreshTokenRepo() RefreshTokenRepository
	UserTagRepo() UserTagRepository
	VerificationTokenRepo() VerificationTokenRepository
	ProfileRepo() UserProfileRepository
	ViewRepo() ViewRepository
	LikeRepo() LikeRepository
	BlockRepo() BlockRepository
	UserDataRepo() UserDataRepository
	ReportRepo() ReportRepository
}
