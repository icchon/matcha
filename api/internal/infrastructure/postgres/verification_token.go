package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/icchon/matcha/api/internal/domain"
)

type verificationTokenRepository struct {
	db DBTX
}

func NewVerificationTokenRepository(db DBTX) domain.VerificationTokenRepository {
	return &verificationTokenRepository{db: db}
}

func (r *verificationTokenRepository) Save(ctx context.Context, token *domain.VerificationToken) error {
	query := `
		INSERT INTO verification_tokens (token, user_id, expires_at)
		VALUES (:token, :user_id, :expires_at)
		ON CONFLICT (token) DO UPDATE SET
			user_id = :user_id,
			expires_at = :expires_at
	`
	_, err := r.db.NamedExecContext(ctx, query, token)
	return err
}

func (r *verificationTokenRepository) Delete(ctx context.Context, token string) error {
	query := "DELETE FROM verification_tokens WHERE token = $1"
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *verificationTokenRepository) Find(ctx context.Context, token string) (*domain.VerificationToken, error) {
	var verificationToken domain.VerificationToken
	query := "SELECT * FROM verification_tokens WHERE token = $1"
	err := r.db.GetContext(ctx, &verificationToken, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &verificationToken, nil
}
