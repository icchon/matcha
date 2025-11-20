package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type userTagRepository struct {
	db DBTX
}

func NewUserTagRepository(db DBTX) domain.UserTagRepository {
	return &userTagRepository{db: db}
}

func (r *userTagRepository) Save(ctx context.Context, userTag *domain.UserTag) error {
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

func (r *userTagRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error) {
	var tags []domain.Tag
	query := `
		SELECT t.id, t.name
		FROM tags t
		INNER JOIN user_tags ut ON t.id = ut.tag_id
		WHERE ut.user_id = $1
	`
	err := r.db.SelectContext(ctx, &tags, query, userID)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
