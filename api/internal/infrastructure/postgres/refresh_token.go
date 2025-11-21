package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type refreshTokenRepository struct {
	db DBTX
}

func NewRefreshTokenRepository(db DBTX) repo.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Find(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	query := "SELECT * FROM refresh_tokens WHERE token_hash = $1"
	err := r.db.GetContext(ctx, &token, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepository) Query(ctx context.Context, q *repo.RefreshTokenQuery) ([]*entity.RefreshToken, error) {
	query := "SELECT * FROM refresh_tokens WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.TokenHash != nil {
		query += fmt.Sprintf(" AND token_hash = $%d", argCount)
		args = append(args, *q.TokenHash)
		argCount++
	}
	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *q.UserID)
		argCount++
	}
	if q.Revoked != nil {
		query += fmt.Sprintf(" AND revoked = $%d", argCount)
		args = append(args, *q.Revoked)
		argCount++
	}

	var tokens []*entity.RefreshToken
	if err := r.db.SelectContext(ctx, &tokens, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tokens, nil
}

func (r *refreshTokenRepository) Create(ctx context.Context, params repo.CreateRefreshTokenParams) (*entity.RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens (token_hash, user_id, expires_at)
		VALUES (:token_hash, :user_id, :expires_at)
		RETURNING *
	`
	var token entity.RefreshToken
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

func (r *refreshTokenRepository) Update(ctx context.Context, token *entity.RefreshToken) error {
	query := `
		UPDATE refresh_tokens SET
			expires_at = :expires_at,
			revoked = :revoked
		WHERE token_hash = :token_hash
	`
	_, err := r.db.NamedExecContext(ctx, query, token)
	return err
}

func (r *refreshTokenRepository) Delete(ctx context.Context, tokenHash string) error {
	query := "DELETE FROM refresh_tokens WHERE token_hash = $1"
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}
