package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
)

type profileService struct {
	uow         repo.UnitOfWork
	profileRepo repo.UserProfileRepository
	pictureRepo repo.PictureQueryRepository
	viewRepo    repo.ViewQueryRepository
	likeRepo    repo.LikeQueryRepository
	notifSvc    service.NotificationService
	fileClient  client.FileClient
}

var _ service.ProfileService = (*profileService)(nil)

func NewProfileService(uow repo.UnitOfWork, profileRepo repo.UserProfileRepository, fileClient client.FileClient, pictureRepo repo.PictureQueryRepository, viewRepo repo.ViewQueryRepository, likeRepo repo.LikeQueryRepository) *profileService {
	return &profileService{uow: uow, profileRepo: profileRepo, fileClient: fileClient, pictureRepo: pictureRepo, viewRepo: viewRepo, likeRepo: likeRepo}
}

func (s *profileService) CreateProfile(ctx context.Context, profile *entity.UserProfile) (*entity.UserProfile, error) {
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		return rm.ProfileRepo().Create(ctx, profile)
	}); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID uuid.UUID, profile *entity.UserProfile) (*entity.UserProfile, error) {
	target, err := s.profileRepo.Find(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	if target == nil {
		return nil, apperrors.ErrNotFound
	}

	{
		if profile.FirstName.Valid {
			target.FirstName = profile.FirstName
		}
		if profile.LastName.Valid {
			target.LastName = profile.LastName
		}
		if profile.Username.Valid {
			target.Username = profile.Username
		}
		if profile.Gender.Valid {
			target.Gender = profile.Gender
		}
		if profile.SexualPreference.Valid {
			target.SexualPreference = profile.SexualPreference
		}
		if profile.Biography.Valid {
			target.Biography = profile.Biography
		}
		if profile.LocationName.Valid {
			target.LocationName = profile.LocationName
		}
	}

	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		return rm.ProfileRepo().Update(ctx, target)
	}); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	return target, nil
}

func (s *profileService) FindProfile(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	profile, err := s.profileRepo.Find(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, apperrors.ErrNotFound
	}
	return profile, nil
}

func (s *profileService) VeiwProfile(ctx context.Context, viewerID, viewedID uuid.UUID) error {
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		view := &entity.View{
			ViewerID: viewerID,
			ViewedID: viewedID,
		}
		return rm.ViewRepo().Create(ctx, view)
	}); err != nil {
		return err
	}
	if _, err := s.notifSvc.CreateAndSendNotofication(ctx, viewerID, viewedID, entity.NotifView); err != nil {
		return err
	}
	return nil
}

func (s *profileService) FindWhoViewedMeList(ctx context.Context, userID uuid.UUID) ([]*entity.View, error) {
	views, err := s.viewRepo.Query(ctx, &repo.ViewQuery{ViewedID: &userID})
	if err != nil {
		return nil, err
	}
	return views, nil
}

func (s *profileService) FindWhoLikedMeList(ctx context.Context, userID uuid.UUID) ([]*entity.Like, error) {
	likes, err := s.likeRepo.Query(ctx, &repo.LikeQuery{LikedID: &userID})
	if err != nil {
		return nil, err
	}
	return likes, nil
}
