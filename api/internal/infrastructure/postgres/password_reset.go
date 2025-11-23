package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type passwordResetRepository struct {
	db DBTX
}

func NewPasswordResetRepository(db DBTX) repo.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, passwordReset *entity.PasswordReset) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at)
		VALUES (:user_id, :token, :expires_at)
		RETURNING *
	`
	return r.db.QueryRowxContext(ctx, query, passwordReset).StructScan(passwordReset)
}

func (r *passwordResetRepository) Update(ctx context.Context, passwordReset *entity.PasswordReset) error {
	query := `
		UPDATE password_resets SET
			token = :token,
			expires_at = :expires_at
		WHERE user_id = :user_id
	`
	_, err := r.db.NamedExecContext(ctx, query, passwordReset)
	return err
}

func (r *passwordResetRepository) Delete(ctx context.Context, token string) error {
	query := "DELETE FROM password_resets WHERE token = $1"
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *passwordResetRepository) Find(ctx context.Context, token string) (*entity.PasswordReset, error) {
	var passwordReset entity.PasswordReset
	query := "SELECT * FROM password_resets WHERE token = $1"
	err := r.db.GetContext(ctx, &passwordReset, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &passwordReset, nil
}

func (r *passwordResetRepository) Query(ctx context.Context, q *repo.PasswordResetQuery) ([]*entity.PasswordReset, error) {
	query := "SELECT * FROM password_resets WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *q.UserID)
		argCount++
	}
	if q.Token != nil {
		query += fmt.Sprintf(" AND token = $%d", argCount)
		args = append(args, *q.Token)
		argCount++
	}
	if q.ExpiresAt != nil {
		query += fmt.Sprintf(" AND expires_at = $%d", argCount)
		args = append(args, *q.ExpiresAt)
		argCount++
	}

	var passwordResets []*entity.PasswordReset
	if err := r.db.SelectContext(ctx, &passwordResets, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return passwordResets, nil
}
