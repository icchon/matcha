package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type connectionRepository struct {
	db DBTX
}

func NewConnectionRepository(db DBTX) domain.ConnectionRepository {
	return &connectionRepository{db: db}
}

func (r *connectionRepository) Save(ctx context.Context, connection *domain.Connection) error {
	query := `
		INSERT INTO connections (user1_id, user2_id, created_at)
		VALUES (:user1_id, :user2_id, :created_at)
		ON CONFLICT (user1_id, user2_id) DO NOTHING
	`
	_, err := r.db.NamedExecContext(ctx, query, connection)
	return err
}

func (r *connectionRepository) Delete(ctx context.Context, user1ID, user2ID uuid.UUID) error {
	query := "DELETE FROM connections WHERE (user1_id = $1 AND user2_id = $2) OR (user1_id = $2 AND user2_id = $1)"
	_, err := r.db.ExecContext(ctx, query, user1ID, user2ID)
	return err
}

func (r *connectionRepository) Find(ctx context.Context, user1ID, user2ID uuid.UUID) (*domain.Connection, error) {
	var connection domain.Connection
	query := "SELECT * FROM connections WHERE (user1_id = $1 AND user2_id = $2) OR (user1_id = $2 AND user2_id = $1)"
	err := r.db.GetContext(ctx, &connection, query, user1ID, user2ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &connection, nil
}

func (r *connectionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Connection, error) {
	var connections []domain.Connection
	query := "SELECT * FROM connections WHERE user1_id = $1 OR user2_id = $1"
	err := r.db.SelectContext(ctx, &connections, query, userID)
	if err != nil {
		return nil, err
	}
	return connections, nil
}
