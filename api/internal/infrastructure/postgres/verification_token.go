package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type verificationTokenRepository struct {
	db DBTX
}

func NewVerificationTokenRepository(db DBTX) repo.VerificationTokenRepository {
	return &verificationTokenRepository{db: db}
}

func (r *verificationTokenRepository) Create(ctx context.Context, params repo.CreateVerificationTokenParams) (*entity.VerificationToken, error) {
	query := `
		INSERT INTO verification_tokens (token, user_id, expires_at)
		VALUES (:token, :user_id, :expires_at)
		RETURNING *
	`
	var token entity.VerificationToken
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	err = stmt.GetContext(ctx, &token, params)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *verificationTokenRepository) Update(ctx context.Context, token *entity.VerificationToken) error {
	query := `
		UPDATE verification_tokens SET
			user_id = :user_id,
			expires_at = :expires_at
		WHERE token = :token
	`
	_, err := r.db.NamedExecContext(ctx, query, token)
	return err
}

func (r *verificationTokenRepository) Delete(ctx context.Context, token string) error {
	query := "DELETE FROM verification_tokens WHERE token = $1"
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *verificationTokenRepository) Find(ctx context.Context, token string) (*entity.VerificationToken, error) {
	var verificationToken entity.VerificationToken
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

func (r *verificationTokenRepository) Query(ctx context.Context, q *repo.VerificationTokenQuery) ([]*entity.VerificationToken, error) {
	query := "SELECT * FROM verification_tokens WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.Token != nil {
		query += fmt.Sprintf(" AND token = $%d", argCount)
		args = append(args, *q.Token)
		argCount++
	}
	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *q.UserID)
		argCount++
	}
	if q.ExpiresAt != nil {
		query += fmt.Sprintf(" AND expires_at = $%d", argCount)
		args = append(args, *q.ExpiresAt)
		argCount++
	}

	var tokens []*entity.VerificationToken
	if err := r.db.SelectContext(ctx, &tokens, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tokens, nil
}
