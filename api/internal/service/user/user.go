package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
)

type userService struct {
	uow            repo.UnitOfWork
	likeRepo       repo.LikeQueryRepository
	viewRepo       repo.ViewQueryRepository
	connectionRepo repo.ConnectionQueryRepository
	notifSvc       service.NotificationService
}

var _ service.UserService = (*userService)(nil)

func NewUserService(uow repo.UnitOfWork, likeRepo repo.LikeQueryRepository, viewRepo repo.ViewQueryRepository, connectionRepo repo.ConnectionQueryRepository, notifSvc service.NotificationService) service.UserService {
	return &userService{
		uow:            uow,
		likeRepo:       likeRepo,
		viewRepo:       viewRepo,
		connectionRepo: connectionRepo,
		notifSvc:       notifSvc,
	}
}

func (u *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return u.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		if err := rm.UserRepo().Delete(ctx, userID); err != nil {
			return err
		}
		return nil
	})
}

// conection, like, view delete
func (s *userService) BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	return s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		if err := rm.ConnectionRepo().Delete(ctx, blockerID, blockedID); err != nil {
			return err
		}
		if err := rm.LikeRepo().Delete(ctx, blockerID, blockedID); err != nil {
			return err
		}
		if err := rm.LikeRepo().Delete(ctx, blockedID, blockerID); err != nil {
			return err
		}
		if err := rm.ViewRepo().Delete(ctx, blockerID, blockedID); err != nil {
			return err
		}
		if err := rm.ViewRepo().Delete(ctx, blockedID, blockerID); err != nil {
			return err
		}
		block := &entity.Block{
			BlockerID: blockerID,
			BlockedID: blockedID,
		}
		if err := rm.BlockRepo().Create(ctx, block); err != nil {
			return err
		}
		return nil
	})
}

func (s *userService) UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	return s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		if err := rm.BlockRepo().Delete(ctx, blockerID, blockedID); err != nil {
			return err
		}
		return nil
	})
}

func (s *userService) FindBlockList(ctx context.Context, userID uuid.UUID) ([]*entity.Block, error) {
	var blocks []*entity.Block
	var err error
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		blocks, err = rm.BlockRepo().Query(ctx, &repo.BlockQuery{BlockerID: &userID})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return blocks, nil
}

func (s *userService) LikeUser(ctx context.Context, likerID, likedID uuid.UUID) (*entity.Connection, error) {
	like, err := s.likeRepo.Find(ctx, likedID, likerID)
	if err != nil {
		return nil, err
	}
	love := (like != nil)
	var connection *entity.Connection
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		like := &entity.Like{
			LikerID: likerID,
			LikedID: likedID,
		}
		if err := rm.LikeRepo().Create(ctx, like); err != nil {
			return err
		}
		if love {
			user1ID, user2ID := likerID, likedID
			if likerID.String() > likedID.String() {
				user1ID, user2ID = likedID, likerID
			}
			connection = &entity.Connection{
				User1ID: user1ID,
				User2ID: user2ID,
			}
			if err := rm.ConnectionRepo().Create(ctx, connection); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if _, err := s.notifSvc.CreateAndSendNotofication(ctx, likerID, likedID, entity.NotifLike); err != nil {
		return nil, err
	}
	if love {
		if _, err := s.notifSvc.CreateAndSendNotofication(ctx, likerID, likedID, entity.NotifMatch); err != nil {
			return nil, err
		}
		if _, err := s.notifSvc.CreateAndSendNotofication(ctx, likedID, likerID, entity.NotifMatch); err != nil {
			return nil, err
		}
	}
	return connection, nil
}

func (s *userService) UnlikeUser(ctx context.Context, likerID, likedID uuid.UUID) error {
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		if err := rm.LikeRepo().Delete(ctx, likerID, likedID); err != nil {
			return err
		}
		if err := rm.ConnectionRepo().Delete(ctx, likerID, likedID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	if _, err := s.notifSvc.CreateAndSendNotofication(ctx, likerID, likedID, entity.NotifUnlike); err != nil {
		return err
	}
	return nil
}

func (s *userService) FindMyLikedList(ctx context.Context, userID uuid.UUID) ([]*entity.Like, error) {
	likes, err := s.likeRepo.Query(ctx, &repo.LikeQuery{LikerID: &userID})
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (s *userService) FindMyViewedList(ctx context.Context, userID uuid.UUID) ([]*entity.View, error) {
	views, err := s.viewRepo.Query(ctx, &repo.ViewQuery{ViewerID: &userID})
	if err != nil {
		return nil, err
	}
	return views, nil
}

func (s *userService) FindConnections(ctx context.Context, userID uuid.UUID) ([]*entity.Connection, error) {
	connection1, err := s.connectionRepo.Query(ctx, &repo.ConnectionQuery{User1ID: &userID})
	if err != nil {
		return nil, err
	}
	connection2, err := s.connectionRepo.Query(ctx, &repo.ConnectionQuery{User2ID: &userID})
	if err != nil {
		return nil, err
	}
	connections := make([]*entity.Connection, 0, len(connection1)+len(connection2))
	connections = append(connections, connection1...)
	connections = append(connections, connection2...)
	return connections, nil
}
