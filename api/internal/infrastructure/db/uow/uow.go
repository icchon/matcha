package uow

import (
	"context"
	"log"

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
	)
	if err = fn(manager); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("CRITICAL: rollback failed: %v (original: %v)", rbErr, err)
		}
		return err
	}
	return tx.Commit()
}
