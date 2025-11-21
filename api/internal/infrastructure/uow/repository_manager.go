package uow

import (
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type RepositoryManager interface {
	UserRepo() repo.UserRepository
	AuthRepo() repo.AuthRepository
	ConnectionRepo() repo.ConnectionRepository
	MessageRepo() repo.MessageRepository
	NotificationRepo() repo.NotificationRepository
	PasswordResetRepo() repo.PasswordResetRepository
	PictureRepo() repo.PictureRepository
	RefreshTokenRepo() repo.RefreshTokenRepository
	UserTagRepo() repo.UserTagRepository
	VerificationTokenRepo() repo.VerificationTokenRepository
}

type repositoryManager struct {
	userRepo              repo.UserRepository
	authRepo              repo.AuthRepository
	connectionRepo        repo.ConnectionRepository
	messageRepo           repo.MessageRepository
	notificationRepo      repo.NotificationRepository
	passwordResetRepo     repo.PasswordResetRepository
	pictureRepo           repo.PictureRepository
	refreshTokenRepo      repo.RefreshTokenRepository
	userTagRepo           repo.UserTagRepository
	verificationTokenRepo repo.VerificationTokenRepository
}

func NewRepositoryManager(
	userRepo repo.UserRepository,
	authRepo repo.AuthRepository,
	connectionRepo repo.ConnectionRepository,
	messageRepo repo.MessageRepository,
	notificationRepo repo.NotificationRepository,
	passwordResetRepo repo.PasswordResetRepository,
	pictureRepo repo.PictureRepository,
	refreshTokenRepo repo.RefreshTokenRepository,
	userTagRepo repo.UserTagRepository,
	verificationTokenRepo repo.VerificationTokenRepository,
) RepositoryManager {
	return &repositoryManager{
		userRepo:              userRepo,
		authRepo:              authRepo,
		connectionRepo:        connectionRepo,
		messageRepo:           messageRepo,
		notificationRepo:      notificationRepo,
		passwordResetRepo:     passwordResetRepo,
		pictureRepo:           pictureRepo,
		refreshTokenRepo:      refreshTokenRepo,
		userTagRepo:           userTagRepo,
		verificationTokenRepo: verificationTokenRepo,
	}
}

func (r *repositoryManager) UserRepo() repo.UserRepository {
	return r.userRepo
}

func (r *repositoryManager) AuthRepo() repo.AuthRepository {
	return r.authRepo
}

func (r *repositoryManager) ConnectionRepo() repo.ConnectionRepository {
	return r.connectionRepo
}

func (r *repositoryManager) MessageRepo() repo.MessageRepository {
	return r.messageRepo
}

func (r *repositoryManager) NotificationRepo() repo.NotificationRepository {
	return r.notificationRepo
}

func (r *repositoryManager) PasswordResetRepo() repo.PasswordResetRepository {
	return r.passwordResetRepo
}

func (r *repositoryManager) PictureRepo() repo.PictureRepository {
	return r.pictureRepo
}

func (r *repositoryManager) RefreshTokenRepo() repo.RefreshTokenRepository {
	return r.refreshTokenRepo
}

func (r *repositoryManager) UserTagRepo() repo.UserTagRepository {
	return r.userTagRepo
}

func (r *repositoryManager) VerificationTokenRepo() repo.VerificationTokenRepository {
	return r.verificationTokenRepo
}
