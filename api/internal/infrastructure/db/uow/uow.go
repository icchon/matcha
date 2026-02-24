package uow

import (
	"context"
	"errors"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/infrastructure/db/postgres"
	"github.com/jmoiron/sqlx"
)

type unitOfWork struct {
	db *sqlx.DB
}

var _ repo.UnitOfWork = (*unitOfWork)(nil)

func NewUnitOfWork(db *sqlx.DB) *unitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) Do(ctx context.Context, fn func(m repo.RepositoryManager) error) error {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	manager := NewRepositoryManager(
		postgres.NewUserRepository(tx),
		postgres.NewAuthRepository(tx),
		postgres.NewConnectionRepository(tx),
		postgres.NewMessageRepository(tx),
		postgres.NewNotificationRepository(tx),
		postgres.NewPasswordResetRepository(tx),
		postgres.NewPictureRepository(tx),
		postgres.NewRefreshTokenRepository(tx),
		postgres.NewUserTagRepository(tx),
		postgres.NewVerificationTokenRepository(tx),
		postgres.NewUserProfileRepository(tx),
		postgres.NewViewRepository(tx),
		postgres.NewLikeRepository(tx),
		postgres.NewBlockRepository(tx),
		postgres.NewUserDataRepository(tx),
		postgres.NewReportRepository(tx),
	)
	if err = fn(manager); err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.Join(err, txErr)
		}
		return err
	}
	return tx.Commit()
}
