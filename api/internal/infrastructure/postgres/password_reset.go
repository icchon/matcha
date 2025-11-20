package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/icchon/matcha/api/internal/domain"
)

type passwordResetRepository struct {
	db DBTX
}

func NewPasswordResetRepository(db DBTX) domain.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Save(ctx context.Context, passwordReset *domain.PasswordReset) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at)
		VALUES (:user_id, :token, :expires_at)
		ON CONFLICT (user_id) DO UPDATE SET
			token = :token,
			expires_at = :expires_at
	`
	_, err := r.db.NamedExecContext(ctx, query, passwordReset)
	return err
}

func (r *passwordResetRepository) Delete(ctx context.Context, token string) error {
	query := "DELETE FROM password_resets WHERE token = $1"
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *passwordResetRepository) Find(ctx context.Context, token string) (*domain.PasswordReset, error) {
	var passwordReset domain.PasswordReset
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
