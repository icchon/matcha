package uow

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(tx *sqlx.Tx) error) error
}

type unitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) Do(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	if err = fn(tx); err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}
