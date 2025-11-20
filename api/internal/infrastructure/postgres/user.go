package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type userRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, created_at, last_connection)
		VALUES (:id, :created_at, :last_connection)
		ON CONFLICT (id) DO UPDATE SET
			last_connection = :last_connection
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *userRepository) Find(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := "SELECT * FROM users WHERE id = $1"
	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
