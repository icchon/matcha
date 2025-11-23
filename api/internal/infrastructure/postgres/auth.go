package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type authRepository struct {
	db DBTX
}

func NewAuthRepository(db DBTX) repo.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Create(ctx context.Context, auth *entity.Auth) error {
	query := `
		INSERT INTO auths (user_id, provider, provider_uid, email, is_verified, password_hash)
		VALUES (:user_id, :provider, :provider_uid, :email, :is_verified, :password_hash)
		RETURNING *
	`
	return r.db.QueryRowxContext(ctx, query, auth).StructScan(auth)
}

func (r *authRepository) Update(ctx context.Context, auth *entity.Auth) error {
	query := `
		UPDATE auths SET
			email = :email,
			is_verified = :is_verified,
			password_hash = :password_hash
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, auth)
	return err
}

func (r *authRepository) Delete(ctx context.Context, userID uuid.UUID, provider entity.AuthProvider) error {
	query := "DELETE FROM auths WHERE user_id = $1 AND provider = $2"
	_, err := r.db.ExecContext(ctx, query, userID, provider)
	return err
}

func (r *authRepository) Find(ctx context.Context, userID uuid.UUID, provider entity.AuthProvider) (*entity.Auth, error) {
	var auth entity.Auth
	query := "SELECT * FROM auths WHERE user_id = $1 AND provider = $2"
	err := r.db.GetContext(ctx, &auth, query, userID, provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) Query(ctx context.Context, q *repo.AuthQuery) ([]*entity.Auth, error) {
	query := "SELECT * FROM auths WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *q.UserID)
		argCount++
	}
	if q.Email != nil {
		query += fmt.Sprintf(" AND email = $%d", argCount)
		args = append(args, *q.Email)
		argCount++
	}
	if q.Provider != nil {
		query += fmt.Sprintf(" AND provider = $%d", argCount)
		args = append(args, *q.Provider)
		argCount++
	}
	if q.ProviderUID != nil {
		query += fmt.Sprintf(" AND provider_uid = $%d", argCount)
		args = append(args, *q.ProviderUID)
		argCount++
	}
	if q.IsVerified != nil {
		query += fmt.Sprintf(" AND is_verified = $%d", argCount)
		args = append(args, *q.IsVerified)
		argCount++
	}

	var auths []*entity.Auth
	if err := r.db.SelectContext(ctx, &auths, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return auths, nil
}
