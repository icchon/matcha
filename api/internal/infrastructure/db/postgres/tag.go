package postgres

import (
	"context"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type tagRepository struct {
	db DBTX
}

func NewTagRepository(db DBTX) repo.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) GetAll(ctx context.Context) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	query := "SELECT * FROM tags"
	err := r.db.SelectContext(ctx, &tags, query)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
