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

type connectionRepository struct {
	db DBTX
}

func NewConnectionRepository(db DBTX) repo.ConnectionRepository {
	return &connectionRepository{db: db}
}

func (r *connectionRepository) Create(ctx context.Context, connection *entity.Connection) error {
	query := `
		INSERT INTO connections (user1_id, user2_id)
		VALUES (:user1_id, :user2_id)
		ON CONFLICT (user1_id, user2_id) DO NOTHING
		RETURNING *
	`
	return r.db.QueryRowxContext(ctx, query, connection).StructScan(connection)
}

func (r *connectionRepository) Delete(ctx context.Context, user1ID, user2ID uuid.UUID) error {
	query := "DELETE FROM connections WHERE (user1_id = $1 AND user2_id = $2) OR (user1_id = $2 AND user2_id = $1)"
	_, err := r.db.ExecContext(ctx, query, user1ID, user2ID)
	return err
}

func (r *connectionRepository) Find(ctx context.Context, user1ID, user2ID uuid.UUID) (*entity.Connection, error) {
	var connection entity.Connection
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

func (r *connectionRepository) Query(ctx context.Context, q *repo.ConnectionQuery) ([]*entity.Connection, error) {
	query := "SELECT * FROM connections WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.User1ID != nil {
		query += fmt.Sprintf(" AND user1_id = $%d", argCount)
		args = append(args, *q.User1ID)
		argCount++
	}
	if q.User2ID != nil {
		query += fmt.Sprintf(" AND user2_id = $%d", argCount)
		args = append(args, *q.User2ID)
		argCount++
	}
	if q.CreatedAt != nil {
		query += fmt.Sprintf(" AND created_at = $%d", argCount)
		args = append(args, *q.CreatedAt)
		argCount++
	}

	var connections []*entity.Connection
	if err := r.db.SelectContext(ctx, &connections, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return connections, nil
}
