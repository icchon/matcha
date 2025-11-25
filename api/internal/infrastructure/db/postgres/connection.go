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
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.QueryRowxContext(ctx, connection).StructScan(connection); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// This is not an error in the case of a conflict where nothing is returned.
			// The connection already exists.
			return nil
		}
		return err
	}
	return nil
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

	// Handle User1ID and User2ID as an OR condition if either is provided
	if q.User1ID != nil && q.User2ID != nil {
		query += fmt.Sprintf(" AND ((user1_id = $%d AND user2_id = $%d) OR (user1_id = $%d AND user2_id = $%d))", argCount, argCount+1, argCount+1, argCount)
		args = append(args, *q.User1ID, *q.User2ID, *q.User1ID, *q.User2ID)
		argCount += 4 // Need 4 arguments for the OR condition
	} else if q.User1ID != nil {
		query += fmt.Sprintf(" AND (user1_id = $%d OR user2_id = $%d)", argCount, argCount)
		args = append(args, *q.User1ID)
		argCount++
	} else if q.User2ID != nil {
		query += fmt.Sprintf(" AND (user1_id = $%d OR user2_id = $%d)", argCount, argCount)
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
