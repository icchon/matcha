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

type userTagRepository struct {
	db DBTX
}

func NewUserTagRepository(db DBTX) repo.UserTagRepository {
	return &userTagRepository{db: db}
}

func (r *userTagRepository) Create(ctx context.Context, userTag *entity.UserTag) error {
	query := `
		INSERT INTO user_tags (user_id, tag_id)
		VALUES (:user_id, :tag_id)
		ON CONFLICT (user_id, tag_id) DO NOTHING
	`
	_, err := r.db.NamedExecContext(ctx, query, userTag)
	return err
}

func (r *userTagRepository) Delete(ctx context.Context, userID uuid.UUID, tagID int32) error {
	query := "DELETE FROM user_tags WHERE user_id = $1 AND tag_id = $2"
	_, err := r.db.ExecContext(ctx, query, userID, tagID)
	return err
}

func (r *userTagRepository) Query(ctx context.Context, q *repo.UserTagQuery) ([]*entity.Tag, error) {
	query := "SELECT t.id, t.name FROM tags t JOIN user_tags ut ON t.id = ut.tag_id WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.UserID != nil {
		query += fmt.Sprintf(" AND ut.user_id = $%d", argCount)
		args = append(args, *q.UserID)
		argCount++
	}
	if q.TagID != nil {
		query += fmt.Sprintf(" AND ut.tag_id = $%d", argCount)
		args = append(args, *q.TagID)
		argCount++
	}
	if q.TagName != nil {
		query += fmt.Sprintf(" AND t.name = $%d", argCount)
		args = append(args, *q.TagName)
		argCount++
	}

	var tags []*entity.Tag
	if err := r.db.SelectContext(ctx, &tags, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tags, nil
}
