package uow

import (
	"context"
	"github.com/icchon/matcha/api/internal/infrastructure/postgres"
	"github.com/jmoiron/sqlx"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(m RepositoryManager) error) error
}

type unitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) Do(ctx context.Context, fn func(m RepositoryManager) error) error {
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
	)
	if err = fn(manager); err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}
