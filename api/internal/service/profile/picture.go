package profile

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

func (s *profileService) DeletePicture(ctx context.Context, pictureID int32, userID uuid.UUID) error {
	picture, err := s.pictureRepo.Find(ctx, pictureID)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	if picture == nil || picture.UserID != userID {
		return apperrors.ErrNotFound
	}
	return s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		return rm.PictureRepo().Delete(ctx, pictureID)
	})
}

func (s *profileService) FindPicture(ctx context.Context, pictureID int32) (*entity.Picture, error) {
	return s.pictureRepo.Find(ctx, pictureID)
}

func (s *profileService) FindPictures(ctx context.Context, userID uuid.UUID) ([]*entity.Picture, error) {
	return s.pictureRepo.Query(ctx, &repo.PictureQuery{UserID: &userID})
}

func (s *profileService) UpdatePictureStatus(ctx context.Context, userID uuid.UUID, pictureID int32, isProfilePic bool) error {
	pictures, err := s.pictureRepo.Query(ctx, &repo.PictureQuery{UserID: &userID})
	if err != nil {
		return apperrors.ErrInternalServer
	}
	if pictures == nil {
		return apperrors.ErrNotFound
	}
	return s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		for _, pic := range pictures {
			if isProfilePic {
				if pic.ID == pictureID {
					if isProfilePic == pic.IsProfilePic.Bool {
						continue
					}
					pic.IsProfilePic = sql.NullBool{Bool: isProfilePic, Valid: true}
				} else {
					if !pic.IsProfilePic.Bool {
						continue
					}
					pic.IsProfilePic = sql.NullBool{Bool: false, Valid: true}
				}
			} else {
				if pic.ID == pictureID {
					if isProfilePic == pic.IsProfilePic.Bool {
						continue
					}
					pic.IsProfilePic = sql.NullBool{Bool: isProfilePic, Valid: true}
				} else {
					continue
				}
			}
			if err := rm.PictureRepo().Update(ctx, pic); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *profileService) UploadPicture(ctx context.Context, userID uuid.UUID, image []byte) (*entity.Picture, error) {
	if len(image) == 0 {
		return nil, apperrors.ErrInvalidInput
	}
	url, err := s.fileClient.SaveImage(image, uuid.NewString())
	if err != nil {
		return nil, err
	}
	picture := &entity.Picture{
		UserID:       userID,
		URL:          url,
		IsProfilePic: sql.NullBool{Bool: false, Valid: true},
	}
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		return rm.PictureRepo().Create(ctx, picture)
	}); err != nil {
		return nil, err
	}
	return picture, nil
}

func (s *profileService) UploadPicutures(ctx context.Context, userID uuid.UUID, images [][]byte) ([]*entity.Picture, error) {
	n := len(images)
	if n == 0 {
		return nil, apperrors.ErrInvalidInput
	}
	var urls []string
	for _, img := range images {
		url, err := s.fileClient.SaveImage(img, uuid.NewString())
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	var pictures []*entity.Picture
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		for _, url := range urls {
			pic := &entity.Picture{
				UserID:       userID,
				URL:          url,
				IsProfilePic: sql.NullBool{Bool: false, Valid: true},
			}
			if err := rm.PictureRepo().Create(ctx, pic); err != nil {
				return err
			}
			pictures = append(pictures, pic)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return pictures, nil
}
