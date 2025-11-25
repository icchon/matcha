package uow

import (
	"github.com/icchon/matcha/api/internal/domain/repo"
)

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
	profileRepo           repo.UserProfileRepository
	viewRepo              repo.ViewRepository
	likeRepo              repo.LikeRepository
	blockRepo             repo.BlockRepository
	userDataRepo          repo.UserDataRepository
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
	profileRepo repo.UserProfileRepository,
	viewRepo repo.ViewRepository,
	likeRepo repo.LikeRepository,
	blockRepo repo.BlockRepository,
	userDataRepo repo.UserDataRepository,
) repo.RepositoryManager {
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
		profileRepo:           profileRepo,
		viewRepo:              viewRepo,
		likeRepo:              likeRepo,
		blockRepo:             blockRepo,
		userDataRepo:          userDataRepo,
	}
}

func (r *repositoryManager) UserDataRepo() repo.UserDataRepository {
	return r.userDataRepo
}

func (r *repositoryManager) BlockRepo() repo.BlockRepository {
	return r.blockRepo
}

func (r *repositoryManager) LikeRepo() repo.LikeRepository {
	return r.likeRepo
}

func (r *repositoryManager) ViewRepo() repo.ViewRepository {
	return r.viewRepo
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

func (r *repositoryManager) ProfileRepo() repo.UserProfileRepository {
	return r.profileRepo
}
