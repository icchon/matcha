package postgres

import (
	"context"

	"database/sql"

	"errors"

	"github.com/google/uuid"

	"github.com/icchon/matcha/api/internal/domain"
)

type authRepository struct {
	db DBTX
}

func NewAuthRepository(db DBTX) domain.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Save(ctx context.Context, auth *domain.Auth) error {
	query := `
		INSERT INTO auths (user_id, email, provider, provider_uid, password_hash)
		VALUES (:user_id, :email, :provider, :provider_uid, :password_hash)
		ON CONFLICT (provider, provider_uid) DO UPDATE SET
			email = :email,
			password_hash = :password_hash
	`
	_, err := r.db.NamedExecContext(ctx, query, auth)
	return err
}

func (r *authRepository) Delete(ctx context.Context, userID uuid.UUID, provider domain.AuthProvider) error {
	query := "DELETE FROM auths WHERE user_id = $1 AND provider = $2"
	_, err := r.db.ExecContext(ctx, query, userID, provider)
	return err
}

func (r *authRepository) Find(ctx context.Context, userID uuid.UUID, provider domain.AuthProvider) (*domain.Auth, error) {
	var auth domain.Auth
	query := "SELECT * FROM authentications WHERE user_id = $1 AND provider = $2"
	err := r.db.GetContext(ctx, &auth, query, userID, provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Auth, error) {
	var auths []domain.Auth
	query := "SELECT * FROM auths WHERE user_id = $1"
	err := r.db.SelectContext(ctx, &auths, query, userID)
	if err != nil {
		return nil, err
	}
	return auths, nil
}
